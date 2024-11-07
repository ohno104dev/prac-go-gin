package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{
		middlewareChain: make([]middleware, 0, 10),
		mux:             make(map[string]http.Handler, 10),
	}
}

func (router *Router) Use(m middleware) {
	router.middlewareChain = append(router.middlewareChain, m)
}

func (router *Router) Add(path string, handler http.Handler) {
	var mergedHandler = handler
	// 先執行最內層 limitMiddleWare(timeMiddleWare(handler))
	for i := len(router.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = router.middlewareChain[i](mergedHandler) // 中間件層層嵌套
	}
	router.mux[path] = mergedHandler
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.RequestURI() // 取得請求路徑
	if handler, exists := router.mux[requestPath]; !exists {
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		handler.ServeHTTP(w, r)
	}
}

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

func getBoy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi boy"))
}

func getGirl(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi girl"))
}

func main() {
	router := NewRouter()
	// 先use外層: limitMiddleWare(timeMiddleWare(...))
	router.Use(limitMiddleWare)
	router.Use(timeMiddleWare)
	router.Add("/", http.HandlerFunc(getBoy))
	router.Add("/home", http.HandlerFunc(getGirl))

	if err := http.ListenAndServe("127.0.0.1:8000", router); err != nil {
		fmt.Println(err)
	}
}
