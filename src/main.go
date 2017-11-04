package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"safechat"
	"fmt"
)




func main() {

	http.Handle("/", websocket.Handler(safechat.Echo))

	fmt.Println("App started at port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
