package eventstore

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
)

func TestShouldShowDatabase(t *testing.T) {
	dbs := showDatabases()
	if len(dbs) == 0 {
		t.Fail()
	}
}

func TestShouldCreateDatabaseWhenNotExist(t *testing.T) {
	createDatabase("test_db")
	if shouldCreateDatabase("test_db") {
		t.Fail()
	}
	executeStatement(`DROP DATABASE "test_db"`)
}

func TestShouldCreateRetentionPolicy(t *testing.T) {
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
}

func TestShouldPushEventToInflux(t *testing.T) {
	createDatabase("event_manager")
	createRetentionPolicy("platform_events", "event_manager")
	e := domain.Event{
		Name: "name",
		Payload: map[string]interface{}{
			"origin":  "1",
			"destiny": "2",
		},
	}
	err := influxPush(e)
	if err != nil {
		t.Fatal(err)
	}
}
