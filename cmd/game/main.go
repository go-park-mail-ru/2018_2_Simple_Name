package main

import (
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//"SimpleGame/2018_2_Simple_Name/internal/session"
	"SimpleGame/internal/session"
	"SimpleGame/internal/db/postgres"
	"SimpleGame/internal/game"
	//"SimpleGame/session"
	"fmt"
	"github.com/gorilla/websocket"
	//"google.golang.org/grpc"
	"log"
//	"net"
	"net/http"
	//"strconv"
	//"time"
)

//var sessManager session.AuthCheckerClient
//var ctx context.Context

var gameService = game.NewGame()

func main() {

	err := db.OpenConn()

	if err != nil {
		return
	}

	grpcConn, err := session.OpenConn()

	if err != nil || grpcConn == nil {
		return
	}

	defer grpcConn.Close()

	mux := http.NewServeMux()

	go gameService.Run()

	mux.HandleFunc("/startgame", startGame)

	fmt.Println("Starting game server at :8082")

	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}

}



func startGame(w http.ResponseWriter, r *http.Request) {
	//sugar.Info("Startgame signal from user")
	sess, err := session.FindSession(r)
	if err != nil {
		fmt.Println("Failed get session", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(sess.Email)

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	upgrader := websocket.Upgrader{}
	//upgrader.CheckOrigin = true
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
		//origin := r.Header["Origin"]
		//if len(origin) == 0 {
		//	return true
		//}
		//u, err := url.Parse(origin[0])
		//if err != nil {
		//	return false
		//}
		//originUrl := "simplegame.now.sh"
		//return strings.EqualFold(u.Host, originUrl)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player := game.NewPlayer(user.Nick, conn)

	gameService.Connection <- player
}