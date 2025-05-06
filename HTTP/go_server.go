package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("host", r.Host)

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	_, _ = fmt.Fprintf(w, "Hello World!")
}

func main(){
	http.HandleFunc("/", sayHello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}