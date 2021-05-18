package go_avrocodec_wrapper

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/linkedin/goavro/v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type CodecWrapper interface {
	Encode(map[string]interface{}) ([]byte, error)
	Decode([]byte) (map[string]interface{}, error)
}

type codec struct {
	codecs     map[int]*goavro.Codec
	latest     int
	address    string
	timeUpdate time.Duration
}

func NewFromRegistry(schemaAddress string, timeUpdate time.Duration) (CodecWrapper, error) {
	var codec = codec{
		address:    schemaAddress,
		codecs:     map[int]*goavro.Codec{},
		timeUpdate: timeUpdate,
	}
	var versions, err = codec.getVersionsFromRegistry()
	if err != nil {
		return nil, err
	}

	err = codec.getSchemaByVersionFromRegistry(versions)
	go codec.update()
	return &codec, err
}

func (r *codec) update() {
	for true {
		<-time.After(r.timeUpdate)
		var versions, err = r.getVersionsFromRegistry()
		if err != nil {
			continue
		}

		if len(versions) == len(r.codecs) {
			continue
		}

		err = r.getSchemaByVersionFromRegistry(versions)
	}
}

func (r *codec) getVersionsFromRegistry() ([]int, error) {
	var versions = make([]int, 0)
	var err = getDataFromRegistry(r.address, &versions)
	return versions, err
}

func (r *codec) getSchemaByVersionFromRegistry(versions []int) error {
	for idx, version := range versions {
		var schemaMap = map[string]interface{}{}
		if err := getDataFromRegistry(r.address+"/"+strconv.Itoa(version), &schemaMap); err != nil {
			return err
		}

		c, err := goavro.NewCodec(schemaMap["schema"].(string))
		if err != nil {
			return err
		}

		schemaId := int(schemaMap["id"].(float64))

		r.codecs[schemaId] = c

		if idx+1 == len(versions) {
			r.latest = schemaId
		}
	}
	return nil
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

func (r *codec) Encode(value map[string]interface{}) ([]byte, error) {
	var payload = make([]byte, 0)
	var binaryValue []byte
	var binarySchemaId = make([]byte, 4)

	binary.BigEndian.PutUint32(binarySchemaId, uint32(r.latest))

	binaryPayload, err := r.codecs[r.latest].BinaryFromNative(payload, value)

	binaryValue = append(binaryValue, byte(0))
	binaryValue = append(binaryValue, binarySchemaId...)
	binaryValue = append(binaryValue, binaryPayload...)

	return binaryValue, err
}

func (r *codec) Decode(value []byte) (map[string]interface{}, error) {
	var error error
	for _, codec := range r.codecs {
		payload, _, err := codec.NativeFromBinary(value[5:])
		if err == nil {
			if payload != nil {
				return payload.(map[string]interface{}), nil
			} else {
				error = err
			}
		}

		payload, _, err = codec.NativeFromBinary(value)
		if err == nil {
			if payload != nil {
				return payload.(map[string]interface{}), nil
			} else {
				error = err
			}
		}

	}
	return nil, error
}
