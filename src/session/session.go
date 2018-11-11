package session

import (
	"SimpleGame/models"
	"SimpleGame/generator"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"time"
)

func Create(service models.UserSessionService, user *models.User, w *http.ResponseWriter) (error) {
	sess := new(http.Cookie)
	sess.Name = "session_id"
	sess.Value = generator.UidGen()
	sess.Expires = time.Now().Add(time.Hour*5)

	//sess.HttpOnly = true
	//sess.Secure = true

	http.SetCookie(*w, sess)

	err := service.Create(sess.Value, user.Email)

	fmt.Println("Session value: ", sess.Value)

	if err != nil {
		return err
	}

	return nil
}

func Delete(service models.UserSessionService, sessionId string, w *http.ResponseWriter) (string, error) {
	sess := new(http.Cookie)
	sess.Name = "session_id"
	sess.Value = sessionId
	sess.Expires = time.Now()
	sess.HttpOnly = true

	http.SetCookie(*w, sess)

	err := service.Delete(sessionId)
	if err != nil {
		return "", err
	}

	return "", nil
}

func Get(service models.UserSessionService, sessionId string) (*models.UserSession, error) {
	uSession := new(models.UserSession)
	var err error

	uSession.Id = sessionId
	uSession.Email, err = service.Get(sessionId)
	if err == redis.ErrNil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return uSession, nil
}
