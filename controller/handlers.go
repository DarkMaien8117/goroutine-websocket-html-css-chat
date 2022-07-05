package controller

import (
	"log"
	"net/http"

	"Maien/model"
)

var SERVER_PORT = ":8117"

// This is basically a init function for the server
func Serve() {
	HandleView()

	http.HandleFunc(model.SOCKET_PATH, handlerSocket) //It is default Golang function

	go model.HandleMessages()

	log.Println("Server running at ", SERVER_PORT)

	if err := http.ListenAndServe(SERVER_PORT, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}

// Here will be serve our client
func HandleView() {
	//It will be set browser's url . So you must http://192.168.104.30:8117/ to get service.
	fs := http.StripPrefix("/", http.FileServer(http.Dir("View")))

	http.Handle("/", fs)
}

// Socket Handler

//It is default Golang function
func handlerSocket(w http.ResponseWriter, r *http.Request) {
	// We use the HTTP Conn as a socket one
	ws, err := model.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	// The conn will close when the app finishes
	defer ws.Close()

	// We add our client
	model.Clients[ws] = true

	for {
		var msg model.Message

		//When accept one client's message...
		err := ws.ReadJSON(&msg) //set the message at msg variable

		// In error case we close connection with our client
		if err != nil {
			log.Printf("Connection Err: %v", err)
			delete(model.Clients, ws)

			break
		}

		model.Broadcast <- msg //Upgrade model/socket.go->Broadcast variable
	}
}
