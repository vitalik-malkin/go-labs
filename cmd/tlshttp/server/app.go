package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

func main() {
	log.Println(os.Getwd())
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	tlsConfig := tlsServerConfig()
	log.Println("Listening...")
	tlsLn, _ := tls.Listen("tcp", ":48525", tlsConfig)
	defer tlsLn.Close()
	go runClient()
	log.Fatal(http.Serve(tlsLn, mux))
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("Request accepted...")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{ Message string }{`Hello & welcome`})
}

func tlsServerConfig() *tls.Config {
	pubCert, err := ioutil.ReadFile("../../../tools/cert/server-localhost01/server-localhost01-cert.cer")
	if err != nil {
		log.Fatal(err)
	}
	pubCertKey, err := ioutil.ReadFile("../../../tools/cert/server-localhost01/server-localhost01-key.pem")
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
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth:   tls.NoClientCert,
	}
	return cfg
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
		InsecureSkipVerify:       true,
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
	// var reqURI fasthttp.URI
	// reqURI.Parse(nil, []byte("https://localhost:48525/hello"))
	// req.SetRequestURI(reqURI.String())
	req.SetRequestURI("https://localhost:48525/hello")
	_ = req.URI()
	for {
		time.Sleep(10 * time.Second)
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
