package usecases

import (
	"encoding/json"
	"io"
	"log"
	"main/domain"
	"main/pkg"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTPHandler struct {
	Client          *mongo.Client
	HTTPCollection  *mongo.Collection
	HTTPSCollection *mongo.Collection
}

func (h *HTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HTTP requests:\n"))

	results, err := pkg.GetAllHTTPRecords(h.HTTPCollection)

	if err != nil {
		log.Println(err)
		return
	}

	for _, result := range results {
		err := pkg.WriteHTTPRecord(w, result)

		if err != nil {
			log.Println(err)
			return
		}
	}

	w.Write([]byte("HTTPS requests:\n"))

	secondResults, err := pkg.GetAllHTTPSRecords(h.HTTPSCollection)

	if err != nil {
		log.Println(err)
		return
	}

	for _, result := range secondResults {
		err := pkg.WriteHTTPSRecord(w, result)

		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (h *HTTPHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		log.Println(err)
		return
	}

	result, err := pkg.GetCurrentHTTPRecordByID(h.HTTPCollection, id)

	if err == nil {
		err := pkg.WriteHTTPRecord(w, result)

		if err != nil {
			log.Println(err)
			return
		}

		return
	}

	if err != mongo.ErrNoDocuments {
		log.Println(err)
		return
	}

	secondResult, err := pkg.GetCurrentHTTPSRecordByID(h.HTTPSCollection, id)

	if err == mongo.ErrNoDocuments {
		w.Write([]byte("There is no record\n"))
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	err = pkg.WriteHTTPSRecord(w, secondResult)

	if err != nil {
		log.Println(err)
		return
	}
}

func (h *HTTPHandler) RepeatCurrent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		log.Println(err)
		return
	}

	result, err := pkg.GetCurrentHTTPRecordByID(h.HTTPCollection, id)

	if err == mongo.ErrNoDocuments {
		h.repeatHTTPS(w, r, id)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	url := &url.URL{
		Scheme: result.Request.Scheme,
		Host:   result.Request.Host,
		Path:   result.Request.Path,
	}

	req := &http.Request{
		Method: result.Request.Method,
		URL:    url,
		Header: result.Request.Headers,
	}

	client := http.Client{
		Timeout: time.Second * 3,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte("ID запроса: " + strconv.Itoa(result.ID) + "\nЗапрос:\n"))

	json.NewEncoder(w).Encode(result.Request)

	w.Write([]byte("Ответ:\n" + resp.Proto + " " + resp.Status + "\n"))

	for key, value := range resp.Header {
		w.Write([]byte(key + ": " + value[0] + "\n"))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte("\n"))
	w.Write(body)
}

func (h *HTTPHandler) repeatHTTPS(w http.ResponseWriter, r *http.Request, id int) {
	result, err := pkg.GetCurrentHTTPSRecordByID(h.HTTPSCollection, id)

	if err == mongo.ErrNoDocuments {
		w.Write([]byte("There is no record"))
		return
	}

	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte("ID запроса: " + strconv.Itoa(id) + "\nЗапрос:\n" + result.ClientRequest))

	dest_conn.Write([]byte(result.ClientRequest))

	w.Write([]byte("\nОтвет:\n"))

	defer dest_conn.Close()
	io.Copy(w, dest_conn)
}

func (h *HTTPHandler) ScanCurrent(w http.ResponseWriter, r *http.Request) {
	paths, err := pkg.ReadFromFile("dicc.txt")

	if err != nil {
		log.Println(err)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		log.Println(err)
		return
	}

	result, err := pkg.GetCurrentHTTPRecordByID(h.HTTPCollection, id)

	if err == mongo.ErrNoDocuments {
		h.scanHTTPS(w, r, id, paths)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	wg := &sync.WaitGroup{}

	client := http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	w.Write([]byte("starts scanning url " + result.Request.Scheme + "://" + result.Request.Host + "\n"))
	runtime.GOMAXPROCS(5)

	for _, path := range paths {
		wg.Add(1)

		go func() {
			defer wg.Done()

			url := &url.URL{
				Scheme: result.Request.Scheme,
				Host:   result.Request.Host,
				Path:   path,
			}

			req := &http.Request{
				Method: result.Request.Method,
				URL:    url,
				Header: result.Request.Headers,
			}

			resp, err := client.Do(req)

			if err != nil {
				log.Println(err)
				return
			}

			code := resp.StatusCode

			if code != 404 {
				err := json.NewEncoder(w).Encode(domain.ScanResponse{
					Code: code,
					Path: path,
				})

				if err != nil {
					log.Println(err)
					return
				}
			}
		}()
	}

	wg.Wait()
}

func (h *HTTPHandler) scanHTTPS(w http.ResponseWriter, r *http.Request, id int, paths []string) {
	result, err := pkg.GetCurrentHTTPSRecordByID(h.HTTPSCollection, id)

	if err == mongo.ErrNoDocuments {
		w.Write([]byte("There is no record"))
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	wg := &sync.WaitGroup{}
	w.Write([]byte("starts scanning url " + result.Request.Scheme + result.Request.Host + "\n"))
	runtime.GOMAXPROCS(5)

	for _, path := range paths {

		wg.Add(1)

		go func() {
			defer wg.Done()

			resp, err := http.Get("https://" + result.Request.Host + "/" + path)

			if err != nil {
				//log.Println(err)
				return
			}

			code := resp.StatusCode

			if code != 404 {
				err := json.NewEncoder(w).Encode(domain.ScanResponse{
					Code: code,
					Path: path,
				})

				if err != nil {
					//log.Println(err)
					return
				}
			}
		}()
	}

	wg.Wait()
}
