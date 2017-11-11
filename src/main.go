package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"safechat"
	"fmt"
)




func main() {

	http.Handle("/wsapp", websocket.Handler(safechat.Echo))

	fmt.Println("App started at port 8090")

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
