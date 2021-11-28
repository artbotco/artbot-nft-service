package main

import (
	"net/http"
	routes "github.com/artbotco/artbot-nft-service/routes"
)

func main() {
	http.HandleFunc("/", routes.Handler)
	http.HandleFunc("/about/", routes.About)
	http.ListenAndServe(":8080", nil)

	// Hello world, the web server

	//helloHandler := func(w http.ResponseWriter, req *http.Request) {
	//	io.WriteString(w, "Hello, world!\n")
	//}
	//
	//http.HandleFunc("/hello", helloHandler)
    //log.Println("Listing for requests at http://localhost:8000/hello")
	//log.Fatal(http.ListenAndServe(":8000", nil))
}
