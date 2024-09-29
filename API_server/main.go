package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"main/usecases"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))

	if err != nil {
		fmt.Println(err)
		return
	}

	collection := client.Database("queriesdb").Collection("HTTPqueries")
	collectionS := client.Database("queriesdb").Collection("HTTPSqueries")

	handler := &usecases.HTTPHandler{
		Client:          client,
		HTTPCollection:  collection,
		HTTPSCollection: collectionS,
	}

	router := mux.NewRouter()
	router.HandleFunc("/requests", handler.GetAll)
	router.HandleFunc("/requests/{id:[0-9]+}", handler.GetCurrent)
	router.HandleFunc("/repeat/{id:[0-9]+}", handler.RepeatCurrent)
	router.HandleFunc("/scan/{id:[0-9]+}", handler.ScanCurrent)

	server := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		TLSConfig: &tls.Config{
			GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert, err := tls.LoadX509KeyPair("./certs/certbundlе.pem", "./certs/cert.key")

				if err != nil {
					return nil, err
				}

				return &cert, nil
			},
		},
	}

	fmt.Println("starts serving at port 8000")
	log.Fatal(server.ListenAndServe())
}
