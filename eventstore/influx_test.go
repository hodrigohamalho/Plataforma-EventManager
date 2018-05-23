package eventstore

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldShowDatabase(t *testing.T) {
	Convey("Should show databases", t, func() {
		dbs := showDatabases()
		if len(dbs) == 0 {
			t.Fail()
		}
	})

}

func TestShouldCreateDatabaseWhenNotExist(t *testing.T) {
	Convey("Should create database when not exist", t, func() {
		createDatabase("test_db")
		if shouldCreateDatabase("test_db") {
			t.Fail()
		}
		executeStatement(`DROP DATABASE "test_db"`)
	})

}

func TestShouldCreateRetentionPolicy(t *testing.T) {
	Convey("Should create retention policy", t, func() {
		createDatabase("test_ret")
		createRetentionPolicy("t_retention", "test_ret")
		policies := showRetentionPolicy("test_ret")
		contains := false
		for _, policy := range policies {
			if policy == "t_retention" {
				contains = true
				break
			}
		}
		if !contains {
			t.Fail()
		}
		executeStatement(`DROP DATABASE "test_ret"`)
	})

}

func TestShouldPushEventToInflux(t *testing.T) {
	Convey("Should push event to influx", t, func() {
		createDatabase("teste")
		createRetentionPolicy("teste", "teste")
		e := domain.Event{
			Name:      "name",
			Owner:     "a",
			AppOrigin: "b",
			Tag:       "1",
			Scope:     "execution",
			Branch:    "master",
			Payload: map[string]interface{}{
				"origin":  "1",
				"destiny": "2",
			},
		}
		err := influxPush(e)
		if err != nil {
			t.Fatal(err)
		}
	})

}

func TestShouldCountEventsInflux(t *testing.T) {

	Convey("Should count events in influx", t, func() {
		createDatabase("teste")
		createRetentionPolicy("teste", "teste")
		pushEvents(5, "evt1", "a", "gol1")
		pushEvents(15, "evt2", "b", "gol1")
		if totalEventsByField("name", "evt1", "1h") != 5 {
			t.Fail()
		}
		if len(queryEvents("name", "evt1", "1h")) != 5 {
			t.Fail()
		}
		executeStatement(`DROP DATABASE "teste"`)
	})

}

func pushEvents(n int, name, owner, appOrigin string) {
	for i := 0; i < n; i++ {
		e := domain.Event{
			Name:      name,
			Owner:     owner,
			AppOrigin: appOrigin,
			Branch:    "master",
			Scope:     "execution",
			Tag:       "tag",
			Payload: map[string]interface{}{
				"origin":  "1",
				"destiny": "2",
			},
		}
		influxPush(e)
	}
}
