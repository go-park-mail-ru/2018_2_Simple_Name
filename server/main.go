package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

type User struct {
	Uid      string
	Email    string
	Password string
}

type Users map[string]User

var users = make(Users)

func main() {

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/profile", profileHandler)

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("../src/static")))
	http.Handle("/static/", staticHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("../src/index.html")
	w.Write(file)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	if r.Method != http.MethodPost {
		http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
		return
	}

	user, err := getBodyReq(r)
	if err != nil {
		http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
		return
	}

	if user.Email != "" && user.Password != "" {
		if !checkExist(*user) {
			addUser(*user)
			session := new(http.Cookie)
			session.Name = "session_id"
			session.Value = uidGen()
			session.Expires = time.Now().Add(time.Minute)
			session.HttpOnly = true
			http.SetCookie(w, session)
			http.Redirect(w, r, "/profile", http.StatusOK)
		} else {
			http.Redirect(w, r, r.Referer(), http.StatusAlreadyReported)
			return
		}
	} else {
		http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
		return
	}
}

func getBodyReq(r *http.Request) (*User, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // важный пункт!
	if err != nil {
		return nil, err
	}
	user := new(User)
	err = json.Unmarshal(body, user)
	return user, err
}

func uidGen() string {
	uid, _ := exec.Command("uuidgen").Output()
	suid := string(uid[:])
	suid = suid[:len(suid)-1]
	return suid
}

func addUser(user User) bool {
	user.Uid = uidGen()
	users[user.Email] = user
	fmt.Println(time.Now().UTC(), "Added user", user)
	return true
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC(), "Request from", r.URL.String())
	fmt.Println("Method", r.Method)

	if r.Method != http.MethodPost {
		http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
		return
	}
	_, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie)
	if !loggedIn {
		user, err := getBodyReq(r)
		if err != nil {
			http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
			return
		}

		if validateUser(*user) {
			session := new(http.Cookie)
			session.Name = "session_id"
			session.Value = uidGen()
			session.Expires = time.Now().Add(time.Minute)
			session.HttpOnly = true
			http.SetCookie(w, session)
			http.Redirect(w, r, "/profile", http.StatusOK)
		} else {
			http.Redirect(w, r, r.Referer(), http.StatusBadRequest)
			return
		}
	} else {
		http.Redirect(w, r, "/profile", http.StatusOK)
	}
}
func checkExist(user User) bool {
	if _, ok := users[user.Email]; ok {
		return true
	}
	return false
}
func validateUser(user User) bool {
	if mapUser, ok := users[user.Email]; ok {
		if user.Password == mapUser.Password {
			return true
		}
	}
	return false
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request from", r.URL.String())
	fmt.Println("Method", r.Method)
	if r.Method != http.MethodPost {
		http.Redirect(w, r, r.URL.String(), http.StatusBadRequest)
		return
	}
	session, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie)

	if loggedIn {
		fmt.Fprintln(w, `<a href="/logout">logout</a>`)
		fmt.Fprintln(w, "Welcome, "+session.Value)
	} else {
		fmt.Fprintln(w, `<a href="/login">login</a>`)
		fmt.Fprintln(w, "You need to login")
	}
}

// app.post('/signup', function (req, res) {
// 	const password = req.body.password;
// 	const email = req.body.email;
// 	const age = req.body.age;
// 	if (
// 		!password || !email || !age ||
// 		!password.match(/^\S{4,}$/) ||
// 		!email.match(/@/) ||
// 		!(typeof age === 'number' && age > 10 && age < 100)
// 	) {
// 		return res.status(400).json({
// 			error: 'Не валидные данные пользователя'
// 		});
// 	}
// 	if (users[ email ]) {
// 		return res.status(400).json({
// 			error: 'Пользователь уже существует'
// 		});
// 	}

// 	const id = uuid();
// 	const user = {
// 		password,
// 		email,
// 		age,
// 		score: 0
// 	};
// 	ids[ id ] = email;
// 	users[ email ] = user;

// 	res.cookie('sessionid', id, {
// 		expires: new Date(Date.now() + 1000 * 60 * 10)
// 	});
// 	res.status(201).json({
// 		id
// 	});
// });

// app.post('/login', function (req, res) {
// 	const password = req.body.password;
// 	const email = req.body.email;
// 	if (!password || !email) {
// 		return res.status(400).json({
// 			error: 'Не указан E-Mail или пароль'
// 		});
// 	}
// 	if (!users[ email ] || users[ email ].password !== password) {
// 		return res.status(400).json({
// 			error: 'Не верный E-Mail и/или пароль'
// 		});
// 	}

// 	const id = uuid();
// 	ids[ id ] = email;

// 	res.cookie('sessionid', id, {
// 		expires: new Date(Date.now() + 1000 * 60 * 10)
// 	});
// 	res.status(201).json({
// 		id
// 	});
// });

// app.get('/me', function (req, res) {
// 	const id = req.cookies[ 'sessionid' ];
// 	const email = ids[ id ];
// 	if (!email || !users[ email ]) {
// 		return res.status(401).end();
// 	}

// 	users[ email ].score += 1;

// 	res.json(users[ email ]);
// });

// app.get('/users', function (req, res) {
// 	const scorelist = Object.values(users)
// 		.sort((l, r) => r.score - l.score)
// 		.map(user => {
// 			return {
// 				email: user.email,
// 				age: user.age,
// 				score: user.score
// 			};
// 		});

// 	res.json(scorelist);
// });
