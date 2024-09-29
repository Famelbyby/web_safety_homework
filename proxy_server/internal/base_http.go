package homework1

import (
	"context"
	"fmt"
	"main/pkg"
	"net"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RequestsHandler struct {
	Client *mongo.Client
	ID     int
	Mutex  *sync.Mutex
}

func (h *RequestsHandler) HandleHTTP(w http.ResponseWriter, r *http.Request) error {
	header := make(map[string][]string)

	for name, value := range r.Header {
		if name != "Proxy-Connection" {
			header[name] = value
		}
	}

	req := &http.Request{
		Method: r.Method,
		URL:    r.URL,
		Header: header,
		Host:   r.Host,
	}

	client := http.Client{
		Timeout: time.Second * 3,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	h.Mutex.Lock()
	id := h.ID
	h.Mutex.Unlock()

	collection := h.Client.Database("queriesdb").Collection("HTTPqueries")
	myctx, can := context.WithTimeout(context.Background(), 10*time.Second)
	defer can()

	requestData, err := pkg.ParseHTTPRequest(req)

	if err != nil {
		return err
	}

	answerData, body, err := pkg.ParseHTTPAnswer(resp)

	if err != nil {
		return err
	}

	_, err = collection.InsertOne(myctx, bson.D{{"request", requestData}, {"answer", answerData}, {"_id", id}})

	if err != nil {
		return err
	}

	h.Mutex.Lock()
	h.ID++
	h.Mutex.Unlock()

	w.Write([]byte(resp.Proto + " " + resp.Status + "\n"))

	for key, value := range resp.Header {
		w.Write([]byte(key + ": " + value[0] + "\n"))
	}

	w.Write([]byte("\n"))
	w.Write(body)

	return nil
}

func (h *RequestsHandler) HandleHTTPS(w http.ResponseWriter, r *http.Request) error {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)

	if err != nil {
		return err
	}

	w.Header().Add("Transfer-Encoding", "gzip, chunked")
	w.Header().Add("Content-Encoding", "gzip")
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)

	if !ok {
		return fmt.Errorf("error with hijacker")
	}

	client_conn, _, err := hijacker.Hijack()

	if err != nil {
		return err
	}

	collection := h.Client.Database("queriesdb").Collection("HTTPSqueries")
	myctx, can := context.WithTimeout(context.Background(), 10*time.Second)
	defer can()

	var clientRequestData string
	var answerData string

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go pkg.Transfer(dest_conn, client_conn, wg, &clientRequestData)
	go pkg.Transfer(client_conn, dest_conn, wg, &answerData)

	wg.Wait()

	fmt.Println([]byte(clientRequestData), []byte(answerData))

	bsonRequest, err := pkg.ParseHTTPRequest(r)

	if err != nil {
		return err
	}

	h.Mutex.Lock()
	id := h.ID
	h.Mutex.Unlock()

	_, err = collection.InsertOne(myctx, bson.D{
		{"request", bsonRequest},
		{"client_request", clientRequestData},
		{"answer", answerData},
		{"_id", id},
	})

	h.Mutex.Lock()
	h.ID++
	h.Mutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}
