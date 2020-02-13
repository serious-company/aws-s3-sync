package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dafiti-group/aws-s3-sync-api/pkg/sync"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
)

type FilesResponse struct {
	Files []string `json:"files"`
}

type Handler struct {
	Log logr.Logger
}

//
func (h *Handler) GetAllFiles(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithName("handler")
	log.Info("Start GetAllfiles")
	vars := mux.Vars(r)
	file := sync.Sync{
		Path: vars["path"],
		Log:  log,
	}
	f, err := file.List()
	if err != nil {
		log.Error(err, "GetAllFiles Failed")
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		log.Info("GetAllFiles Sucess")
		respondJSON(w, http.StatusOK, map[string][]string{"files": f})
	}
}

// Sync ...
func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithName("handler")
	log.Info("Start syncing")
	file := sync.Sync{
		Log: log,
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&file); err != nil {
		log.Error(err, "Sync Decode Failed")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err := file.AwsSync()
	if err != nil {
		log.Error(err, "Sync Failed")
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		log.Info("Sync Sucess")
		respondJSON(w, http.StatusOK, map[string][]string{"message": []string{"Updated"}})
	}
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
