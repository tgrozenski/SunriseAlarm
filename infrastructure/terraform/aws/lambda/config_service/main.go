package main

import (
	"fmt"
	"net/http"
)

func main_route(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello world\n")
}

func main() {
	http.HandleFunc("/", main_route)
	http.ListenAndServe(":8080", nil)
	//lambda.Start(handler)
}
