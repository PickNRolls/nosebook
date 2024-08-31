package application_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"nosebook/src/tests/testlib"
	"testing"
)

func TestLikePost(t *testing.T) {
	expect := testlib.CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/like/post", bytes.NewReader(reqBody))
	testlib.AddSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"ok": true,
		"data": J{
			"postId": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7",
			"liked":  true,
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)

	res, _ = http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ = io.ReadAll(res.Body)

	expected = J{
		"ok": true,
		"data": J{
			"postId": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7",
			"liked":  false,
		},
	}
	actual = J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}

func TestLikeComment(t *testing.T) {
	expect := testlib.CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"id": "620c79b7-3927-48b7-a308-1ffd3db6036f",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/like/comment", bytes.NewReader(reqBody))
	testlib.AddSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"ok": true,
		"data": J{
			"commentId": "620c79b7-3927-48b7-a308-1ffd3db6036f",
			"liked":     true,
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)

	res, _ = http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ = io.ReadAll(res.Body)

	expected = J{
		"ok": true,
		"data": J{
			"commentId": "620c79b7-3927-48b7-a308-1ffd3db6036f",
			"liked":     false,
		},
	}
	actual = J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}
