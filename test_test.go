package go_avro_codec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNla(t *testing.T) {
	_, err := NewFromRegistryMock("/subjects/payments-value/versions/7", "{\"subject\":\"payments-value\",\"version\":2,\"id\":7,\"schema\":\"{\\\"type\\\":\\\"record\\\",\\\"name\\\":\\\"LongList\\\",\\\"aliases\\\":[\\\"LinkedLongs\\\"],\\\"fields\\\":[{\\\"name\\\":\\\"value\\\",\\\"type\\\":\\\"long\\\"}]}\"}")
	assert.Nil(t, err)
}
