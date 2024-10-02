package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Request is http call wrapper
func Request(client *http.Client, method, url string, queries, headers map[string]string, payload []byte) (int, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return 0, nil, err
	}

	// set queries string
	q := req.URL.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// set headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	statusCode := res.StatusCode

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return statusCode, nil, err
	}

	return statusCode, body, nil
}
