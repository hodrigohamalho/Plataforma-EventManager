package domain

//SplitState is a struct to control splited events execution state
type SplitState struct {
	Tag    string       `json:"tag"`
	Events []EventState `json:"events"`
}

//Pending is a pending status
const Pending = "pending"

//Error is a error status when events ends with error or exception event
const Error = "error"

//Success is a succes status when a event is executed with no errros with done event
const Success = "success"

//EventState represents an execution state for some event
type EventState struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Branch string `json:"branch"`
	Scope  string `json:"scope"`
}
