package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	homework1 "main/internal"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	RequestHandler *homework1.RequestsHandler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method == http.MethodConnect {
		err = h.RequestHandler.HandleHTTPS(w, r)
	} else {
		err = h.RequestHandler.HandleHTTP(w, r)
	}

	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	myctx, can := context.WithTimeout(context.Background(), 10*time.Second)
	defer can()

	cursor, err := client.Database("queriesdb").Collection("HTTPqueries").Find(myctx, bson.D{})

	if err != nil {
		log.Println(err)
		return
	}

	var HTTPresults []interface{}

	if err = cursor.All(myctx, &HTTPresults); err != nil {
		log.Println(err)
		return
	}

	cursor, err = client.Database("queriesdb").Collection("HTTPSqueries").Find(myctx, bson.D{})

	if err != nil {
		log.Println(err)
		return
	}

	var HTTPSresults []interface{}

	if err = cursor.All(myctx, &HTTPSresults); err != nil {
		log.Println(err)
		return
	}

	handler := &Handler{
		RequestHandler: &homework1.RequestsHandler{
			ID:     len(HTTPresults) + len(HTTPSresults) + 1,
			Mutex:  &sync.Mutex{},
			Client: client,
		},
	}

	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error

			if r.Method == http.MethodConnect {
				err = handler.RequestHandler.HandleHTTPS(w, r)
			} else {
				err = handler.RequestHandler.HandleHTTP(w, r)
			}

			if err != nil {
				log.Println(err)
			}
		}),
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

	fmt.Println("starts serving at port 8080")

	log.Fatal(server.ListenAndServe())
}
