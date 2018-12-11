package main

import (
	"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//"SimpleGame/2018_2_Simple_Name/internal/auth"
	//"SimpleGame/2018_2_Simple_Name/internal/db/postgres"
	//leaders "SimpleGame/2018_2_Simple_Name/internal/leaderboard"
	//middle "SimpleGame/2018_2_Simple_Name/internal/middleware"
	//"SimpleGame/2018_2_Simple_Name/internal/profile"
	//"SimpleGame/2018_2_Simple_Name/internal/session"
	"SimpleGame/internal/auth"
	"SimpleGame/internal/db/postgres"
	leaders "SimpleGame/internal/leaderboard"
	middle "SimpleGame/internal/middleware"
	"SimpleGame/internal/profile"
	"SimpleGame/internal/session"
	"SimpleGame/internal/stat"
	"log"
	"net/http" //"github.com/gorilla/mux"
	//"github.com/gorilla/websocket"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	siteMux.HandleFunc("/api/signup", middle.CORSsettings(auth.SignupHandler))
	siteMux.HandleFunc("/api/signin", middle.CORSsettings(auth.SigninHandler))
	siteMux.HandleFunc("/api/profile", middle.CORSsettings(profile.ProfileHandler))
	siteMux.HandleFunc("/api/leaders", middle.CORSsettings(leaders.LeadersHandler))
	siteMux.HandleFunc("/api/islogged", middle.CORSsettings(auth.Islogged))
	siteMux.HandleFunc("/api/logout", middle.CORSsettings(auth.LogOut))
	//siteMux.HandleFunc("/startgame", startGame)
	siteMux.HandleFunc("/api/leaderscount", middle.CORSsettings(leaders.LeadersCount))
	siteMux.HandleFunc("/api/getAvatar", middle.CORSsettings(profile.GetAvatar))
	siteMux.Handle("/api/metrics", promhttp.Handler())

	var HitStat = stat.NewPrometheus()
	siteHandler := middle.AccessLogMiddleware(siteMux, sugar, HitStat)

	port := "8080"

	sugar.Infow("starting server at :" + port)

	//fmt.Println("starting server at :8080")
	if err := http.ListenAndServe(":"+port, siteHandler); err != nil {
		log.Fatalf("cannot listen: %s", err)
	}

}
