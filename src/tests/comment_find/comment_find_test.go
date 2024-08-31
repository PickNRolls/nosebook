package postfind

import (
	"io"
	"net/http"
	"nosebook/src/tests/testlib"
	"testing"
)

func TestCommentFind(t *testing.T) {
	expect := testlib.CreateMatcher(t, true)
	req, _ := http.NewRequest("GET", "http://backend:8080/comments?postId=c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7", nil)
	testlib.AddSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expectedJSONs := TestCommentFindJSONs()
	expect(string(body)).ToBe(expectedJSONs[0])

	req, _ = http.NewRequest("GET", "http://backend:8080/comments?postId=c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7&next=d0023f4d-8d7f-4907-9438-d2ed2a9661f4/2024-02-16T15:40:55Z", nil)
	testlib.AddSessionId(req)
	res, _ = http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ = io.ReadAll(res.Body)

	expect(string(body)).ToBe(expectedJSONs[1])
}
