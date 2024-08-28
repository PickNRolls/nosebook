package application_tests

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type ResponseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Response[T any] struct {
	Errors []ResponseError `json:"errors"`
	Data   T               `json:"data"`
}

type UserDTO struct {
	Id             uuid.UUID `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Nick           string    `json:"nick"`
	Passhash       string    `json:"passhash"`
	CreatedAt      time.Time `json:"createdAt"`
	LastActivityAt time.Time `json:"lastActivityAt"`
}

type SessionDTO struct {
	SessionId uuid.UUID `json:"sessionId"`
	UserId    uuid.UUID `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type LoginDTO struct {
	User    UserDTO    `json:"user"`
	Session SessionDTO `json:"session"`
}

func TestNotAuthenticated(t *testing.T) {
	expect := CreateMatcher(t, false)
	res, _ := http.Get("http://backend:8080/whoami")

	expect(res.StatusCode).ToBe(403).ElseFail()
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	expect(err).ToBe(nil)
	expected := Response[struct{}]{}
	expected.Errors = append(expected.Errors, ResponseError{
		Type:    "Not Authenticated",
		Message: "You are not authenticated",
	})

	actual := Response[struct{}]{}
	err = json.Unmarshal(body, &actual)
	expect(err).ToBe(nil)
	expect(actual).ToDeepEqual(expected).ElseFail()
}

func TestLogin(t *testing.T) {
	expect := CreateMatcher(t, true)
	req, _ := http.NewRequest("POST", "http://backend:8080/login", strings.NewReader(`{"nick": "test_tester", "password": "123123123" }`))
	res, err := http.DefaultClient.Do(req)

	expect(res.StatusCode).ToBe(200).ElseFail()
	expect(err).ToBe(nil)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	expectedResponse := Response[LoginDTO]{}
	json.Unmarshal([]byte(`
		{
    		"errors": [],
    		"data": {
        		"user": {
            		"id": "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
            		"firstName": "Test",
            		"lastName": "Tester",
            		"nick": "test_tester",
            		"passhash": "$2a$04$PFIkrnjZ62TLHhcU3a6Breh1sLUVMXzwlrLNo2dqQSTM9d02py.oa",
            		"createdAt": "2024-08-28T13:00:39.440309Z",
            		"lastActivityAt": "2024-08-28T13:00:39.440309Z"
        		},
        		"session": {
            		"sessionId": "b58eb155-5933-43eb-b905-a9edfe8d0744",
            		"userId": "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
            		"createdAt": "2024-08-28T13:06:30.624200793Z",
            		"expiresAt": "2024-08-30T13:06:30.624200793Z"
        		}
    		}
		}
	`), &expectedResponse)

	actual := Response[LoginDTO]{}
	err = json.Unmarshal(body, &actual)
	expect(err).ToBe(nil)

	expect(len(actual.Errors)).ToBe(0)

	expected := expectedResponse.Data
	expect(actual.Data.User.Id).ToBe(expected.User.Id)
	expect(actual.Data.User.FirstName).ToBe(expected.User.FirstName)
	expect(actual.Data.User.LastName).ToBe(expected.User.LastName)
	expect(actual.Data.User.Nick).ToBe(expected.User.Nick)
	expect(actual.Data.User.Passhash).ToBe(expected.User.Passhash)

	expect(actual.Data.Session.UserId).ToBe(expected.Session.UserId)
}
