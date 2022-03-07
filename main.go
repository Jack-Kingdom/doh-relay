package main

import (
	"doh-relay/resolver"
	D "github.com/miekg/dns"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func dnsQueryHandler(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("content-type") != resolver.DohMimeType {
		http.Error(w, "content type mismatch", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var m = D.Msg{}
	err = m.Unpack(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	msg, err := resolver.Exchange(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msgBytes, err := msg.Pack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", resolver.DohMimeType)
	_, err = w.Write(msgBytes)
	zap.L().Error("err on write response body", zap.Error(err))
	return
}

func main() {
	http.HandleFunc("/dns-query", dnsQueryHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		zap.L().Error("err on start service", zap.Error(err))
	}
}
