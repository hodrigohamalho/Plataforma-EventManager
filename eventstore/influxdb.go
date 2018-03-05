package eventstore

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
)

func init() {
	if shouldCreateDatabase("event_manager") {
		createDatabase("event_manager")
		createRetentionPolicy("platform_events", "event_manager")
	}
}

//Push event to Influx
func influxPush(event domain.Event) error {
	s := influxWrite(compileEventToLintProtocol(event))
	if s != "" {
		return errors.New(s)
	}
	fmt.Println(s)
	return nil
}

func compileEventToLintProtocol(event domain.Event) string {
	//cpu_load_short,host=server02,region=us-west value=0.55 1422568543702900257
	s, _ := json.Marshal(event)
	encoded := base64.StdEncoding.EncodeToString(s)
	return fmt.Sprintf(`events,name=%s count=1,data="%s"`, event.Name, string(encoded))
}

func getBaseUrl() string {
	return fmt.Sprintf("http://%s:%s", infra.GetEnv("INFLUX_HOST", "localhost"), infra.GetEnv("INFLUX_PORT", "8086"))
}

func createDatabase(name string) {
	executeStatement(fmt.Sprintf("CREATE DATABASE %s", name))
}

func createRetentionPolicy(name, db string) {
	executeStatement(fmt.Sprintf("CREATE RETENTION POLICY %s ON %s DURATION 500w REPLICATION 1", name, db))
}

func shouldCreateDatabase(name string) bool {
	dbs := showDatabases()
	for _, db := range dbs {
		if db == name {
			return false
		}
	}
	return true
}

func showRetentionPolicy(dbname string) []string {
	r := executeStatement(fmt.Sprintf("SHOW RETENTION POLICIES ON %s", dbname))
	jsonParsed, _ := gabs.ParseJSON([]byte(r))
	b := jsonParsed.Path("results.series.values").String()
	//TODO refactor
	b = strings.Replace(b, "[", "", -1)
	b = strings.Replace(b, "]", "", -1)
	b = strings.Replace(b, "\"", "", -1)
	return strings.Split(b, ",")
}

func showDatabases() []string {
	r := executeStatement("SHOW DATABASES")
	jsonParsed, _ := gabs.ParseJSON([]byte(r))
	b := jsonParsed.Path("results.series.values").String()
	//TODO refactor
	b = strings.Replace(b, "[", "", -1)
	b = strings.Replace(b, "]", "", -1)
	b = strings.Replace(b, "\"", "", -1)
	return strings.Split(b, ",")
}

func executeStatement(stmt string) string {
	_url := fmt.Sprintf("%s/query?u=%s&p=%s", getBaseUrl(), infra.GetEnv("INFLUX_USER", ""), infra.GetEnv("INFLUX_PASSWORD", ""))
	payload := strings.NewReader(url.PathEscape(fmt.Sprintf("q=%s", stmt)))
	req, _ := http.NewRequest("POST", _url, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	return string(b)
}

func influxWrite(point string) string {
	_url := fmt.Sprintf("%s/write?u=%s&p=%s&db=%s&rp=%s", getBaseUrl(), infra.GetEnv("INFLUX_USER", ""), infra.GetEnv("INFLUX_PASSWORD", ""), "event_manager", "platform_events")
	payload := strings.NewReader(point)
	req, _ := http.NewRequest("POST", _url, payload)
	res, _ := http.DefaultClient.Do(req)
	st := res.StatusCode
	fmt.Println(st)
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	return string(b)
}
