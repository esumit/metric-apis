package httprqrs

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-TenantType", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"status": 500, "error":"Resource path not mapped"}`))
}

