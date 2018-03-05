package executor

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/client"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
)

//PushEvent to executor
func PushEvent(event domain.Event) error {
	url := fmt.Sprintf("%s/event", baseUrl())
	_, err := client.Post(url, event)
	return err
}

func baseUrl() string {
	scheme := infra.GetEnv("EXECUTOR_SCHEME", "http")
	host := infra.GetEnv("EXECUTOR_HOST", "localhost")
	port := infra.GetEnv("EXECUTOR_PORT", "8000")
	return fmt.Sprintf("%s://%s:%s", scheme, host, port)
}
