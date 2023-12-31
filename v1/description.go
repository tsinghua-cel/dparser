package v1

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	ecdsaprysm "github.com/prysmaticlabs/prysm/v4/crypto/ecdsa"
	"github.com/tsinghua-cel/dparser/types"
	"log"
)

type P2PInfo struct {
	PrivateKey *ecdsa.PrivateKey
	IP         int
	P2PId      peer.ID
}

func ValidateDescription(d types.Description) error {
	cached := make(map[string]interface{})
	for _, e := range d.Topology.Executor {
		if e.Name == "" {
			return types.ErrInvalidName
		}

		if e.Version == "" {
			return types.ErrInvalidVersion
		}

		if _, ok := cached[e.Name]; ok {
			return types.ErrDuplicatedModule
		}
		cached[e.Name] = e
	}

	for _, beacon := range d.Topology.Beacons {
		if beacon.Name == "" {
			return types.ErrInvalidName
		}

		if beacon.Version == "" {
			return types.ErrInvalidVersion
		}

		if _, ok := cached[beacon.Name]; ok {
			return types.ErrDuplicatedModule
		}
		if beacon.Executor == "" || cached[beacon.Executor] == nil {
			return types.ErrNotFindDependency
		}
		cached[beacon.Name] = beacon
	}
	for _, beacon := range d.Topology.Beacons {
		for _, peer := range beacon.Peers {
			if cached[peer] == nil {
				log.Fatalf("not found %s.peer(%s)\n", beacon.Name, peer)
				return types.ErrNotFindDependency
			}
		}
	}

	for _, validator := range d.Topology.Validators {
		if validator.Name == "" {
			return types.ErrInvalidName
		}

		if validator.Version == "" {
			return types.ErrInvalidVersion
		}

		if _, ok := cached[validator.Name]; ok {
			return types.ErrDuplicatedModule
		}

		if validator.Beacon == "" || cached[validator.Beacon] == nil {
			log.Fatalf("not found %s.beacon(%s)\n", validator.Name, validator.Beacon)
			return types.ErrNotFindDependency
		}

		cached[validator.Name] = validator
	}

	return nil
}

func GetBeaconP2PInfo(d types.Description) map[string]P2PInfo {
	p2pInfo := make(map[string]P2PInfo)
	ipbase := 2
	for _, beacon := range d.Topology.Beacons {
		var privkey = new(ecdsa.PrivateKey)
		if beacon.P2PKey == "" {
			priv, err := GenerateP2PKey()
			if err != nil {
				log.Fatal("Failed to generate private key")
			}
			privkey = priv
		} else {
			priv, err := GetP2PKeyFromHex(beacon.P2PKey)
			if err != nil {
				log.Fatal("Failed to retrieve private key")
			}
			privkey = priv
		}
		id := GetPid(privkey)
		p2pInfo[beacon.Name] = P2PInfo{
			PrivateKey: privkey,
			P2PId:      id,
			IP:         ipbase,
		}
		ipbase++
	}
	return p2pInfo
}

func GetP2PKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {
	dst := make([]byte, hex.DecodedLen(len(hexKey)))
	_, err := hex.Decode(dst, []byte(hexKey))
	if err != nil {
		return nil, errors.New("failed to decode hex string")
	}
	unmarshalledKey, err := crypto.UnmarshalSecp256k1PrivateKey(dst)
	if err != nil {
		return nil, err
	}
	return ecdsaprysm.ConvertFromInterfacePrivKey(unmarshalledKey)
}

func GenerateP2PKey() (*ecdsa.PrivateKey, error) {
	priv, _, err := crypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	return ecdsaprysm.ConvertFromInterfacePrivKey(priv)
}

func GetPid(priKey *ecdsa.PrivateKey) peer.ID {
	ifaceKey, err := ecdsaprysm.ConvertToInterfacePrivkey(priKey)
	if err != nil {
		log.Fatal("Failed to retrieve private key")
	}
	id, err := peer.IDFromPublicKey(ifaceKey.GetPublic())
	if err != nil {
		log.Fatal("Failed to retrieve peer id")
	}
	return id
}
