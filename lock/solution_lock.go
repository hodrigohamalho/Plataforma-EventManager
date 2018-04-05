package lock

import (
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-EventManager/domain"
)

var memoryDB map[string]bool

var mux = sync.Mutex{}

func init() {
	memoryDB = make(map[string]bool)
}

//SolutionIsLocked verifies if a specific solution in locked to receive events
func SolutionIsLocked(solutionId string) (bool, error) {
	lock, err := findBySolutionID(solutionId)
	if err != nil && err.Error() == "not found" {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return lock.SolutionID == solutionId, err
	}

}

//Lock a solution to receive
func Lock(solutionID string, event *domain.Event) error {
	mux.Lock()
	defer mux.Unlock()
	if lock, err := findBySolutionID(solutionID); err != nil {
		return err
	} else if lock.SolutionID == solutionID {
		return fmt.Errorf("lock already exist for solution %s", solutionID)
	} else {
		nLock := domain.NewLock(solutionID, event)
		return lockMongo(nLock)
	}
	return nil
}

func Unlock(solutionID string) error {
	mux.Lock()
	defer mux.Unlock()

	return unlockMongo(solutionID)
}
