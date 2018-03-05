package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func Put(url string, body interface{}) (string, error) {
	return doRequest("PUT", url, body)
}

func Post(url string, body interface{}) (string, error) {
	return doRequest("POST", url, body)
}

func Get(url string) (string, error) {
	if resp, err := http.Get(url); err != nil {
		return "", err
	} else if response, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} else {
		return string(response), nil
	}
}

func doRequest(method, url string, body interface{}) (string, error) {
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
			} else if resp.StatusCode != http.StatusOK {
				return "", errors.New(string(response))
			} else {
				return string(response), nil
			}
		}
	}
}
