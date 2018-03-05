package domain

//Event define a basic platform event contract
type Event struct {
	Name         string                 `json:"name"`
	Payload      map[string]interface{} `json:"payload"`
	Reproduction map[string]interface{} `json:"reproduction"`
	Reprocessing map[string]interface{} `json:"reprocessing"`
}
