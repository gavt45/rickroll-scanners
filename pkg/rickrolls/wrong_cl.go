package rickrolls

import "net/http"

func WrongContentLength(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Length", "1000")
	res.Write([]byte(""))
}
