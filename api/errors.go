package api

import "net/http"

func (a *APIServer) errorResponse(w http.ResponseWriter, _ *http.Request, status int, err error) {
	w.Header().Set("X-App-Error", err.Error())
	http.Error(w, err.Error(), status)
}
