package sdk

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/client"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
)

func getUrl() string {
	return fmt.Sprintf("%s://%s:%s/core",
		infra.GetEnv("APICORE_SCHEME", "http"),
		infra.GetEnv("APICORE_HOST", "localhost"),
		infra.GetEnv("APICORE_PORT", "9110"))
}

//GetSolutionIDByEventName on apicore
func GetSolutionIDByEventName(eventName string) (string, error) {
	url := fmt.Sprintf("%s/operation?filter=byEvent&event=%s", getUrl(), eventName)

	arr := make([]*domain.Operation, 0, 0)

	if err := client.GetJSON(url, &arr); err != nil {
		return "", err
	}
	if len(arr) == 0 {
		return "", infra.NewSubscriberNotFoundException(fmt.Sprintf("Event %s has no subscribers", eventName))
	}
	return arr[0].SystemID, nil
}

//EventHasBindings verify if exist any operation that have inbound or outbound event binding with same name
func EventHasBindings(eventName string) (bool, error) {
	url := fmt.Sprintf("%s/operation?filter=bindingEvent&event=%s", getUrl(), eventName)

	arr := make([]*domain.Operation, 0, 0)

	if err := client.GetJSON(url, &arr); err != nil {
		return false, err
	}
	if len(arr) == 0 {
		return false, infra.NewSubscriberNotFoundException(fmt.Sprintf("Event %s has no subscribers", eventName))
	}
	return true, nil
}
