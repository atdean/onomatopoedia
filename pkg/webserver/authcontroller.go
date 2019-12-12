package webserver

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"

	"github.com/atdean/onomatopoedia/pkg/repositories"
)

type AuthController struct {
	App *App
}

func NewAuthController(app *App) *AuthController {
	return &AuthController{
		App: app,
	}
}

// GET /signup

// POST /signup
// Code to hash password, for later use
// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)

// GET /login
func (ctrl *AuthController) GetLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/auth/login.html",
	))
	if err := renderer.ExecuteTemplate(w, "base", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<div>Sorry, an internal server error occurred.</div>"))
	}
}

// POST /login
func (ctrl *AuthController) PostLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if len(username) <= 0 || len(password) <= 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<div>Username or password not provided.</div>"))
	} else {
		repo := repositories.NewUserRepository(ctrl.App.SqlPool)

		user, err := repo.GetByUsername(username)
		if err != nil || user == nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("<div>Username not found.</div>"))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("<div>Sorry, an internal server error occurred.</div>"))
			}
		} else {
			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("<div>Incorrect password."))
			} else {
				// Create a random-number session token
				sessionToken := uuid.NewV4().String()
				fmt.Println(sessionToken)

				_, err := ctrl.App.RedisConn.Do("SETEX", sessionToken, "120", username)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name: "session_token",
					Value: sessionToken,
					Expires: time.Now().Add(120 * time.Second),
				})

				http.Redirect(w, r, "/", http.StatusOK)
			}
		}
	}
}