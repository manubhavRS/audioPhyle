package server

import (
	"net/http"
)

func Anyfunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
