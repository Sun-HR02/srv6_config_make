package main

import (
	"SRv6.Config.Builder/buildJson"
	"fmt"
	"log"
)

type test struct {
	Hello string `json:"hello, omitempty"`
}

func main() {
	bytes, err := buildJson.WriteJson()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v", string(bytes))

	return
}
