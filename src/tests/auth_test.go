package application_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

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

func TestRegister(t *testing.T) {
	expect := CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"firstName": "test",
		"lastName":  "test",
		"nick":      "some_unusual_nick",
		"password":  "123123123",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/register", bytes.NewReader(reqBody))
	res, _ := http.DefaultClient.Do(req)

	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"errors": []any{},
		"data": J{
			"user": J{
				"firstName": "test",
				"lastName":  "test",
				"nick":      "some_unusual_nick",
			},
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)

	expect(actual).ToContain(expected)
	user := actual["data"].(J)["user"].(J)
	expect(user["id"]).Not().ToBe("")
	expect(user["passhash"]).Not().ToBe("")

	session := actual["data"].(J)["session"].(J)
	expect(session).Not().ToBe(nil)
	expect(session["sessionId"]).Not().ToBe("")
	expect(session["userId"]).ToBe(user["id"])
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

func TestLogout(t *testing.T) {
	expect := CreateMatcher(t, true)
	req, _ := http.NewRequest("POST", "http://backend:8080/logout", nil)
	req.Header.Add("X-Auth-Session-Id", "bb23af03-be50-4bce-b729-b259b2e02e54")
	res, _ := http.DefaultClient.Do(req)

	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"errors": []any{},
		"data": J{
			"sessionId": "bb23af03-be50-4bce-b729-b259b2e02e54",
			"userId":    "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)

	expect(actual).ToContain(expected)
}
