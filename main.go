package main

import (
	v1 "dparser/v1"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
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
	if err = description.Validate(); err != nil {
		log.Fatalf("Failed to validate description data: %v", err)
	}
	// generate Dockerfile
	buildScript, err := os.OpenFile("generated/build.sh", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open build.sh: %v", err)
	}

	// 1. generate Dockerfile for each execute
	{
		gethDockTml, err := template.ParseFiles("templates/geth.Dockerfile.tmpl")
		if err != nil {
			log.Fatalf("Failed to parse geth.Dockerfile.tmpl: %v", err)
		}
		braches := make(map[string]bool)
		for _, execute := range description.Topology.Executor {
			if _, ok := braches[execute.Version]; ok {
				continue
			}
			fs, err := os.OpenFile(fmt.Sprintf("generated/geth.Dockerfile.%s", execute.Version), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Failed to open geth.Dockerfile.%s: %v", execute.Version, err)
			}

			if err = gethDockTml.Execute(fs, execute); err != nil {
				log.Fatalf("Failed to execute geth.Dockerfile.tmpl: %v", err)
			}
			fs.Close()
			braches[execute.Version] = true
			buildScript.WriteString(fmt.Sprintf("docker build -t geth:%s -f generated/geth.Dockerfile.%s .\n", execute.Version, execute.Version))
		}
	}

	// 2. generate Dockerfile for each beacon
	{
		beaconDockTml, err := template.ParseFiles("templates/beacon.Dockerfile.tmpl")
		if err != nil {
			log.Fatalf("Failed to parse beacon.Dockerfile.tmpl: %v", err)
		}
		braches := make(map[string]bool)
		for _, beacon := range description.Topology.Beacons {
			if _, ok := braches[beacon.Version]; ok {
				continue
			}
			fs, err := os.OpenFile(fmt.Sprintf("generated/beacon.Dockerfile.%s", beacon.Version), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Failed to open beacon.Dockerfile.%s: %v", beacon.Version, err)
			}

			if err = beaconDockTml.Execute(fs, beacon); err != nil {
				log.Fatalf("Failed to beacon geth.Dockerfile.tmpl: %v", err)
			}
			fs.Close()
			braches[beacon.Version] = true
			buildScript.WriteString(fmt.Sprintf("docker build -t beacon:%s -f generated/beacon.Dockerfile.%s .\n", beacon.Version, beacon.Version))
		}

	}
	// 3. generate Dockerfile for each validator
	{
		validatorDockTml, err := template.ParseFiles("templates/validator.Dockerfile.tmpl")
		if err != nil {
			log.Fatalf("Failed to parse validator.Dockerfile.tmpl: %v", err)
		}
		braches := make(map[string]bool)
		for _, validator := range description.Topology.Validators {
			if _, ok := braches[validator.Version]; ok {
				continue
			}
			fs, err := os.OpenFile(fmt.Sprintf("generated/validator.Dockerfile.%s", validator.Version), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Failed to open validator.Dockerfile.%s: %v", validator.Version, err)
			}

			if err = validatorDockTml.Execute(fs, validator); err != nil {
				log.Fatalf("Failed to validator geth.Dockerfile.tmpl: %v", err)
			}
			fs.Close()
			braches[validator.Version] = true
			buildScript.WriteString(fmt.Sprintf("docker build -t validator:%s -f generated/validator.Dockerfile.%s .\n", validator.Version, validator.Version))
		}
	}
	buildScript.Close()

	// generate docker-compose.yml

}
