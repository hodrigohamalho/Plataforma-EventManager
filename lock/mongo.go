package lock

import (
	"fmt"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-EventManager/infra"

	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

func UpMongo() {
	log.Info("Connecting to Mongo")
	s, err := mgo.Dial(infra.GetEnv("MONGO_HOST", "localhost:27017"))
	remaining := 10
	delay := 1 * time.Second
	for err != nil {
		if remaining <= 0 {
			panic("Cannot connect to mongodb EventManager cannot start")
		}
		log.Error(err)
		log.Error(fmt.Sprintf("Trying to connecto to mongodb in %s: remaining retries %d", delay, remaining))
		remaining--
		time.Sleep(delay)
		s, err = mgo.Dial("server1.example.com,server2.example.com")
	}
	log.Info("Connected to Mongo")
	session = s
	session.SetMode(mgo.Monotonic, true)

}

func lockMongo(lock *domain.LockSolution) error {
	c := session.DB("event_manager").C("locks")
	err := c.Insert(lock)
	if err != nil {
		return err
	}
	return nil
}

func unlockMongo(solutionID string) error {
	return session.DB("event_manager").C("locks").Remove(bson.M{"solutionid": solutionID})
}

func findBySolutionID(solutionID string) (*domain.LockSolution, error) {
	lock := new(domain.LockSolution)
	c := session.DB("event_manager").C("locks")
	err := c.Find(bson.M{"solutionid": solutionID}).One(lock)
	if err != nil {
		if err.Error() == "not found" {
			return lock, nil
		}
		return lock, err
	}
	return lock, nil
}
