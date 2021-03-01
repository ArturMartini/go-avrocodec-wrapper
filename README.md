# go-avrocodec-wrapper
go-avrocodec-wrapper is a library to simplify process of load schema from registry, encode and decode data

## Features 
* Load schemas  
* Encode data by latest schema
* Decode data by schema id from message
* Update schema by custom time
* Mock server to simplify unit tests

Usage:
```go
package main
codec "github.com/arturmartini/go-avrocodec-wrapper"

func NewCodec() codec.CodecWrapper {
    //This method load schemas and return a reference for use to encode and decode messages
    //Internally has process to update schemas based by custom time  
    codec, err := codec.NewFromRegistry("http://your-domain/subjects/entity-value/versions", time.Minute*5)
    if err != nil {
        return nil, err
    }
    return codec, nil
}

//This method encode message by latest schema in memory
//You need to parse your domain to map[string]interface{} referenced by avro protocol
func (l Listner) Sender(domain Domain) error {
    mapValue, err := l.yourMarshal(domain)
    if err != nil {
	    return err
	}

    binary, err := l.codec.Encode(value)
    if err != nil {
		return err
	}

    return l.kafka.send(binary)
}

//This method receive data from kafka and decode by schema id in message
//You need parse map[strin]interface{] to your domain
func (l Listner) Receiver(payload []byte) (Domain, error) {
    entity, err := l.codec.Decode(payload)
    if err != nil {
		return Domain{}, err
	}

    domain, err := l.yourUnmarshal(entity)
    if err != nil {
		return Domain{}, err
	}
    
    return domain, nil
}
```