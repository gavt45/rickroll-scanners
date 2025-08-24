package rickrolls

import (
	"bytes"
	_ "embed"
	"io"
	"net/http"
)

//go:embed billion_lols.xml
var lols []byte

var billionsOfLol *bytes.Buffer

func init() {
	billionsOfLol = bytes.NewBuffer(lols)
}

func BillionLaughsSimple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "10000000000000")
	io.Copy(w, billionsOfLol)
}
