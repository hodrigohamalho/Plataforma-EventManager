package processor

import (
	"regexp"
	"strings"

	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
)

//Processor is a router structure to route events comming from API
type Processor struct {
	dispatcher  bus.Dispatcher
	routes      map[string]func(*Context) error
	middlewares []middleware
}

type middleware struct {
	pattern string
	action  func(*Context) error
}

//NewProcessor creates a new instance of events processor
func NewProcessor(dispatcher bus.Dispatcher) *Processor {
	p := new(Processor)
	p.dispatcher = dispatcher
	p.routes = make(map[string]func(*Context) error)
	p.middlewares = make([]middleware, 0, 0)
	return p
}

//When register a event match pattern to bind to a callback function
func (p *Processor) When(pattern string, callback func(*Context) error) {
	p.routes[pattern] = callback
}

//Use register a middleware based on pattern matching on event name
func (p *Processor) Use(pattern string, callback func(*Context) error) {
	p.middlewares = append(p.middlewares, middleware{pattern: pattern, action: callback})
}

//Push publish an event to Processor router
func (p *Processor) Push(event *domain.Event) (err error) {
	action, err := getMatchActions(event.Name, p.routes)
	middlewares, err := getMiddlewaresByPattern(event.Name, p.middlewares)
	ctx := Context{Event: event, dispatcher: p.dispatcher, Session: make(map[string]string)}
	if action != nil && middlewares != nil {
		for _, middleware := range middlewares {
			if proceed := middleware.action(&ctx); proceed != nil {
				return proceed
			}
		}
	}
	if action != nil {
		return action(&ctx)
	}
	return err
}

func getMiddlewaresByPattern(pattern string, middlewares []middleware) ([]middleware, error) {
	matches := make([]middleware, 0, 0)
	for _, middleware := range middlewares {
		if matched, err := regexp.MatchString(convertPattern(middleware.pattern), pattern); err != nil {
			return nil, infra.NewSystemException(err.Error())
		} else if matched {
			matches = append(matches, middleware)
		}
	}
	return matches, nil
}

func getMatchActions(eventName string, _map map[string]func(*Context) error) (func(*Context) error, error) {
	for k, v := range _map {
		if matched, err := regexp.MatchString(convertPattern(k), eventName); err != nil {
			return nil, infra.NewSystemException(err.Error())
		} else if matched {
			return v, nil
		}
	}
	return nil, nil
}

func convertPattern(s string) string {
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, "*", "\\W*", -1)
	return s
}
