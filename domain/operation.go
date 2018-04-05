package domain

type Operation struct {
	Event     string `json:"event_in"`
	SystemID  string `json:"systemId"`
	ProcessID string `json:"processId"`
}
