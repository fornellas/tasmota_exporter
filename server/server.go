package server

import (
	"io"
	"net/http"
)

func probe(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func NewServer(addr string) http.Server {

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/probe", probe)

	return http.Server{
		Addr:    addr,
		Handler: serveMux,
	}
}
