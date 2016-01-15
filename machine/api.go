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

func AddAdminEndpoint(store *VoiceMessageStorage, route string, ) (*AdminWebAPI, error) {
	api := &AdminWebAPI { store, route }
	api.registerAdminWebAPI(route)
	return api, nil
}


func (api *AdminWebAPI) registerAdminWebAPI(route string) {

	http.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		glog.V(3).Infof("Admin API call: %s %s", req.Method, req.URL.String())

		if req.Method != "GET" {
			glog.V(2).Infof("Method %s not supported", req.Method)
			sendBadRequest(w, "only GET requests are supported")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("["))

		// Display voice messages
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





