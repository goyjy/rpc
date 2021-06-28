package comF

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Compress(data string) string {
	var in bytes.Buffer
	gz := gzip.NewWriter(&in)
	_, _ = gz.Write([]byte(data))
	_ = gz.Flush()
	_ = gz.Close()
	res := in.Bytes()
	return string(res)
}

func DeCompress(data string) string {
	in := bytes.NewReader([]byte(data))
	gz, _ := gzip.NewReader(in)
	res, _ := ioutil.ReadAll(gz)
	return string(res)
}
