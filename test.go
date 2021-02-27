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
	idx := strings.LastIndex(address, "/")

	version := address[idx+1:]
	path := address[:idx-1]
	
	s := http.NewServeMux()
	s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[" + version + "]"))
		w.WriteHeader(http.StatusOK)
	})
	
	s.HandleFunc(path + "/" + version, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(schema))
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(s)
	defer server.Close()
	
	return NewFromRegistry(server.URL + path)
}
