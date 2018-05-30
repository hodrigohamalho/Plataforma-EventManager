package eventstore

import (
	"testing"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldShowDatabase(t *testing.T) {
	Convey("should show databases", t, func() {
		dbs := showDatabases()
		So(len(dbs), ShouldBeGreaterThan, 0)
	})

}

func TestShouldCreateDatabaseWhenNotExist(t *testing.T) {
	Convey("should create database when not exist", t, func() {
		createDatabase("test_db")
		So(shouldCreateDatabase("test_db"), ShouldBeFalse)
		executeStatement(`DROP DATABASE "test_db"`)
	})

}

func TestShouldCreateRetentionPolicy(t *testing.T) {
	Convey("should create retention policy", t, func() {
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
	Convey("should push event to influx", t, func() {
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
		So(err, ShouldBeNil)
	})

}

func TestShouldCountEventsInflux(t *testing.T) {

	Convey("should count events in influx", t, func() {
		createDatabase("teste")
		createRetentionPolicy("teste", "teste")
		pushEvents(5, "evt1", "a", "gol1")
		pushEvents(15, "evt2", "b", "gol1")
		So(Count("name", "evt1", "1h"), ShouldEqual, 5)
		So(len(Query("name", "evt1", "1h")), ShouldEqual, 5)
		executeStatement(`DROP DATABASE "teste"`)
	})

}

func TestShouldInstallInflux(t *testing.T) {

	Convey("should count events in influx", t, func() {
		Install()
		i := 10
		ok := false
		for i >= 0 {
			ok = len(showDatabases()) > 0
			if ok {
				break
			}
			time.Sleep(1 * time.Second)
		}
		So(ok, ShouldBeTrue)
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
		Push(e)
	}
}
