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
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"data": nil,
		"errors": []any{
			J{
				"type":    "Not Authenticated",
				"message": "You are not authenticated",
			},
		},
	}
	actual := J{}

	json.Unmarshal(body, &actual)
	expect(actual).ToBe(expected).ElseFail()
}

func TestLogin(t *testing.T) {
	expect := CreateMatcher(t, true)
	req, _ := http.NewRequest("POST", "http://backend:8080/login", strings.NewReader(`{"nick": "test_tester", "password": "123123123" }`))
	res, err := http.DefaultClient.Do(req)

	expect(res.StatusCode).ToBe(200).ElseFail()
	expect(err).ToBe(nil)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"errors": []any{},
		"data": J{
			"user": J{
				"id":        "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
				"firstName": "Test",
				"lastName":  "Tester",
				"nick":      "test_tester",
				"passhash":  "$2a$04$PFIkrnjZ62TLHhcU3a6Breh1sLUVMXzwlrLNo2dqQSTM9d02py.oa",
			},
			"session": J{
				"userId": "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
			},
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)

	expect(actual).ToContain(expected)
}
