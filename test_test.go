package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	schema = `{
	"type":"record",
	"name":"LongList",
	"aliases": ["LinkedLongs"],
	"fields":[
		{"name":"value","type":"string"}
	]
}`
)

func TestEncodeDecodeWapper(t *testing.T) {
	codec, err := NewFromRegistryMock(newAvroSchema(7, 2, "entity-value", schema))
	var value = map[string]interface{}{
		"value": "xpto",
	}

	binary, err := codec.Encode(value)
	assert.Nil(t, err)

	valueDecoded, err := codec.Decode(binary)
	assert.Nil(t, err)
	assert.Equal(t, value, valueDecoded)
}

func TestMustReturnErrorOnInvalidPayload(t *testing.T) {
	codec, err := NewFromRegistryMock(newAvroSchema(7, 2, "entity-value", schema))
	valueDecoded, err := codec.Decode([]byte(""))
	assert.Error(t, err)
	assert.Nil(t, valueDecoded)
}

func newAvroSchema(id, version int, subject, schema string) string {
	return fmt.Sprintf(`{
			"id":%d,
			"version":%d,
			"subject":"%s",
			"schema":"%s"
		}`,
		id,
		version,
		subject,
		jsonEscape(schema),
	)
}

func jsonEscape(str string) string {
	b, err := json.Marshal(str)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}
