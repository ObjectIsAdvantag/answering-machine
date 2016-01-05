package server

import (
	"fmt"
	"time"
	"os"
	"os/signal"

	"net/http"

	"github.com/golang/glog"
	"github.com/ObjectIsAdvantag/answering-machine/machine"
)


func Run(port string, version string) error {

	service := &Server{ port, version}
	if err := service.Start(); err != nil {
		glog.Errorf("Failed to start server: %v\n", err)
		return err
	}

	return nil
}

type Server struct {
	port 		string
	version 	string
}

func (svc *Server) Start() error {

	// start http server
	go func() {

		// register health check
		start := time.Now()
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			glog.V(1).Infof("hit healthcheck endpoint\n")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(w, `{ "name":"%s", "version":"%s", "port":"%s", "started":"%s"}`, "Answering Machine", svc.version, svc.port, start.Format(time.RFC3339))
		})

		// register the TropoApplication
		app := machine.NewAnsweringMachine()
		app.RegisterHandlers()

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

