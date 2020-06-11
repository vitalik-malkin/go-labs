package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	tlsConfig := tlsConfig()
	log.Println("Listening...")
	tlsLn, _ := tls.Listen("tcp", ":48525", tlsConfig)
	defer tlsLn.Close()
	log.Fatal(http.Serve(tlsLn, mux))
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{ Message string }{"Hello & welcome"})
}

func tlsConfig() *tls.Config {

}
