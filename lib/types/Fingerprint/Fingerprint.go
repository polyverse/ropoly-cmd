package Fingerprint

import (
    "crypto/md5"
    "encoding/binary"
	"encoding/json"
	"github.com/polyverse/ropoly-cmd/lib/architectures"
	"github.com/polyverse/ropoly-cmd/lib/types/BinDump"
	"github.com/polyverse/ropoly-cmd/lib/types/Gadget"
	"io/ioutil"
)

type Contents map[uint64][]uint64

type Fingerprint struct {
	Contents Contents
	GadgetSpec Gadget.UserSpec
}

func GenerateFingerprintFromGadgets(gadgets []*Gadget.GadgetInstance, gadgetSpec Gadget.UserSpec) Fingerprint {
    setContents := make(map[uint64]map[uint64]bool)
    for _, gadget := range gadgets {
        hash := md5.Sum(gadget.Gadget.Bytes())
        keyBytes := hash[8:16]
        key := binary.LittleEndian.Uint64(keyBytes)
        if setContents[key] == nil {
            setContents[key] = make(map[uint64]bool)
        }
        setContents[key][gadget.Address] = true
    }
    contents := make(Contents)
    for gadgetKey, addresses := range setContents {
        for address := range addresses {
            contents[gadgetKey] = append(contents[gadgetKey], address)
        }
    }
	return Fingerprint{
		contents,
		gadgetSpec,
	}
}

func GenerateFingerprintFromBinDump(binary BinDump.BinDump, gadgetSpec Gadget.UserSpec) (Fingerprint, error) {
	gadgets, err := architectures.GetGadgets(binary, gadgetSpec)
	if err != nil {
		return Fingerprint{}, err
	}

	return GenerateFingerprintFromGadgets(gadgets, gadgetSpec), nil
}

func OpenFingerprint(path string) (Fingerprint, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return Fingerprint{}, err
	}

	var fingerprint Fingerprint
	err = json.Unmarshal(b, &fingerprint)
	if err != nil {
		return Fingerprint{}, err
	}

	return fingerprint, nil
}

func EncodeFingerprint(fingerprint Fingerprint) (string, error) {
	b, err := json.MarshalIndent(fingerprint, "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}