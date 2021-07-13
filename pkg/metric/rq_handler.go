package metric

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type metricApiRqHandler struct {
	data DataManager
}

func NewMetricApiRqHandler(data DataManager) *metricApiRqHandler {
	return &metricApiRqHandler{data}
}

func (h *metricApiRqHandler) Save(w http.ResponseWriter, r *http.Request) error {
	var mc MetricCreateRq
	vars := mux.Vars(r)
	key := vars["key"]
	mc.Key = key
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mc)
	if err != nil {
		log.Println(time.Now().String(), ":", "JSON Error:", err)
		w.Write([]byte(`{"status": 400, "error":"Bad Request"}`))
		return nil
	}
	log.Println("Save Value:" , mc.Value)
	log.Println("Save Key:" , mc.Key)
	
	
	mcr, err := h.data.Save(r.Context(), &mc)

	w.Header().Set("Content-Type", "application/json")
	if mcr != nil {
		pcrJson, _ := json.MarshalIndent(mcr, "", "    ")
		w.Write([]byte(pcrJson))
	} else {
		w.Write([]byte(err.Error()))
	}
	
	return nil
}

func (h *metricApiRqHandler) Get(w http.ResponseWriter, r *http.Request) error {
	
	vars := mux.Vars(r)
	key := vars["key"]
	log.Println(key)
	pcr, err := h.data.Get(r.Context(), key)
	
	w.Header().Set("Content-Type", "application/json")
	if pcr != nil {
		pcrJson, _ := json.MarshalIndent(pcr, "", "    ")
		w.Write([]byte(pcrJson))
	} else {
		w.Write([]byte(err.Error()))
	}

	return nil
}

