package shared

import (
	"fmt"
	"os"
)

type FilePart struct {
	Part        uint64
	TotalParts  uint64
	Path        string
	Length      uint64
	TotalLength uint64
	Data        []byte
}

const chunkSize = 32 * 1024

func LoadFile(path string) ([][]byte, error) {
	f, err := os.Open(path)

	if err != nil {
		HandleError(err, false)
		return nil, err
	}

	defer f.Close()

	stat, err := f.Stat()
	parts := stat.Size() % chunkSize

	result := make([][]byte, parts)

	for i := int64(0); i < parts; i++ {
		//f.Seek...
		part := []byte{}
		_, err = f.ReadAt(part, i*chunkSize)

		if err != nil {
			HandleErrorWithMsg(err, true, fmt.Sprintf("Failed to read chunk#=%d, of file=%v", i, f.Name()))
		}

		result = append(result, part)
	}

	return result, nil
}
