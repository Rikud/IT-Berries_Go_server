package controllers

import (
	"IT-Berries/auth/models"
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

type LoginHandle struct {}
type RegistrationHandle struct {}
type MeHandle struct {}
type ScoreboardHandle struct {}

func (handle LoginHandle)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(notSafeMethod)
		//TODO write coorect redirect fo this request;
		return
	}
	fmt.Fprintln(w, "Trying to login user", r.URL.String())
	fmt.Println("Trying to login user")
	loginForm := models.LoginForm{r.FormValue("email"), r.FormValue("password")}
	user := loginForm.
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