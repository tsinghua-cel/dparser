package main

import (
	v1 "dparser/v1"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	fileName = flag.String("file", "description.json", "description file name")
)

func main() {
	flag.Parse()
	descriptionData, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Fatalf("Failed to read description file: %v", err)
	}
	description := new(v1.Description)
	if err = json.Unmarshal(descriptionData, description); err != nil {
		log.Fatalf("Failed to unmarshal description data: %v", err)
	}
	fmt.Printf("Description: %+v\n", description)
	fmt.Println("Hello, playground")
}
