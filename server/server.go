package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fornellas/tasmota_exporter/tasmota"
)

var defaultTimeout = time.Second * 10

func getProbeFn(logger *logrus.Logger) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error

		// Method
		if req.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Target
		if !req.URL.Query().Has("target") {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "missing 'target' parameter")
			return
		}
		tasmotaUrl, err := url.Parse(req.URL.Query().Get("target"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid 'target' parameter: %v", err)
			return
		}
		if tasmotaUrl.RequestURI() != "/" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid 'target' parameter, path must be '/' without query parameters: %v", tasmotaUrl.RequestURI())
			return
		}
		tasmotaUrl.Path = "/cm"
		tasmotaUrl.RawQuery = "cmnd=status%200"

		// Timeout
		timeout := defaultTimeout
		if req.URL.Query().Has("timeout") {
			timeout, err = time.ParseDuration(req.URL.Query().Get("timeout"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `invalid 'timeount', must be integer + unit ("ns", "us" (or "µs"), "ms", "s", "m", "h"): %v`, err)
				return
			}
			if timeout <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, `invalid 'timeount', must be > 0`)
				return
			}
		}

		// GET
		client := http.Client{
			Timeout: timeout,
		}
		logger.Infof("GET %s", tasmotaUrl.String())
		resp, err := client.Get(tasmotaUrl.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `Failed to GET %s: %v`, tasmotaUrl.String(), err)
			return
		}
		if resp.StatusCode != 200 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `GET %s returned !200 status code: %d`, tasmotaUrl.String(), resp.StatusCode)
			return
		}
		if resp.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `GET %s returned Content-Type != application/json: %v`, tasmotaUrl.String(), resp.Header.Get("Content-Type"))
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `GET %s: failed to read body: %v`, tasmotaUrl.String(), err)
			return
		}

		// Parse
		var status tasmota.Status
		if err := json.Unmarshal(body, &status); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `GET %s: failed to parse body: %v`, tasmotaUrl.String(), err)
			return
		}

		// Output
		fmt.Fprintf(w, "%s", status.TimeSeriesGroup().String())
	}
}

func NewServer(addr string, logger *logrus.Logger) http.Server {

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/probe", getProbeFn(logger))

	return http.Server{
		Addr:    addr,
		Handler: serveMux,
	}
}
