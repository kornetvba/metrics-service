package agent

import (
	"bytes"
	"compress/gzip"
)

func CompressData(data *[]byte) ([]byte, error) {
	var dataCompress bytes.Buffer
	gz, err := gzip.NewWriterLevel(&dataCompress, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	gz.Write(*data)

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return dataCompress.Bytes(), nil
}
