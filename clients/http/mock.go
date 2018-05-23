package http

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-EventManager/infra"
)

//ReponseMock is mock configure struct
type ReponseMock struct {
	Method        string
	URL           string
	ReponseBody   string
	requestBody   string
	ResponseError error
	executed      int
}

//CalledTimes return how many times mock was called
func (resp *ReponseMock) CalledTimes() int {
	return resp.executed
}

//RequestBody returns request body that mock received
func (resp *ReponseMock) RequestBody() string {
	return resp.requestBody
}

var mocks map[string]*ReponseMock
var once sync.Once

//RegisterMock to return an specific response
func RegisterMock(mock *ReponseMock) {
	once.Do(func() {
		mocks = make(map[string]*ReponseMock)
	})
	key := fmt.Sprintf("%s:%s", mock.Method, mock.URL)
	mocks[key] = mock
}

func doRequestMock(method, url string, body interface{}) (string, error) {
	key := fmt.Sprintf("%s:%s", method, url)
	mock, ok := mocks[key]
	if !ok {
		return "", infra.NewException("test_exception", fmt.Sprintf("mock for %s not defined exception", key))
	}
	mock.executed++
	jsonBody, _ := json.Marshal(body)
	mock.requestBody = string(jsonBody)
	return mock.ReponseBody, mock.ResponseError
}
