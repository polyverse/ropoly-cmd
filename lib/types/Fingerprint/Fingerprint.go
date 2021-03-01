package Fingerprint

import (
	"encoding/json"
	"github.com/polyverse/ropoly-cmd/lib/architectures"
	"github.com/polyverse/ropoly-cmd/lib/types/BinDump"
	"github.com/polyverse/ropoly-cmd/lib/types/Gadget"
	"io/ioutil"
)

type Contents map[string][]uint64

type Fingerprint struct {
	Contents Contents
	GadgetSpec Gadget.UserSpec
}

func GenerateFingerprintFromGadgets(gadgets []*Gadget.GadgetInstance, gadgetSpec Gadget.UserSpec) Fingerprint {
	contents := make(Contents)
	for _, gadget := range gadgets {
		key := gadget.Gadget.Bytes().String()
		contents[key] = append(contents[key], gadget.Address)
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