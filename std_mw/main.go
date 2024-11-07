package main

import (
	"log"
	"net/http"
	"time"
)

var limitCh = make(chan struct{}, 10)

func limitMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitCh <- struct{}{}
		next.ServeHTTP(w, r)
		<-limitCh
	})
}

func timeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(begin)
		log.Printf("%s use time %d ms\n", r.URL.Path, elapsed.Milliseconds())
	})
}

func boy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi boy"))
}

func girl(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi girl"))
}

func main() {
	http.Handle("/boy", timeMiddleWare(limitMiddleWare(http.HandlerFunc(boy))))
	http.Handle("/girl", timeMiddleWare(limitMiddleWare(http.HandlerFunc(girl)))) // too many nests
	http.ListenAndServe("127.0.0.1:8000", nil)
}
