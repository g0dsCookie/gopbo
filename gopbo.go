package main

import (
	"fmt"

	"github.com/g0dsCookie/gopbo/pbo"
)

const file = `/home/g0dscookie/Downloads/@ExileServer-1.0.3f/Arma 3 Server/@ExileServer/addons/exile_server_config.pbo`

func main() {
	f, err := pbo.Load(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(f)
}
