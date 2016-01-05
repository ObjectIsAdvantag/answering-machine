package service

import (
	"fmt"
	"time"
	"os"
	"os/signal"

	"net/http"

	"github.com/golang/glog"
)


func Run(apiKey string, port string, version string) error {

	service := &Service{ apiKey, port, version}
	if err := service.Start(); err != nil {
		glog.Errorf("Failed to start service: %v\n", err)
		return err
	}

	return nil
}

type Service struct {
	apiKey      string
	port 		string
	version 	string
}

func (svc *Service) Start() error {

	// start http server
	go func() {

		// register health check
		start := time.Now()
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			glog.V(1).Infof("hit healthcheck endpoint\n")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(w, `{ "name":"%s", "version":"%s", "port":"%s", "started":"%s"}`, "Answering Machine", svc.version, svc.port, start.Format(time.RFC3339))
		})

		// add a default route if the proxy is not registered on /
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			glog.V(1).Infof("hit default endpoint\n")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
		})

		glog.Infof("Listening on http://:%s\n", svc.port)
		if err := http.ListenAndServe(":" + svc.port, nil); err != nil {
			glog.Fatalf("Service died unexpectedly\n",err)
		}
	}()

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	return nil
}

