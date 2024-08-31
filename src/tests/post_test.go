package application_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"nosebook/src/tests/testlib"
	"testing"
)

func TestPost(t *testing.T) {
	expect := testlib.CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"ownerId": "1ae02f69-ea1a-4308-b825-0e5896e652e4",
		"message": "my test message",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/posts/publish", bytes.NewReader(reqBody))
	testlib.AddSessionId(req)
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
	postId := actual["data"].(J)["id"].(string)
	expect(postId).Not().ToBe("")

	reqBody, _ = json.Marshal(J{
		"id": postId,
	})
	req, _ = http.NewRequest("POST", "http://backend:8080/posts/remove", bytes.NewReader(reqBody))
	testlib.AddSessionId(req)
	res, _ = http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ = io.ReadAll(res.Body)

	expected = J{
		"ok": true,
		"data": J{
			"id": postId,
		},
	}
	actual = J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}
