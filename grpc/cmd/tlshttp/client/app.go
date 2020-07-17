package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	clientValidCN = regexp.MustCompile(`^(dev[0-9]{1,3})$`)
)

func main() {
	go runClient()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	_ = <-c
	log.Println("Goodbye...")
}

func verifyServerCert(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	log.Print("Server cert verification routine called...")
	return nil
}

func tlsClientConfig() *tls.Config {
	pubCert, err := ioutil.ReadFile("../../../tools/cert/client-dev00/client-dev00-cert.cer")
	if err != nil {
		log.Fatal(err)
	}
	pubCertKey, err := ioutil.ReadFile("../../../tools/cert/client-dev00/client-dev00-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	tlsCert, err := tls.X509KeyPair(pubCert, pubCertKey)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
		Certificates:             []tls.Certificate{tlsCert},
		InsecureSkipVerify:       false,
		VerifyPeerCertificate:    verifyServerCert,
	}
	return cfg
}

func runClient() {
	tlsConfig := tlsClientConfig()
	httpCli := &fasthttp.HostClient{
		Addr:      "localhost:48525",
		IsTLS:     true,
		TLSConfig: tlsConfig}
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()
	req.Header.SetMethod("GET")
	req.SetRequestURI("https://localhost:48525/hello")
	//_ = req.URI()
	for {
		time.Sleep(3 * time.Second)
		log.Printf("\n\n=========\n")
		log.Println("Doing client request...")
		err := httpCli.Do(req, res)
		if err == nil {
			resData := &struct{ Message string }{}
			json.Unmarshal(res.Body(), resData)
			log.Printf("Client request success: %v\n", resData)

		} else {
			log.Printf("Client request fail: %v\n", err)
		}
	}
}
