package postfind

import (
	"encoding/json"
	"io"
	"net/http"
	"nosebook/src/tests/testlib"
	"testing"
)

type J = testlib.J

func TestPostFind(t *testing.T) {
	expect := testlib.CreateMatcher(t, true)
	req, _ := http.NewRequest("GET", "http://backend:8080/posts?ownerId=1ae02f69-ea1a-4308-b825-0e5896e652e4", nil)
	testlib.AddSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"ok": true,
		"data": J{
			"data": []any{
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e2",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 5",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c842",
					"author": J{
						"id":        "48683858-796c-45ad-a361-9e3d6d003354",
						"firstName": "Marina",
						"lastName":  "Graf",
						"nick":      "mmm",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 5",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-011a90f7c8e2",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 11",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-001a90f7c8e2",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 10",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-901a90f7c8e0",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 9",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e0",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 4",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-801a90f7c8e9",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 8",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e9",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 3",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c829",
					"author": J{
						"id":        "48683858-796c-45ad-a361-9e3d6d003354",
						"firstName": "Marina",
						"lastName":  "Graf",
						"nick":      "mmm",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 3",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-701a90f7c8e8",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 7",
				},
			},
			"next": "c7b7bf17-38f9-4ed5-b0a8-701a90f7c8e8/2024-02-16T14:36:48Z",
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)

	next := actual["data"].(J)["next"].(string)
	req, _ = http.NewRequest("GET", "http://backend:8080/posts?ownerId=1ae02f69-ea1a-4308-b825-0e5896e652e4&cursor="+next, nil)
	testlib.AddSessionId(req)
	res, _ = http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ = io.ReadAll(res.Body)

	expected = J{
		"ok": true,
		"data": J{
			"data": []any{
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e8",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 2",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c818",
					"author": J{
						"id":        "48683858-796c-45ad-a361-9e3d6d003354",
						"firstName": "Marina",
						"lastName":  "Graf",
						"nick":      "mmm",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 2",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-601a90f7c8e7",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message 6",
				},
				J{
					"id": "c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7",
					"author": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"owner": J{
						"id":        "1ae02f69-ea1a-4308-b825-0e5896e652e4",
						"firstName": "Ilya",
						"lastName":  "Blinkov",
						"nick":      "drugtester",
					},
					"message": "post message",
				},
			},
		},
	}
	actual = J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}

func TestPostComments(t *testing.T) {
	expect := testlib.CreateMatcher(t, true)
	req, _ := http.NewRequest("GET", "http://backend:8080/posts?ownerId=1ae02f69-ea1a-4308-b825-0e5896e652e4&authorId=1ae02f69-ea1a-4308-b825-0e5896e652e4&cursor=c7b7bf17-38f9-4ed5-b0a8-601a90f7c8e7/2024-02-16T14:36:38Z", nil)
	testlib.AddSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expectedJSON := TestPostCommentsJSON()
	expect(string(body)).ToBe(expectedJSON)
}
