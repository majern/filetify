package shared

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

func Encode[T any](obj *T) []byte {

	buffer := bytes.Buffer{}
	enc := gob.NewEncoder(&buffer)
	HandleErrorWithMsg(enc.Encode(*obj), true, fmt.Sprintf("Failed to encode object '%T'", obj))

	result := buffer.Bytes()

	//log.Printf("Object %v ecnoded. Size: %d\n", obj, len(buffer.Bytes()))

	return result
}

func Decode[T any](buffer []byte) *T {
	var result *T
	dec := gob.NewDecoder(bytes.NewReader(buffer))
	HandleErrorWithMsg(dec.Decode(&result), true, fmt.Sprintf("Failed to decode object '%T'", result))

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
