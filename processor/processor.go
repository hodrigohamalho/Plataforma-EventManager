package processor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"

	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

type Processor struct {
	dispatcher     bus.Dispatcher
	currentPattern string
	keyOrder       []string
	executionFlow  map[string][]func(*domain.Event) error
	cutOfRules     []func(*domain.Event) error
}

func NewProcessor(dispatcher bus.Dispatcher) *Processor {
	p := new(Processor)
	p.dispatcher = dispatcher
	p.keyOrder = make([]string, 0, 0)
	p.executionFlow = make(map[string][]func(*domain.Event) error)
	p.cutOfRules = make([]func(*domain.Event) error, 0, 0)
	return p
}

func (p *Processor) Where(pattern string) *Processor {
	p.currentPattern = pattern
	p.executionFlow[pattern] = make([]func(*domain.Event) error, 0, 0)
	p.keyOrder = append(p.keyOrder, pattern)
	return p
}

func (p *Processor) CutOff(action func(event *domain.Event) error) *Processor {
	p.cutOfRules = append(p.cutOfRules, action)
	return p
}

func (p *Processor) Execute(action func(event *domain.Event) error) *Processor {
	p.executionFlow[p.currentPattern] = append(p.executionFlow[p.currentPattern], action)
	return p
}

func (p *Processor) Dispatch(routingKey string) *Processor {

	p.executionFlow[p.currentPattern] = append(p.executionFlow[p.currentPattern], func(event *domain.Event) error {
		/*
			Quando um evento no escopo de reprocessamento chega com o branch master eu aplico o processo de split no documento de controle do reprocessamento
			então para cada branch derivada do master também será adicionado, ao documento de controle, os "eventos" splitados
		*/
		if event.Scope == "reprocessing" && event.Branch == "master" {
			reprocessingOrigin := new(domain.ReprocessingData)
			reprocessingOrigin.ParseEvent(event)
			if response, err := sdk.GetDocument("reprocessing", map[string]string{"id": reprocessingOrigin.ID}); err != nil {
				return err
			} else {
				reprocessingList := make([]domain.Reprocessing, 0, 1)
				if err := json.Unmarshal([]byte(response), &reprocessingList); err != nil {
					return infra.NewComponentException(err.Error())
				}
				reprocessing := reprocessingList[0]
				for _, command := range event.Commands {
					clone := new(domain.Event)
					infra.Clone(event, clone)
					clone.Branch = command.Branch
					clone.Reprocessing["executed"] = false
					reprocessing.Events = append(reprocessing.Events, clone)
				}
				sdk.ReplaceDocument("reprocessing", map[string]string{"id": reprocessingOrigin.ID}, reprocessing)
			}
		}
		err := p.dispatcher.Publish(routingKey, event.ToCeleryMessage())
		if len(event.Bindings) > 0 && !event.IsEndingEvent() {
			binding := event.Bindings[0]
			if binding.Reprocessable && event.HasCommands() {
				log.Debug("Process is reprocessable")
				log.Debug("Dispatching splited event")

				for _, command := range event.Commands {
					log.Debug(fmt.Sprintf("   ---> Event %s on Branch %s", command.Name, command.Branch))
					if err = p.dispatcher.Publish(routingKey, command.ToCeleryMessage()); err != nil {
						return infra.NewComponentException(err.Error())
					}
				}
			}
		}
		return err
	})
	return p
}

func (p *Processor) Else(routingKey string) *Processor {
	p.executionFlow["else"] = append(p.executionFlow[p.currentPattern], func(event *domain.Event) error {
		return p.dispatcher.Publish(routingKey, event.ToCeleryMessage())
	})
	return p
}

func (p *Processor) Push(event *domain.Event) (err error) {
	p.currentPattern = ""
	suitable := false
	event.ApplyDefaultFields()
	buildCommands(event)
	log.Info(fmt.Sprintf("Received event %s Scope %s Branch %s", event.Name, event.Scope, event.Branch))
	for _, cutOffRule := range p.cutOfRules {
		if err = cutOffRule(event); err != nil {
			log.Error(fmt.Sprintf("Cutting off event %s with error %s", event.Name, err.Error()))
			return
		}
	}
	publishEventToStore(p.dispatcher, event)
	for _, k := range p.keyOrder {
		actions := p.executionFlow[k]
		if matched, err := regexp.MatchString(convertPattern(k), event.Name); err != nil {
			return infra.NewSystemException(err.Error())
		} else if matched {
			suitable = true
			for _, fnc := range actions {
				if err := fnc(event); err != nil {
					switch t := err.(type) {
					case *infra.Exception:
						return t
					default:
						log.Error(t)
						return infra.NewSystemException(t.Error())
					}
				}
			}
			break
		}
	}
	if !suitable {
		fncs, ok := p.executionFlow["else"]
		if ok && len(fncs) > 0 {
			return fncs[0](event)
		}
		return infra.NewArgumentException(fmt.Sprintf("Event didn't match any clause"))
	}
	return nil
}

func publishEventToStore(dispatcher bus.Dispatcher, event *domain.Event) error {
	if err := dispatcher.Publish("store", event.ToCeleryMessage()); err != nil {
		for _, command := range event.Commands {
			err = dispatcher.Publish("store", command.ToCeleryMessage())
			if err != nil {
				return infra.NewComponentException(err.Error())
			}
		}
		return infra.NewComponentException(err.Error())
	}
	return nil
}

//BuildCommands builds command chain on event
func buildCommands(e *domain.Event) (err error) {
	e.Bindings, err = sdk.EventBindings(e.Name)
	if len(e.Bindings) > 0 {
		systemID := e.Bindings[0].SystemID
		if branches, err := sdk.GetOpenBranchesBySystem(systemID); err != nil {
			return err
		} else {
			for _, branch := range branches {
				command := e.GetCommand()
				command.Branch = branch.Name
				e.Commands = append(e.Commands, command)
			}
		}
	}
	return
}

func convertPattern(s string) string {
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, "*", "\\W*", -1)
	return s
}
