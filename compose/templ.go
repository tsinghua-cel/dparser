package compose

var (
	totalValidators = 64
	composeHeader   = `
version: "3.9"

services:`
	composeNetwork = `

networks:
  meta:
    driver: bridge
    ipam:
      config:
        - subnet: 172.99.0.0/16
`
	executeTmpl = `
  {{ .ExecuteName }}:
    image: {{ .ExecuteImage }}
    container_name: {{ .ExecuteName }}
    entrypoint: /usr/local/bin/execution.sh
    deploy:
      restart_policy:
        condition: on-failure
        delay: 1s
        max_attempts: 100
        window: 120s
    ports:
      - "{{ .ExecuteAuthPort }}:8551"
      - "{{ .ExecuteRPCPort }}:8545"
    volumes:
      - ./config:/root/config
      - ./data/{{ .ExecuteDataPath }}:/root/gethdata
    networks:
      - meta
`
	beaconTmpl = `
  {{ .BeaconName}}:
    image: {{ .BeaconImage }}
    container_name: {{ .BeaconName }}
    entrypoint: /usr/local/bin/beacon-node.sh
    environment:
      - ALLPEERS={{ .BeaconPeers }}
      - EXECUTE={{ .ExecuteName }}
      - MAXPEERS={{ .BeaconMaxPeers }}
    deploy:
      restart_policy:
        condition: on-failure
        delay: 1s
        max_attempts: 100
        window: 120s
    volumes:
      - ./config:/root/config
      - ./data/{{ .BeaconDataPath }}:/root/beacondata
    depends_on:
      - {{ .ExecuteName }}
    networks:
      meta:
        ipv4_address: {{ .BeaconIP }}
`
	validatorTmpl = `
  {{ .ValidatorName }}:
    image: {{ .ValidatorImage }}
    container_name: {{ .ValidatorName }}
    entrypoint: /usr/local/bin/validator.sh
    environment:
      - VALIDATORS_NUM= {{ .ValidatorNum }}
      - VALIDATORS_INDEX= {{ .ValidatorStartIndex }}
      - BEACONRPC={{ .BeaconName }}:4000
    deploy:
      restart_policy:
        condition: on-failure
        delay: 1s
        max_attempts: 100
        window: 120s
    volumes:
      - ./config:/root/config
      - ./data/{{ .ValidatorDataPath }}:/root/validatordata
    depends_on:
      - {{ .BeaconName }}
    networks:
      - meta
`
)

type ExecuteConfig struct {
	ExecuteName     string
	ExecuteImage    string
	ExecuteAuthPort string
	ExecuteRPCPort  string
	ExecuteDataPath string
}

type BeaconConfig struct {
	BeaconName     string
	BeaconImage    string
	BeaconIP       string
	BeaconPeers    string
	BeaconDataPath string
	ExecuteName    string
	BeaconMaxPeers int
}

type ValidatorConfig struct {
	ValidatorName       string
	ValidatorImage      string
	ValidatorNum        int
	ValidatorStartIndex int
	ValidatorDataPath   string
	BeaconName          string
}
