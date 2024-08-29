package application_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestRequestSend(t *testing.T) {
	expect := CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"responderId": "48683858-796c-45ad-a361-9e3d6d003354",
		"message":     "test add",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/friendship/send-request", bytes.NewReader(reqBody))
	addSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"ok": true,
		"data": J{
			"requesterId": "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
			"responderId": "48683858-796c-45ad-a361-9e3d6d003354",
			"message":     "test add",
			"viewed":      false,
			"accepted":    false,
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}

func TestRequestAccept(t *testing.T) {
	expect := CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"requesterId": "1ae02f69-ea1a-4308-b825-0e5896e652e4",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/friendship/accept-request", bytes.NewReader(reqBody))
	addSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"ok": true,
		"data": J{
			"requesterId": "1ae02f69-ea1a-4308-b825-0e5896e652e4",
			"responderId": "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
			"message":     "test request",
			"viewed":      true,
			"accepted":    true,
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}

func TestRequestDeny(t *testing.T) {
	expect := CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"requesterId": "baa0e8bc-385f-4314-9580-29855aff2229",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/friendship/deny-request", bytes.NewReader(reqBody))
	addSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"ok": true,
		"data": J{
			"requesterId": "baa0e8bc-385f-4314-9580-29855aff2229",
			"responderId": "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
			"message":     "test request",
			"viewed":      true,
			"accepted":    false,
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}

func TestRemoveFriend(t *testing.T) {
	expect := CreateMatcher(t, true)
	reqBody, _ := json.Marshal(J{
		"friendId": "37d28fdf-99bc-44b5-8df9-6a3b1a36f177",
	})
	req, _ := http.NewRequest("POST", "http://backend:8080/friendship/remove-friend", bytes.NewReader(reqBody))
	addSessionId(req)
	res, _ := http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	expected := J{
		"ok": true,
		"data": J{
			"requesterId": "37d28fdf-99bc-44b5-8df9-6a3b1a36f177",
			"responderId": "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
			"message":     "test",
			"viewed":      true,
			"accepted":    false,
		},
	}
	actual := J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)

	reqBody, _ = json.Marshal(J{
		"friendId": "2db640fd-7aa4-4bba-8ee6-3935b700297a",
	})
	req, _ = http.NewRequest("POST", "http://backend:8080/friendship/remove-friend", bytes.NewReader(reqBody))
	addSessionId(req)
	res, _ = http.DefaultClient.Do(req)
	expect(res.StatusCode).ToBe(200)
	body, _ = io.ReadAll(res.Body)

	expected = J{
		"ok": true,
		"data": J{
			"requesterId": "2db640fd-7aa4-4bba-8ee6-3935b700297a",
			"responderId": "ed1a3fd0-4d0b-4961-b4cd-cf212357740d",
			"message":     "test",
			"viewed":      true,
			"accepted":    false,
		},
	}
	actual = J{}
	json.Unmarshal(body, &actual)
	expect(actual).ToContain(expected)
}
