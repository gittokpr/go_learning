package main

import (
	"context"
	"e2/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	//creating a new logger
	log := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewProduct(log)

	//creating a new servemux for better control
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", hh.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", hh.UpdateProduct)
	putRouter.Use(hh.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", hh.AddProduct)
	postRouter.Use(hh.MiddlewareProductValidation)

	//sm.Handle("/", hh)

	//creating server for manual control
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// to not block this code
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// make an os chanel
	//signal.Notify will broadcast a msg to sigChan when any of the mentioned os commands are called.
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// reading from a channel will block until a msg is available
	sig := <-sigChan
	log.Println("Received Terminate, graceful shutdown", sig)

	//graceful shutdown once msg is consumed ( waiting 30 seconds before shutdown )
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
