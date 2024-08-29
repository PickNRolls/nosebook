package application_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestCommentPost(t *testing.T) {
	expect := CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/comments/publish-on-post", bytes.NewReader(reqBody))
	addSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"ok":   true,
		"data": J{},
	}
	actual := J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
	expect(actual["data"].(J)["id"]).ToBeTypeOf("string")
	commentId := actual["data"].(J)["id"].(string)
	expect(commentId).Not().ToBe("")

	reqBody, _ = json.Marshal(J{
		"id": commentId,
	})
	req, _ = http.NewRequest("POST", "http://backend:8080/comments/remove", bytes.NewReader(reqBody))
	addSessionId(req)
	res, _ = http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ = io.ReadAll(res.Body)

	expected = J{
		"ok": true,
		"data": J{
			"id": commentId,
		},
	}
	actual = J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}
