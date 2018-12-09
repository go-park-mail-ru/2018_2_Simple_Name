package main

import (
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//"SimpleGame/2018_2_Simple_Name/internal/auth"
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//leaders "SimpleGame/2018_2_Simple_Name/internal/leaderboard"
	//middle "SimpleGame/2018_2_Simple_Name/internal/middleware"
	//"SimpleGame/2018_2_Simple_Name/internal/profile"
	//"SimpleGame/2018_2_Simple_Name/internal/session"
	"SimpleGame/internal/session"
	"SimpleGame/internal/auth"
	"SimpleGame/internal/db/postgres"
	leaders "SimpleGame/internal/leaderboard"
	middle "SimpleGame/internal/middleware"
	"SimpleGame/internal/profile"

	"log"
	"net/http"

	//"github.com/gorilla/mux"
	//"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

//var postgres models.UserService = &db.PostgresUserService{}
//var redis models.UserSessionService = &session.RedisSessionService{}
var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

//var gameService = game.NewGame()


//var chatService = chat.NewChat()


func main() {

	defer logger.Sync()

	//err := postgres.InitService()

	err := db.OpenConn()

	if err != nil {
		return
	}

	grpcConn, err := session.OpenConn()

	if err != nil || grpcConn == nil {
		return
	}

	defer grpcConn.Close()


	//_, err = redis.InitService()

	//if err != nil {
	//	//logging.ErrorLog("Failed open redis", err, sugar)
	//	sugar.Errorw("Failed open redis",
	//		"error", err,
	//		"time", strconv.Itoa(time.Now().Hour())+":"+strconv.Itoa(time.Now().Minute()))
	//	////fmt.Println(err.Error())
	//	return
	//}

	//defer obj.Close() // Не будет работать

	//go gameService.Run()

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/signup", middle.CORSsettings(auth.SignupHandler))
	siteMux.HandleFunc("/signin", middle.CORSsettings(auth.SigninHandler))
	siteMux.HandleFunc("/profile", middle.CORSsettings(profile.ProfileHandler))
	siteMux.HandleFunc("/leaders", middle.CORSsettings(leaders.LeadersHandler))
	siteMux.HandleFunc("/islogged", middle.CORSsettings(auth.Islogged))
	siteMux.HandleFunc("/logout", middle.CORSsettings(auth.LogOut))
	//siteMux.HandleFunc("/startgame", startGame)
	siteMux.HandleFunc("/leaderscount", middle.CORSsettings(leaders.LeadersCount))
	siteMux.HandleFunc("/getAvatar", middle.CORSsettings(profile.GetAvatar))

	siteHandler := middle.AccessLogMiddleware(siteMux, sugar)

	port := "8080"

	sugar.Infow("starting server at :" + port)

	//fmt.Println("starting server at :8080")
	if err := http.ListenAndServe(":"+port, siteHandler); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}

}
