package domain

import (
	"time"
)

type LockSolution struct {
	SolutionID string
	Event      *Event
	Timestamp  time.Time
}

func NewLock(solutionID string, event *Event) *LockSolution {
	l := new(LockSolution)
	l.SolutionID = solutionID
	l.Event = event
	l.Timestamp = time.Now().UTC()
	return l
}
