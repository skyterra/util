package primitive

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Gzip 采用gzip进行压缩
func Gzip(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	var buff bytes.Buffer

	gz := gzip.NewWriter(&buff)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// Gunzip 解压gzip数据
func Gunzip(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer r.Close()
	return ioutil.ReadAll(r)
}
