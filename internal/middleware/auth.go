package middleware

import (
	//"SimpleGame/2018_2_Simple_Name/internal/session"
	"SimpleGame/internal/session"
	"net/http"
)

func IsLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.FindSession(r)

		if err != nil {
			//fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if sess != nil {
			//w.WriteHeader(http.StatusOK)
			//return

			next.ServeHTTP(w, r)
		} else {
			//w.WriteHeader(http.StatusBadRequest)
			next.ServeHTTP(w, r)
			return
		}
	})

}