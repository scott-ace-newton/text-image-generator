package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	le := logrus.NewEntry(l)

	s := NewService(le)
	h := NewHandler(le, s)

	r := mux.NewRouter()
	h.RegisterHandlers(r)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)


	httpServer := &http.Server{
		Addr:         fmt.Sprint("0.0.0.0:8080"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: r,
	}

	go func() {
		log.Infof("Listening on port %v", 8080)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Errorf("HTTP server got shut down error: %v", err)
		}
		sig <- os.Interrupt
	}()
	<-sig
	log.Info("shutting down HTTP server...")
	time.Sleep(2 * time.Second)
	os.Exit(0)


}


