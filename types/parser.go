package types

import "crypto/ecdsa"

type BaseModule interface {
	Name() string
	Version() string
}

type ExecutorDes interface {
	BaseModule
}

type BeaconDes interface {
	BaseModule
	Execute() string
	Peers() []string
	MaxPeers() int
	P2PKey() string
	SetP2PKey(*ecdsa.PrivateKey)
	P2PId() string
}

type ValidatorDes interface {
	BaseModule
	Beacon() string
}

type Parser interface {
	Parse(string) Description
}

type Description interface {
	ExecuteNodes() []ExecutorDes
	BeaconNodes() []BeaconDes
	Validators() []ValidatorDes
}

type dummyParser struct{}

func (d *dummyParser) Parse(string) Description {
	return nil
}

func NewParser() Parser {
	return &dummyParser{}
}
