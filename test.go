package main

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

type AvroSchema struct {
	Subject string `json:"subject"`
	Id      int    `json:"id"`
	Version int    `json:"version"`
	Schema  string `json:"schema"`
}

func NewFromRegistryMock(schema string) (CodecWrapper, error) {
	avroSchema := AvroSchema{}
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
