package sdk

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/PMoneda/http"
	log "github.com/sirupsen/logrus"
)

//IsRecording checks if current system is in recording state
func IsRecording(systemID string) (bool, error) {
	scheme := env.Get("REPLAY_SCHEME", "http")
	host := env.Get("REPLAY_HOST", "localhost")
	port := env.Get("REPLAY_PORT", "6081")
	url := fmt.Sprintf("%s://%s:%s/v1/tape/%s/recording", scheme, host, port, systemID)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return false, nil
	}
	if resp.Status != 200 && resp.Status != 404 {
		return false, fmt.Errorf("Error on checking is Recording with systemID=%s data=%s", systemID, string(resp.Body))
	}
	if resp.Status == 404 {
		return false, nil
	}
	return true, nil
}

//RecordEvent records an event synchronously at replay service
func RecordEvent(event *domain.Event) error {
	scheme := env.Get("REPLAY_SCHEME", "http")
	host := env.Get("REPLAY_HOST", "localhost")
	port := env.Get("REPLAY_PORT", "6081")
	url := fmt.Sprintf("%s://%s:%s/v1/tape/%s/record", scheme, host, port, event.SystemID)
	resp, err := http.Post(url, event)
	if err != nil {
		log.Error(err)
		return err
	}
	if resp.Status != 200 && resp.Status != 404 {
		return fmt.Errorf("Error on recording event with systemID=%s data=%s", event.SystemID, string(resp.Body))
	}
	return nil
}
