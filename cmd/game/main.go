package main

import (
	"SimpleGame/internal/game"
	//"SimpleGame/session"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
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
	grpcConn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())

	if err != nil || grpcConn == nil {
		fmt.Println(err.Error())
		return
	}

	mux := http.NewServeMux()

	go gameService.Run()

	mux.HandleFunc("/startgame", startGame)

	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}

	fmt.Println("Starting sess server at :8082")
}



func startGame(w http.ResponseWriter, r *http.Request) {
	//sugar.Info("Startgame signal from user")
	sess, err := findSession(r)
	if err != nil {
		fmt.Println("Failed get session", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := postgres.GetUser(sess.Email)

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