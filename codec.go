package go_avro_codec

import (
	"encoding/json"
	"errors"
	"github.com/linkedin/goavro/v2"
	"io/ioutil"
	"net/http"
	"strconv"
)

type CodecWrapper interface {
	Translate([]byte) (map[string]interface{}, error)
}

type codec struct {
	versions map[int]*goavro.Codec
}

func NewFromRegistry(schemaAddress string) (CodecWrapper, error) {
	var codec codec
	var codecs = map[int]*goavro.Codec{}
	var versions = make([]int, 0)
	var err = getDataFromRegistry(schemaAddress, &versions)
	if err != nil {
		return nil, err
	}

	for _, version := range versions {
		var schemaMap = map[string]interface{}{}
		if err := getDataFromRegistry(schemaAddress+"/"+strconv.Itoa(version), &schemaMap); err != nil {
			return nil, err
		}

		if codecs[version], err = goavro.NewCodec(schemaMap["schema"].(string)); err != nil {
			return nil, err
		}
	}

	codec.versions = codecs
	return &codec, nil
}

func getDataFromRegistry(schema string, rawMap interface{}) error {
	var response, err = http.Get(schema)

	if err != nil || response.StatusCode != http.StatusOK {
		return errors.New("error when getting schema from registry")
	}

	rawData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("cannot extract schema from response")
	}

	err = json.Unmarshal(rawData, &rawMap)
	if err != nil {
		return errors.New("error when unmarshal schema")
	}

	return err
}

func (r codec) Translate(value []byte) (map[string]interface{}, error) {
	var error error
	for _, codec := range r.versions {
		payload, _, err := codec.NativeFromBinary(value)
		if payload != nil {
			return payload.(map[string]interface{}), nil
		} else {
			error = err
		}
	}
	return nil, error
}	