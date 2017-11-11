package safechat

import (
	"fmt"
	"golang.org/x/net/websocket"
	"encoding/json"
)

type Message struct{
	Action string 		`json:"action"`
	From string			`json:"from"`
	Content string		`json:"content"`
}

type User struct{
	Name string			`json:"name"`
	ws *websocket.Conn	`json:"-"`
}

const BACKEND_USER_NAME = "go-safe-chat-backend"
const ACTION_USER_LIST 	= "userList"
const ACTION_NEW_USER 	= "newUser"
const ACTION_USER_EXIT  = "userLeft"
const ACTION_MESSAGE 	= "message"

var channels = make(map[string][]User)


func initChannel(channelName string, initializedChannel []User ){
	channels[channelName] = initializedChannel
}
func addUserToChannel(channelName string, user User){
	var users = channels[channelName]
	users = append(users, user)
	channels[channelName] = users
}

func broadCastMessage(channelName string, message Message, exceptionUserName string){
	var users, _ = channels[channelName]
	outputJson, _ := json.Marshal(message)

	// Loop over the users in the current room and send the message
	for _, currentUser := range users {

		if currentUser.ws != nil && currentUser.Name != exceptionUserName {
			fmt.Println("Sending to client " + currentUser.Name)
			if err := websocket.Message.Send(currentUser.ws, string(outputJson)); err != nil {
				fmt.Println("Can't send")
			}
		}
	}
}

func userExit(channelName string, currentUserName string){
	userExitMessage := Message{
		Action:  ACTION_USER_EXIT,
		From:    BACKEND_USER_NAME,
		Content: currentUserName,
	}
	broadCastMessage(channelName, userExitMessage, currentUserName)
	for i, v := range channels[channelName] {
		if v.Name == currentUserName {
			v.ws.Close()
			channels[channelName] = append(channels[channelName][:i], channels[channelName][i+1:]...)
			break
		}
	}
}

func debugMap(){
	for k, v := range channels{
		fmt.Println("*******  " + k + "  ***********")
		userList, _ := json.Marshal(v)
		fmt.Println(string(userList))
	}
	fmt.Println("*******          ***********")

}

func Echo(ws *websocket.Conn) {
	var roomName = ws.Config().Protocol[0]

	var emptyChannel []User = []User{
		User{Name: BACKEND_USER_NAME, ws: nil},
	}
	var currentUserName string

	debugMap()
	if _, ok := channels[roomName]; !ok {
		initChannel(roomName, emptyChannel)
	}

	userList, _ := json.Marshal(channels[roomName])
	var userListMessage = Message{
		Action: ACTION_USER_LIST,
		From: BACKEND_USER_NAME,
		Content: string(userList),
	}

	userListJson, err := json.Marshal(userListMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	websocket.Message.Send(ws, string(userListJson))

	for {
		var messageFromClient string
		var outputMessage Message
		if err := websocket.Message.Receive(ws, &messageFromClient); err != nil {
			fmt.Println("Can't receive")
			userExit(roomName, currentUserName)
			break
		}
		fmt.Println("Received back from client: " + messageFromClient)
		var message Message
		if err := json.Unmarshal([]byte(messageFromClient), &message); err != nil {
			fmt.Println(err)
			continue
		}

		switch message.Action {
			case ACTION_NEW_USER :
				currentUserName = message.From
				addUserToChannel(roomName, User{Name: message.From, ws: ws})
				outputMessage = Message{
					Action:  ACTION_NEW_USER,
					From:    BACKEND_USER_NAME,
					Content: message.From,
				}
		    case ACTION_MESSAGE :
				outputMessage = message
		}
		broadCastMessage(roomName, outputMessage, message.From)
		debugMap()
	}
}