package sdk

import (
	"fmt"

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
