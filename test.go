package go_avro_codec

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

type codecTest struct {
	codec
}

func NewFromRegistryMock(address, schema string) (CodecWrapper, error) {
	id := strings.LastIndex(schema, "/")
	
	path := address[:id-2]
	
	s := http.NewServeMux()
	s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[" + address[id:] + "]"))
		w.WriteHeader(http.StatusOK)
	})
	
	s.HandleFunc(path + "/" + address[id:], func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(schema))
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(s)
	defer server.Close()
	
	return NewFromRegistry(schema)
}
