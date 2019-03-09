package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime"
)

var (
	MaxWorker    = runtime.NumCPU()
	MaxQueue     = 1000
	RequestQueue chan Request
)

type Request struct {
	URL     *url.URL
	Request *http.Request
}

type Worker struct {
	WorkerPool     chan chan Request
	RequestChannel chan Request
	Quit           chan bool
}

func NewWorker(workPool chan chan Request) Worker {
	return Worker{
		WorkerPool:     workPool,
		RequestChannel: make(chan Request),
		Quit:           make(chan bool),
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("path", req.RequestURI)
	io.WriteString(w, "Hello World!")
}

func main() {
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("list error", err)
	}
}
