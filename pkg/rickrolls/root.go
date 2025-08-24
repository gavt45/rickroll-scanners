package rickrolls

import "net/http"

type RickRoll http.HandlerFunc

var RickRolls = []RickRoll{
	BillionLaughsSimple,
	GzipBombSimple,
	WrongContentLength,
}
