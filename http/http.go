package http

import (
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Call(path string, method string, headers map[string]string, params map[string]string, body string) ([]byte, error) {
	client := &http.Client{}

	urlVal, err := url.Parse(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}
	// path values
	urlValues := urlVal.Query()
	for k, v := range params {
		key, val := k, v
		urlValues.Set(key, val)
	}

	var reader io.Reader
	if method == http.MethodPost && body != "" {
		reader = strings.NewReader(body)
	}

	urlVal.RawQuery = urlValues.Encode()
	urlPath := urlVal.String()
	req, _ := http.NewRequest(method,
		urlPath, reader)

	// header
	for k, v := range headers {
		key, val := k, v
		req.Header.Add(key, val)
	}

	// do request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "client.Do(req) failed")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "io.ReadAll(resp.Body) failed")
	}

	return respBody, nil
}
