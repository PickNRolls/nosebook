package application_tests

import "net/http"

func addSessionId(req *http.Request) {
	req.Header.Add("X-Auth-Session-Id", "bb23af03-be50-4bce-b729-b259b2e02e54")
}
