package controllers

import (
	"IT-Berries_Go_server/auth/encoder"
	"IT-Berries_Go_server/auth/models"
	"IT-Berries_Go_server/auth/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

var notSafeMethod = []byte(`
<html>
	<body>
	[Error]: not safe method, use POST for send personal data
	</body>
</html>
`)

type JSONError struct {
	Err string `json:"error, string"`
	Code int
}

var Handlers = map[string]Handler{}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type LoginHandle struct{}
type RegistrationHandle struct{}
type MeHandle struct{}
type ScoreboardHandle struct{}
type LogOut struct {}

func (handle LogOut) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Log out failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	log.Println("Trying to log out user.")
	cookie, err := r.Cookie(services.CookieName)
	if err != nil || cookie.Value == "" {
		log.Println("The user isn't authorized!")
		errorResponce(w, "The user isn't authorized!", http.StatusUnauthorized)
		return
	}
	services.DeleteSession(cookie, w)
	w.WriteHeader(http.StatusOK)
}

func (handle LoginHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Login failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	if !checkMethod(w, r) {return}
	log.Println("Trying to login user")
	loginForm := models.LoginForm{r.FormValue("email"), r.FormValue("password")}
	user := services.FindUserByEmail(loginForm.GetEmail())
	if user == nil {
		errorResponce(w, "Wrong email or password!", http.StatusBadRequest)
		log.Println("Couldn't find a user with this email")
		return
	}
	if encoder.ComparePasswords(user.GetPassword(), loginForm.GetPassword()) {
		log.Println("Login success")
		err := services.NewSession(w, user.GetId())
		if err != nil {
			panic("Error while trying to generate session id!")
		}
		user.SetScore(services.GetBestScoreForUserById(user.GetId()))
		result, _ := json.Marshal(user)
		fmt.Fprint(w, string(result))
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Wrong password")
		errorResponce(w, "Wrong email or password!", http.StatusBadRequest)
	}
}

func (handle RegistrationHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Registration failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	if !checkMethod(w, r) {return}
	log.Println("Trying to register new user")
	username := r.FormValue("username")
	if username == "" {
		log.Println("Empty username")
		errorResponce(w, "Specify a correct login!", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	check, err := regexp.MatchString("(.*)@(.*)", email)
	if err != nil {
		log.Println("Email matching error")
		errorResponce(w, "Email matching error!", http.StatusBadRequest)
		return
	}
	if email == "" || !check {
		log.Println("Empty email")
		errorResponce(w, "Specify a valid e-mail!", http.StatusBadRequest)
		return
	}
	if services.FindUserByEmail(email) != nil {
		log.Println("User with this email already exists!")
		errorResponce(w, "User with this email already exists!", http.StatusConflict)
		return
	}
	password := r.FormValue("password")
	if password == "" || len(password) < 4 {
		log.Println("Wrong password length")
		errorResponce(w, "The password field must contain more than 4 characters!", http.StatusBadRequest)
		return
	}
	repPassword := r.FormValue("password_repeat")
	if repPassword == "" || repPassword != password {
		log.Println("Passwords do not match")
		errorResponce(w, "Repeat password correctly!", http.StatusBadRequest)
		return
	}
	r.ParseMultipartForm(32 << 20)
	avatarFile, avatarHeader, err := r.FormFile("avatar")
	if err != nil{
		log.Println("Error reading avatar file!", err)
	} else {
		defer avatarFile.Close()
	}

	avatarName := username + "_avatar"
	if avatarFile != nil && avatarHeader.Filename != "" {
		avatarSave, err := os.OpenFile("/home/ivan/Park/semest2/front/DZ2/2018_1_IT-Berries/" + avatarName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Save avatar error!")
			return
		}
		defer avatarSave.Close()
		io.Copy(avatarSave, avatarFile)
	}
	user := new(models.User)
	user.SetEmail(email)
	user.SetPassword(encoder.HashPassword(password))
	user.SetUsername(username)
	user.SetAvatar(avatarName)
	services.SaveUser(*user)
	log.Println("Register success!")
	result, _ := json.Marshal(user)
	fmt.Fprint(w, string(result))
	w.WriteHeader(http.StatusCreated)
}

func (handle MeHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handle.authentication(w, r)
	} else {
		handle.profileChange(w, r)
	}
}

func (handll MeHandle) authentication(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Registration failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	log.Println("Trying to authenticate user.")
	cookie, err := r.Cookie(services.CookieName)
	if err != nil || cookie.Value == "" {
		log.Println("The user isn't authorized!")
		errorResponce(w, "The user isn't authorized!", http.StatusUnauthorized)
		return
	}
	user := services.GetUserBySessionId(cookie.Value)
	if user == nil {
		log.Println("Can't find such user")
		errorResponce(w, "The user isn't authorized!", http.StatusUnauthorized)
		return
	}
	user.SetScore(services.GetBestScoreForUserById(user.GetId()))
	result, _ := json.Marshal(user)
	fmt.Fprint(w, string(result))
	w.WriteHeader(http.StatusOK)
}

func (handle MeHandle) profileChange(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to change user profile data.")
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
	Handlers["/api/logout"] = LogOut{}
}

func checkMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		w.Write(notSafeMethod)
		//TODO write coorect redirect fo this request;
		return false
	}
	return true
}

func AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, accept, authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func errorResponce(w http.ResponseWriter,  errorMessage string, errorStatus int) {
	error := &JSONError{errorMessage, errorStatus}
	result, _ := json.Marshal(error)
	http.Error(w, string(result), errorStatus)
	w.WriteHeader(errorStatus)
}