package session

import (
	//"SimpleGame/2018_2_Simple_Name/internal/db/redis"
	//"SimpleGame/2018_2_Simple_Name/internal/generator"
	//"SimpleGame/2018_2_Simple_Name/internal/models"
	"SimpleGame/internal/db/redis"
	"SimpleGame/internal/generator"
	"SimpleGame/internal/models"

	"context"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	"time"
)


var sessManager db.AuthCheckerClient
var ctx context.Context

func OpenConn() (*grpc.ClientConn, error) {
	var grpcConn *grpc.ClientConn
	grpcConn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())

	if err != nil || grpcConn == nil {
		fmt.Println(err.Error())
		return nil, err
	}

	sessManager = db.NewAuthCheckerClient(grpcConn)

	ctx = context.Background()

	return grpcConn, nil
}

func FindSession(r *http.Request) (*db.UserSession, error) {
	val := r.Cookies()

	for i := 0; i < len(val); i++ {

		if val[i].Name == "session_id" {
			sessKey := new(db.SessionKey)

			sessValue := new(db.SessionValue)


			sessKey.ID = val[i].Value
			//fmt.Println(sessKey)
			sessValue, err := sessManager.Get(ctx, sessKey)
			fmt.Println(sessValue)
			if err != nil || sessValue.Email == "" {
				return nil, err
			}

			sess := new(db.UserSession)

			sess.ID = sessKey.ID
			sess.Email = sessValue.Email

			//sess, err := session.Get(redis, val[i].Value)
			
			return sess, nil

		} else {
			continue
		}
	}

	return nil, nil
}

func RmCookie(r *http.Request, w *http.ResponseWriter) (error) {
	UserSession, err := FindSession(r)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if UserSession == nil {
		fmt.Println("UserSession = nil")
		return nil
	}

	sess := new(http.Cookie)
	sess.Name = "session_id"
	sess.Value = UserSession.ID
	sess.Expires = time.Now()

	uSess := new(db.SessionKey)
	uSess.ID = sess.Value

	http.SetCookie(*w, sess)

	_, err = sessManager.Delete(ctx, uSess)

	if err != nil {
		return err
	}

	return nil
}

func SetCookie(user *models.User, w *http.ResponseWriter) (error) {
	uSess := new(db.UserSession)
	sess := new(http.Cookie)
	sess.Name = "session_id"
	sess.Value = generator.UidGen()
	sess.Expires = time.Now().Add(time.Hour*5)

	uSess.ID = sess.Value
	uSess.Email = user.Email

	http.SetCookie(*w, sess)

	_, err := sessManager.Create(ctx, uSess)

	if err != nil {
		return err
	}

	return nil
}

