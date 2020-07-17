package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	clientValidCN = regexp.MustCompile(`^(dev[0-9]{1,3})$`)
)

func main() {
	log.Println(os.Getwd())
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	tlsConfig := tlsServerConfig()
	log.Println("Listening...")
	tlsLn, _ := tls.Listen("tcp", ":48525", tlsConfig)
	defer tlsLn.Close()
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
		Certificates:          []tls.Certificate{tlsCert},
		ClientAuth:            tls.RequireAndVerifyClientCert,
		VerifyPeerCertificate: verifyClientCert,
	}
	return cfg
}

func verifyClientCert(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	log.Print("Client cert verification routine called...")
	if verifiedChains == nil {
		return errors.New("client certificate not provided")
	}
	log.Printf("client cert verification: verified certs len %d", len(verifiedChains))
	var authorizedCertFound = false
	for i := 0; i < len(verifiedChains); i++ {
		chain := verifiedChains[i]
		for y := 0; y < len(chain); y++ {
			cert := chain[y]
			if !authorizedCertFound && clientValidCN.MatchString(cert.Subject.CommonName) {
				authorizedCertFound = true
				log.Printf("client cert chain verification: AUTHORIZED cert (%d, %d) — %s, %d, %s", i, y, cert.Subject, cert.SerialNumber, cert.Subject.CommonName)
			} else {
				log.Printf("client cert chain verification: cert (%d, %d) — %s, %d, %s", i, y, cert.Subject, cert.SerialNumber, cert.Subject.CommonName)
			}
		}
	}
	if !authorizedCertFound {
		return errors.New("no any certificate presented by the client can be authorized using the defined rule(s)")
	}

	// if rawCerts == nil {
	// 	log.Print("Client cert verification: no certs")
	// } else {
	// 	log.Printf("Client cert verification: certs len %d", len(rawCerts))
	// 	for i := 0; i < len(rawCerts); i++ {
	// 		cert, err := x509.ParseCertificate(rawCerts[i])
	// 		if err == nil {
	// 			log.Printf("Client cert verification: cert (%d) — %s, %d, %s", i, cert.Subject, cert.SerialNumber, cert.Subject.CommonName)
	// 		} else {
	// 			log.Printf("Client cert verification: error occurred while parsing cert: %v", err)
	// 		}
	// 	}
	// }

	return nil
}
