package go_avrocodec_wrapper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeWapper(t *testing.T) {
	codec, err := NewFromRegistryMock("{\"subject\":\"entity-value\",\"version\":2,\"id\":7,\"schema\":\"{\\\"type\\\":\\\"record\\\",\\\"name\\\":\\\"LongList\\\",\\\"aliases\\\":[\\\"LinkedLongs\\\"],\\\"fields\\\":[{\\\"name\\\":\\\"value\\\",\\\"type\\\":\\\"string\\\"}]}\"}")
	var value = map[string]interface{}{
		"value": "xpto",
	}

	binary, err := codec.Encode(value)
	assert.Nil(t, err)

	valueDecoded, err := codec.Decode(binary)
	assert.Nil(t, err)
	assert.Equal(t, value, valueDecoded)
}