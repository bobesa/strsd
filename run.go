package main

import (
	"fmt"
	"./strsd"
	"net/http"
)

func main() {
	fmt.Println("STRSD Server")

	apiPath := "/api/"

	//Prepare game
	myGame := strsd.MakeGame(apiPath);

	//HTTP
	http.Handle("/", http.FileServer(http.Dir("./client/")))
	http.HandleFunc(apiPath, myGame.ProcessRequest)

	//Bots
	users := []string{/*"bobesa","bot","bot2"*/}
	for _, name := range users {
		player := myGame.MakePlayer(name);
		myGame.AddPlayer(player)
	}

	http.ListenAndServe(":8080", nil)
}
