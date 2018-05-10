package sdk

type ProcessInstance struct {
	Model
	ID     string `json:"id"`
	Status string `json:"status"`
}

func NewProcessInstance() *ProcessInstance {
	procInst := new(ProcessInstance)
	procInst.Metadata.Type = "processInstance"
	return procInst
}
