package webserver

import "net/http"

func (app *App) ServerError(w http.ResponseWriter, err error) {
	app.ErrorLogger.Printf("%s\n", err)

	// Probably fine to show a plain 500 error string for this class of errors
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}