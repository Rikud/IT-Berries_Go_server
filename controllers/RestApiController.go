package controllers

import (
	"IT-Berries/auth/encoder"
	"IT-Berries/auth/models"
	"IT-Berries/auth/services"
	"fmt"
	"net/http"
)

var notSafeMethod = []byte(`
<html>
	<body>
	[Error]: not safe method, use POST for send personal data
	</body>
</html>
`)

var Handlers = map[string]Handler{}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type LoginHandle struct{}
type RegistrationHandle struct{}
type MeHandle struct{}
type ScoreboardHandle struct{}

func (handle LoginHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Fprintln(w, "Login failed.", rec)
		}
	}()
	if r.Method != http.MethodPost {
		w.Write(notSafeMethod)
		//TODO write coorect redirect fo this request;
		return
	}
	fmt.Fprintln(w, "Trying to login user", r.URL.String())
	fmt.Println("Trying to login user")
	loginForm := models.LoginForm{r.FormValue("email"), r.FormValue("password")}
	user := services.FindUserByEmail(loginForm.GetEmail())
	if user == nil {
		panic("Couldn't find a user with this email")
	}
	password, err := encoder.HashPassword(loginForm.GetPassword())
	if err != nil {
		panic("Hashing error")
	}
	if encoder.ComparePasswords(user.GetPassword(), password) {
		fmt.Fprintln(w, "Login success")
	} else {
		panic("Wrong password")
	}
}

func (handle RegistrationHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Trying to register new user", r.URL.String())
}

func (handle MeHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Trying to authorize  new user", r.URL.String())
}

func (handle ScoreboardHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Leaderboard", r.URL.String())
}

func init() {
	fmt.Println("init in RestApiController.go")
	Handlers["/api/login"] = LoginHandle{}
	Handlers["/api/registration"] = RegistrationHandle{}
	Handlers["/api/me"] = MeHandle{}
	Handlers["/api/users/scoreboard"] = ScoreboardHandle{}
}
