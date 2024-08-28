package application_tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type ResponseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Errors []ResponseError `json:"errors"`
	Data   any             `json:"data"`
}

func TestAuth(t *testing.T) {
	expect := CreateMatcher(t)
	res, _ := http.Get("http://backend:8080/whoami")

	expect(res.StatusCode).ToBe(403).ElseFail()
	body, err := io.ReadAll(res.Body)

	expect(err).ToBe(nil)
	expected := Response{}
	expected.Errors = append(expected.Errors, ResponseError{
		Type:    "Not Authenticated",
		Message: "You are not authenticated",
	})

	actual := Response{}
	err = json.Unmarshal(body, &actual)
	expect(err).ToBe(nil)
	expect(actual).ToDeepEqual(expected).ElseFail()
}
