package pkg

import (
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func ParseHTTPRequest(r *http.Request) (bson.D, error) {
	method := r.Method
	path := r.URL.Path
	scheme := r.URL.Scheme
	host := r.URL.Host
	getParams := r.URL.Query()
	headers := r.Header
	cookies := r.Cookies()
	postParams := r.PostForm

	bsonGetParams := bson.D{}

	for key, value := range getParams {
		bsonGetParams = append(bsonGetParams, bson.E{key, value})
	}

	bsonHeaders := bson.D{}

	for key, value := range headers {
		bsonHeaders = append(bsonHeaders, bson.E{key, value})
	}

	bsonCookies := bson.D{}

	for _, value := range cookies {
		bsonCookies = append(bsonCookies, bson.E{value.Name, value.Value})
	}

	bsonPostParams := bson.D{}

	for key, value := range postParams {
		bsonPostParams = append(bsonPostParams, bson.E{key, value})
	}

	return bson.D{
		{"method", method},
		{"scheme", scheme},
		{"host", host},
		{"path", path},
		{"get_params", bsonGetParams},
		{"headers", bsonHeaders},
		{"cookies", bsonCookies},
		{"post_params", bsonPostParams},
	}, nil
}

func ParseHTTPAnswer(resp *http.Response) (bson.D, []byte, error) {
	code := resp.StatusCode
	message := http.StatusText(code)
	headers := resp.Header

	bsonHeaders := bson.D{}

	for key, value := range headers {
		bsonHeaders = append(bsonHeaders, bson.E{key, value})
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return bson.D{}, nil, err
	}

	return bson.D{
		{"code", code},
		{"message", message},
		{"headers", bsonHeaders},
		{"body", string(body)},
	}, body, nil
}
