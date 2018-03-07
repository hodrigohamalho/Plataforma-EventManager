package domain

//Event define a basic platform event contract
type Event struct {
	Timestamp    string                 `json:"timestamp"`
	Name         string                 `json:"name,omitempty"`
	AppOrigin    string                 `json:"appOrigin,omitempty"`
	Owner        string                 `json:"owner,omitempty"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
	Reproduction map[string]interface{} `json:"reproduction,omitempty"`
	Reprocessing map[string]interface{} `json:"reprocessing,omitempty"`
}
