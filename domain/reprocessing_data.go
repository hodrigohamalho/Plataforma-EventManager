package domain

type ReprocessingData struct {
	ID       string
	Executed bool
	Tag      string
	Branch   string
}

//ParseEvent get data from reprocessing envelop
func (r *ReprocessingData) ParseEvent(event *Event) {
	executed, ok := event.Reprocessing["executed"]
	if !ok {
		r.Executed = false
	} else {
		r.Executed = executed.(bool)
	}

	tag, ok := event.Reprocessing["event_tag"]
	if !ok {
		r.Tag = ""
	} else {
		r.Tag = tag.(string)
	}

	id, ok := event.Reprocessing["id"]
	if !ok {
		r.ID = ""
	} else {
		r.ID = id.(string)
	}

	r.Branch = event.Branch
}
