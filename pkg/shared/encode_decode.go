package shared

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io"
	"log"
)

func Encode[T any](obj T) []byte {

	buffer := bytes.Buffer{}
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(obj)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	result := buffer.Bytes()

	//log.Printf("Object %v ecnoded. Size: %d\n", obj, len(buffer.Bytes()))

	return result
}

func Decode[T any](buffer []byte) *T {
	var result *T
	dec := gob.NewDecoder(bytes.NewReader(buffer))
	err := dec.Decode(&result)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	//log.Printf("Object %v decoded. Size: %d", result, len(buffer))

	return result
}

func Compress(s []byte) []byte {

	buffer := bytes.Buffer{}
	zipped := gzip.NewWriter(&buffer)
	zipped.Write(s)
	zipped.Close()

	return buffer.Bytes()
}

func Decompress(buffer []byte) []byte {

	rdr, _ := gzip.NewReader(bytes.NewReader(buffer))
	data, err := io.ReadAll(rdr)
	if err != nil {
		log.Fatal(err)
	}
	rdr.Close()

	return data
}
