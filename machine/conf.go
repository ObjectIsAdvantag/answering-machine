// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

// Configuration API to read effective runtime configuration
package machine

import (
	"net/http"
	"encoding/json"

	"github.com/golang/glog"
)



// Adds an endpoint to display the active runtime configuration of the AnsweringMachine.
// By default, the configuration is accessible at /conf
// Note that the configuration is defined at startup and cannot be changed afterwards
func AddConfEndpoint(machine *AnsweringMachine, route string) {

	// default route
	if route == "" {
		route = "/conf"
	}

	http.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		glog.V(3).Infof("Conf API call: %s %s", req.Method, req.URL.String())

		if req.Method != "GET" {
			glog.V(2).Infof("Method %s not supported", req.Method)
			sendBadRequest(w, "only GET requests are supported")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("{"))

		// write env (mask password)
		publicEnv := *machine.env
		publicEnv.RecorderPassword = "********"
		w.Write([]byte(`"env":`))
		enc := json.NewEncoder(w)
		enc.Encode(publicEnv)

		// write messages
		w.Write([]byte(`, "messages":`))
		enc.Encode(*machine.messages)

		// write routes
		w.Write([]byte(`, "routes":`))
		enc.Encode(*machine.routes)

		w.Write([]byte("}"))
	})
}


