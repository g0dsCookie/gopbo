package main

import (
	"github.com/g0dsCookie/gopbo/pbo"
)

const file = `/home/g0dscookie/Downloads/@ExileServer-1.0.3f/Arma 3 Server/@ExileServer/addons/exile_server.pbo`
const dir = `/home/g0dscookie/Downloads/@ExileServer-1.0.3f/Arma 3 Server/@ExileServer/addons/exile_server`

func main() {
	if err := pbo.UnpackVerbose(file, dir); err != nil {
		panic(err)
	}
}
