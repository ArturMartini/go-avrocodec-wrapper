package go_avrocodec_wrapper

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

type codecTest struct {
	codec
}

type avroSchema struct {
	Subject string `json:"subject"`
	Id      int    `json:"id"`
	Version int    `json:"version"`
	Schema  string `json:"schema"`
}

// NewFromRegistryMock create a new CodecWrapper from a schema string
// the format expected is:
// {
//    "id": 1,
//    "subject": "entity-value",
//    "version": 2,
//    "schema": "{}",
// }
func NewFromRegistryMock(schema string) (CodecWrapper, error) {
	avroSchema := avroSchema{}
	err := json.Unmarshal([]byte(schema), &avroSchema)
	if err != nil {
		return nil, errors.New("invalid schema structure")
	}

	path := "/" + avroSchema.Subject
	id := strconv.Itoa(avroSchema.Id)

	s := http.NewServeMux()
	s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[" + id + "]"))
	})

	s.HandleFunc(path+"/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(schema))
	})

	server := httptest.NewServer(s)
	defer server.Close()

	return NewFromRegistry(server.URL+path, time.Minute*1)
}
