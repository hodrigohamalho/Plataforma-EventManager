package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//Put make a PUT request
func Put(url string, body interface{}) (string, error) {
	return doRequest("PUT", url, body)
}

//Post make a POST request
func Post(url string, body interface{}) (string, error) {
	return doRequest("POST", url, body)
}

//Get make a GET request
func Get(url string) (string, error) {
	if resp, err := http.Get(url); err != nil {
		return "", err
	} else if response, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} else {
		return string(response), nil
	}
}

//GetJSON make a GET request and unmarshal response to JSON
func GetJSON(url string, obj interface{}) error {
	if resp, err := http.Get(url); err != nil {
		return err
	} else if response, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else {
		return json.Unmarshal(response, obj)
	}
}

func doRequest(method, url string, body interface{}) (string, error) {
	if mocks != nil {
		return doRequestMock(method, url, body)
	}
	return httpRequest(method, url, body)
}

func httpRequest(method, url string, body interface{}) (string, error) {
	client := http.DefaultClient
	if j, err := json.Marshal(body); err != nil {
		return "", err
	} else if req, err := http.NewRequest(method, url, strings.NewReader(string(j))); err != nil {
		return "", err
	} else {
		req.Header["Content-Type"] = []string{"application/json"}
		if resp, err := client.Do(req); err != nil {
			return "", err
		} else {
			if response, err := ioutil.ReadAll(resp.Body); err != nil {
				return "", err
			} else if resp.StatusCode >= 300 {
				return "", fmt.Errorf("Status %s: response: %s", resp.Status, string(response))
			} else {
				return string(response), nil
			}
		}
	}
}
