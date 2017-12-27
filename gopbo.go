package main

import (
	"github.com/g0dsCookie/gopbo/pbo"
)

const file = `/home/g0dscookie/Downloads/Arma 3 Server/@ExileServer/addons/exile_server.pbo`
const dir = `exile_server`

func main() {
	if err := pbo.UnpackVerbose(file, dir); err != nil {
		panic(err)
	}
}
