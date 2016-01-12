// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

// Admin API to browse voice messages
package machine

import (
	"net/http"
	"encoding/json"

	"github.com/golang/glog"
)


type AdminWebAPI struct {
	store		*VoiceMessageStorage
	route 		string
}

func CreateAdminWebAPI(store *VoiceMessageStorage, route string, ) (*AdminWebAPI, error) {
	api := &AdminWebAPI { store, route }
	api.addAdminWebAPI(route)
	return api, nil
}



func (api *AdminWebAPI) addAdminWebAPI(route string) {

	http.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		glog.V(3).Infof("Admin API call: %s %s", req.Method, req.URL.String())

		if req.Method != "GET" {
			glog.V(2).Infof("Method %s not supported", req.Method)
			api.sendBadRequest(w, "only GET requests are supported")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("["))

		// Add voice messages
		messages := api.store.FetchAllVoiceMessages()
		first := true
		for callID, msg := range messages {
			if !first {
				w.Write([]byte(","))
			} else {
				first = false
			}
			glog.V(3).Infof("CallID: %s, has message: %s", callID, msg)
			enc := json.NewEncoder(w)
			enc.Encode(msg)
		}
		w.Write([]byte("]"))

	})
}


// Error structure ala google (see Vision API)
type errorWrapper struct {
	Error `json:"error"`
}

type Error struct {
	Code int `json:"code"`
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Details *[]errorDetails `json:"details,omitempty"`
}

type errorDetails struct {
	Type string `json:"@type"`
	Links []struct {
		Description string `json:"description"`
	} `json:"links"`
}

func (api *AdminWebAPI) sendBadRequest(w http.ResponseWriter, message string) {
	api.sendError(w, http.StatusBadRequest, "BAD REQUEST", message)
}

func (api *AdminWebAPI) sendInternalError(w http.ResponseWriter, reason string, message string) {
	api.sendError(w, http.StatusInternalServerError, reason, message)
}

func (api *AdminWebAPI) sendError(w http.ResponseWriter, code int, status string, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	error := Error{
		Code: code,
		Status: status,
		Message: message,
		Details: nil,
	}
	ew := errorWrapper{error}
	enc := json.NewEncoder(w)
	enc.Encode(ew)
}




