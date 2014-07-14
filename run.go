/*
	Default runnable of STRSD
	This provides you with:
	1) single game instance
	2) user game instance's ProcessRequest func with http.HandleFunc
	3) preloading of players
*/
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

	//Player preload
	playerNames := []string{/* enter usernames here */}
	for _, name := range playerNames {
		player := myGame.MakePlayer(name);
		myGame.AddPlayer(player)
	}

	http.ListenAndServe(":8080", nil)
}
