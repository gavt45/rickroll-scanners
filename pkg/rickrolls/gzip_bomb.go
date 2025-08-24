package rickrolls

import (
	"bytes"
	_ "embed"
	"io"
	"net/http"
)

//go:embed bomb.gzip
var b []byte

var bombReader io.Reader

func init() {
	bombReader = bytes.NewReader(b)
}

func GzipBombSimple(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Encoding", "gzip")
	res.Header().Set("Content-Type", "text/html")
	io.Copy(res, bombReader)
}
