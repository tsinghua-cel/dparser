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
	mysqlTmpl = `
  mysql:
    container_name: "ethmysql"
    image: "mysql:latest"
    environment:
      - MYSQL_ROOT_PASSWORD=12345678
    ports:
      - "3306:3306"
    restart: always
    privileged: true
    volumes:
      - "/etc/localtime:/etc/localtime"
      - "./data/mysql/data:/var/lib/mysql"
      - "./config/mysql/conf/my.cnf:/etc/my.cnf"
      - "./config/mysql/init:/docker-entrypoint-initdb.d/"
    networks:
      - meta
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
      - P2PKEY={{ .BeaconP2PKey }}
{{ .BeaconEnv }}
    deploy:
      restart_policy:
        condition: on-failure
        delay: 1s
        max_attempts: 100
        window: 120s
    ports:
      - "{{ .BeaconGrpcPort }}:4000"
      - "{{ .BeaconGrpcGwPort }}:3500"
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
      - VALIDATORS_NUM={{ .ValidatorNum }}
      - VALIDATORS_INDEX={{ .ValidatorStartIndex }}
      - BEACONRPC={{ .BeaconName }}:4000
{{ .ValidatorEnv }}
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
	attackerTmpl = `
  {{ .AttackerName}}:
    image: {{ .AttackerImage }}
    container_name: {{ .AttackerName }}
    entrypoint: /usr/local/bin/attacker.sh
    environment:
      - NAME={{ .AttackerName }}
{{ .AttackerEnv }}
    deploy:
      restart_policy:
        condition: on-failure
        delay: 1s
        max_attempts: 100
        window: 120s
    ports:
      - "{{ .AttackerPort }}:10000"
      - "{{ .SwagPort }}:10001"
    volumes:
      - {{ .AttackerConfig }}:/root/config.toml
      - {{ .AttackerStrategy }}:/root/strategy.json
      - ./data/{{ .AttackerDataPath }}:/root/attackerdata
    depends_on:
      - mysql
    networks:
      meta:
        ipv4_address: {{ .AttackerIP }}
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
	BeaconName       string
	BeaconImage      string
	BeaconIP         string
	BeaconPeers      string
	BeaconDataPath   string
	ExecuteName      string
	BeaconMaxPeers   int
	BeaconP2PKey     string
	BeaconEnv        string
	BeaconGrpcPort   int
	BeaconGrpcGwPort int
}

type ValidatorConfig struct {
	ValidatorName       string
	ValidatorImage      string
	ValidatorNum        int
	ValidatorStartIndex int
	ValidatorDataPath   string
	BeaconName          string
	ValidatorEnv        string
}

type AttackerConfig struct {
	AttackerName     string
	AttackerImage    string
	AttackerDataPath string
	AttackerEnv      string
	AttackerPort     int
	SwagPort         int
	AttackerConfig   string
	AttackerIP       string
	AttackerStrategy string
}
