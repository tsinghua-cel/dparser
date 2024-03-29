package types

type Validator struct {
	Name           string                 `json:"name"`
	Version        string                 `json:"version"`
	Beacon         string                 `json:"beacon"`
	ValidatorCount int                    `json:"validator_count"`
	Env            map[string]interface{} `json:"env"`
}

type Beacon struct {
	Name     string                 `json:"name"`
	Version  string                 `json:"version"`
	Executor string                 `json:"executor"`
	MaxPeers int                    `json:"max-peers"`
	Peers    []string               `json:"peers"`
	P2PKey   string                 `json:"p2p-key"`
	Env      map[string]interface{} `json:"env"`
}

type Execute struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Attacker struct {
	Name     string                 `json:"name"`
	Version  string                 `json:"version"`
	Config   string                 `json:"config"`
	Strategy string                 `json:"strategy"`
	Env      map[string]interface{} `json:"env"`
}

type Config struct {
	TotalTime int `json:"total-time"` // seconds time to run the test.
}

type Topology struct {
	Executor   []Execute   `json:"executors"`
	Beacons    []Beacon    `json:"beacons"`
	Validators []Validator `json:"validators"`
	Attackers  []Attacker  `json:"attackers"`
}

type Description struct {
	Version  string   `json:"version"`
	Config   Config   `json:"config"`
	Topology Topology `json:"topology"`
}
