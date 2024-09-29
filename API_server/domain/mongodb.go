package domain

type Request struct {
	Method     string              `json:"method" bson:"method"`
	Scheme     string              `json:"scheme,-" bson:"scheme"`
	Host       string              `json:"host,-" bson:"host"`
	Path       string              `json:"path" bson:"path"`
	GetParams  map[string][]string `json:"get_params" bson:"get_params"`
	Headers    map[string][]string `json:"headers" bson:"headers"`
	Cookies    map[string]string   `json:"cookies" bson:"cookies"`
	PostParams map[string][]string `json:"post_params" bson:"post_params"`
}

type Answer struct {
	Code    int                 `json:"code" bson:"code"`
	Message string              `json:"message" bson:"message"`
	Headers map[string][]string `json:"headers" bson:"headers"`
	Body    string              `json:"body" bson:"body"`
}

type HTTPEntity struct {
	Request Request `bson:"request"`
	Answer  Answer  `bson:"answer"`
	ID      int     `bson:"_id"`
}

type HTTPSEntity struct {
	Request       Request `json:"request" bson:"request"`
	ClientRequest string  `json:"client_request" bson:"client_request"`
	AnswerData    string  `json:"answer_data" bson:"answer_data"`
	ID            int     `bson:"_id"`
}

type ScanResponse struct {
	Code int    `json:"code"`
	Path string `json:"file"`
}
