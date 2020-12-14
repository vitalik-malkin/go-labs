package main

import (
	"crypto/tls"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	fh "github.com/valyala/fasthttp"
)

func main() {
	repro()
}

func repro() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	level.Info(logger).Log("msg", "initializing")

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       false,
	}
	client := &fh.HostClient{
		Addr:      "en.wikipedia.org",
		IsTLS:     true,
		TLSConfig: tlsConfig}

	req, resp := fh.AcquireRequest(), fh.AcquireResponse()
	defer func() {
		fh.ReleaseRequest(req)
		fh.ReleaseResponse(resp)
	}()

	req.Header.SetMethod("GET")
	req.Header.Set("Content-Type", "plain/text")
	req.SetRequestURI("https://en.wikipedia.org/wiki/Web_API")

	err := client.Do(req, resp)
	if err != nil {
		level.Error(logger).Log("req_err", err)
	}
}
