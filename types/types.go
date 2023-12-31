package types

type Validator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Beacon  string `json:"beacon"`
}

type Beacon struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Executor string   `json:"executor"`
	MaxPeers int      `json:"max-peers"`
	Peers    []string `json:"peers"`
	P2PKey   string   `json:"p2p-key"`
}

type Execute struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Config struct {
	TotalTime int `json:"total-time"` // seconds time to run the test.
}

type Topology struct {
	Executor   []Execute   `json:"executors"`
	Beacons    []Beacon    `json:"beacons"`
	Validators []Validator `json:"validators"`
}

type Description struct {
	Version  string   `json:"version"`
	Config   Config   `json:"config"`
	Topology Topology `json:"topology"`
}
