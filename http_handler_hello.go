package main

import "net/http"

func handlerGetHello(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, World!\n"))
}
