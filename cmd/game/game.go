package main

import (

	"SimpleGame/internal/db/postgres"
	"SimpleGame/internal/game" //"SimpleGame/session"
	"SimpleGame/internal/session"
	"fmt"
	"log" //	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp" 

)


var (
	
	gameService = game.NewGame()
)
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

	mux.HandleFunc("/api/startgame", startGame)
	mux.Handle("/api/metrics", promhttp.Handler())

	fmt.Println("Starting game server at :8082")

	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}

}

func startGame(w http.ResponseWriter, r *http.Request) {
	//sugar.Info("Startgame signal from user")
	sess, err := session.SessionObj.FindSession(r)
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
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	q := r.URL.Query()
	SingleFlag := false
	if q.Get("single") == "true"{
		SingleFlag = true
	}
	fmt.Println(SingleFlag)
	player := game.NewPlayer(user.Nick, conn,SingleFlag)

	gameService.Connection <- player
}
