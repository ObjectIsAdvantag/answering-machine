package server

import (
	"fmt"
	"time"
	"os"
	"os/signal"

	"net/http"

	"github.com/golang/glog"
)

type Service interface {
	RegisterHandlers()
}

func Run(port string, svc Service, version string, name string) error {

	server := &server{ port, svc, version, name}
	if err := server.start(); err != nil {
		glog.Errorf("Failed to start server: %v\n", err)
		return err
	}

	return nil
}

type server struct {
	port 				string
	service				Service
	serviceVersion 		string
	serviceName 		string
}

func (srv *server) start() error {

	// start http server
	go func() {

		// register health check
		launch := time.Now()
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			glog.V(1).Infof("hit healthcheck endpoint\n")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(w, `{ "name":"%s", "version":"%s", "port":"%s", "started":"%s"}`, srv.serviceName, srv.serviceVersion, srv.port, launch.Format(time.RFC3339))
		})

		// register service endpoints
		srv.service.RegisterHandlers()

		glog.Infof("Listening on http://:%s\n", srv.port)
		if err := http.ListenAndServe(":" + srv.port, nil); err != nil {
			glog.Fatalf("Service died unexpectedly\n",err)
		}
	}()

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	return nil
}

