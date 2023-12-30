package v1

import "dparser/types"

type Validator struct {
	name    string `json:"name"`
	version string `json:"version"`
	beacon  string `json:"beacon"`
}

// validator implement ValidatorDes interface
func (v Validator) Name() string {
	return v.name
}
func (v Validator) Version() string {
	return v.version
}
func (v Validator) Beacon() string {
	return v.beacon
}

type Beacon struct {
	name     string   `json:"name"`
	version  string   `json:"version"`
	executor string   `json:"executor"`
	maxPeers int      `json:"max-peers"`
	peers    []string `json:"peers"`
	p2pKey   string   `json:"p2p-key"`
}

// Beacon implement BeaconDes interfaces
func (b Beacon) Name() string {
	return b.name
}
func (b Beacon) Version() string {
	return b.version
}
func (b Beacon) Execute() string {
	return b.executor
}
func (b Beacon) Peers() []string {
	return b.peers
}
func (b Beacon) MaxPeers() int {
	return b.maxPeers
}
func (b Beacon) P2PKey() string {
	return b.p2pKey
}
func (b Beacon) P2PId() string {
	return b.p2pKey
}

type Execute struct {
	name    string `json:"name"`
	version string `json:"version"`
}

// Execute implement ExecuteDes interface
func (e Execute) Name() string {
	return e.name
}
func (e Execute) Version() string {
	return e.version
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

// Description implement Description interface
func (d Description) ExecuteNodes() []types.ExecutorDes {
	var es []types.ExecutorDes
	for _, e := range d.Topology.Executor {
		es = append(es, e)
	}
	return es
}
func (d Description) BeaconNodes() []types.BeaconDes {
	var bs []types.BeaconDes
	for _, b := range d.Topology.Beacons {
		bs = append(bs, b)
	}
	return bs
}
func (d Description) Validators() []types.ValidatorDes {
	var vs []types.ValidatorDes
	for _, v := range d.Topology.Validators {
		vs = append(vs, v)
	}
	return vs
}
