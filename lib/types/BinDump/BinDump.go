package BinDump

import (
	"debug/elf"
	"encoding/json"
	"io/ioutil"
)

type Segment struct {
	Addr uint64     `json:"address"`
	Contents []byte `json:"contents"`
}

type BinDump struct {
	Machine elf.Machine `json:"machine"`
	Segments []Segment  `json:"segments"`
}

func GenerateElfBinDump(path string) (BinDump, error) {
	elfFile, err := elf.Open(path)
	if err != nil {
		return BinDump{}, err
	}

	segments := make([]Segment, 0, len(elfFile.Progs))
	for _, programHeader := range elfFile.Progs {
		if programHeader.Flags & elf.PF_X != 0 {

			size := programHeader.Filesz
			segmentContents := make([]byte, size, size)
			_, err := programHeader.Open().Read(segmentContents)
			if err != nil {
				return BinDump{}, err
			}

			segment := Segment{
				Addr: programHeader.Vaddr,
				Contents: segmentContents,
			}
			segments = append(segments, segment)
		}
	}

	return BinDump{
		Machine: elfFile.Machine,
		Segments: segments,
	}, nil
}

func OpenBinDump(path string) (BinDump, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return BinDump{}, err
	}

	var binDump BinDump
	err = json.Unmarshal(b, &binDump)
	if err != nil {
		return BinDump{}, err
	}

	return binDump, nil
}

func EncodeBinDump(binDump BinDump) (string, error) {
	b, err := json.MarshalIndent(binDump, "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}