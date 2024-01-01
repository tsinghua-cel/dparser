package compose

import (
	"bytes"
	"fmt"
	"github.com/tsinghua-cel/dparser/types"
	v1 "github.com/tsinghua-cel/dparser/v1"
	"log"
	"os"
	"text/template"
)

func BuildCompose(d types.Description, output string) error {
	beaconP2pinfo := v1.GetBeaconP2PInfo(d)
	buffer := bytes.NewBufferString("")
	buffer.WriteString(composeHeader)
	// build all execute
	baseExecuteAuthPort := 10000
	baseExecuteRPCPort := 20000
	for idx, execute := range d.Topology.Executor {
		var config ExecuteConfig
		config.ExecuteName = execute.Name
		config.ExecuteImage = fmt.Sprintf("geth:%s", execute.Version)
		config.ExecuteAuthPort = fmt.Sprintf("%d", baseExecuteAuthPort+idx)
		config.ExecuteRPCPort = fmt.Sprintf("%d", baseExecuteRPCPort+idx)
		config.ExecuteDataPath = fmt.Sprintf("%s", execute.Name)

		tmpl, err := template.New("test").Parse(executeTmpl)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(buffer, config)
		if err != nil {
			log.Fatalf("Failed to execute executeTmpl: %v", err)
		}
	}

	// build all beacon
	// prepare some default config.
	var allPeers = make([]string, 0, len(d.Topology.Beacons))
	for _, beacon := range d.Topology.Beacons {
		allPeers = append(allPeers, beacon.Name)
	}
	var defaultMaxPeers = 70
	for _, beacon := range d.Topology.Beacons {
		var config BeaconConfig
		config.BeaconName = beacon.Name
		config.BeaconImage = fmt.Sprintf("beacon:%s", beacon.Version)
		config.BeaconDataPath = fmt.Sprintf("%s", beacon.Name)
		config.BeaconIP = fmt.Sprintf("172.99.1.%d", beaconP2pinfo[beacon.Name].IP)
		config.ExecuteName = beacon.Executor
		config.BeaconMaxPeers = beacon.MaxPeers
		if config.BeaconMaxPeers == 0 {
			config.BeaconMaxPeers = defaultMaxPeers
		}
		config.BeaconP2PKey = leftPadding(beaconP2pinfo[beacon.Name].PrivateKey.D.Text(16), 64)
		peers := allPeers
		if len(beacon.Peers) > 0 {
			peers = beacon.Peers
		}

		for _, peer := range peers {
			if peer == beacon.Name {
				continue
			}
			// --peer /ip4/172.99.1.1/tcp/13000/p2p/16Uiu2HAmHwS8xvw3T5nMKW6Cq9drWKov2P7fcFECq59d6U86dM59
			config.BeaconPeers += fmt.Sprintf(" --peer /ip4/172.99.1.%d/tcp/13000/p2p/%s ", beaconP2pinfo[peer].IP, beaconP2pinfo[peer].P2PId)
		}

		tmpl, err := template.New("test").Parse(beaconTmpl)
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(buffer, config)
		if err != nil {
			log.Fatalf("Failed to execute beaconTmpl: %v", err)
		}
	}

	// build all validator
	nodeCount := len(d.Topology.Validators)
	for idx, validator := range d.Topology.Validators {
		index := idx
		step := totalValidators / nodeCount
		startIndex := index * step
		endIndex := startIndex + step - 1
		if index == nodeCount-1 {
			endIndex = totalValidators - 1
		}
		validatorsNum := endIndex - startIndex + 1

		var config ValidatorConfig
		config.ValidatorName = validator.Name
		config.ValidatorImage = fmt.Sprintf("validator:%s", validator.Version)
		config.BeaconName = validator.Beacon
		config.ValidatorDataPath = fmt.Sprintf("%s", validator.Name)
		config.ValidatorNum = validatorsNum
		config.ValidatorStartIndex = startIndex

		tmpl, err := template.New("test").Parse(validatorTmpl)
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(buffer, config)
		if err != nil {
			log.Fatalf("Failed to execute validatorTmpl: %v", err)
		}
	}

	buffer.WriteString(composeNetwork)

	fs, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open docker-compose.yaml: %v", err)
	}
	fs.Write(buffer.Bytes())
	fs.Close()

	return nil
}

func leftPadding(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return leftPadding("0"+str, length)
}
