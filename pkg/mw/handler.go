package mw

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HttpRqRsMiddleware(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		err := f(w, r)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"status": 500, "error":"Resource path not mapped"}`))
		}
	})
}

