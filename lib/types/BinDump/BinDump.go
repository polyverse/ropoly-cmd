package BinDump

import (
	"debug/elf"
	"encoding/json"
	"github.com/polyverse/masche/memaccess"
	"github.com/polyverse/masche/process"
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

func GenerateElfBinDumpFromFile(path string) (BinDump, error) {
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

func GenerateBinDumpFromPid(pid uint) (BinDump, error) {
	softerrors := []error{} // TODO: do something with these
	proc := process.GetProcess(int(pid))

	segments := []Segment{}
	pc := uintptr(0)
	for {

		region, harderror, softerrors2 := memaccess.NextMemoryRegionAccess(proc, uintptr(pc), memaccess.Readable+memaccess.Executable)
		softerrors = append(softerrors, softerrors2...)
		if harderror != nil {
			return BinDump{}, harderror
		}
		if region == memaccess.NoRegionAvailable {
			break
		}

		//Make sure we move the Program Counter
		pc = region.Address + uintptr(region.Size)

		contents := make([]byte, region.Size, region.Size)
		harderror, softerrors2 = memaccess.CopyMemory(proc, region.Address, contents)
		softerrors = append(softerrors, softerrors2...)
		if harderror != nil {
			return BinDump{}, harderror
		}

		segment := Segment{
			uint64(region.Address),
			contents,
		}
		segments = append(segments, segment)
	}

	machine := elf.EM_X86_64 // TODO: assuming X86_64 for now

	return BinDump{
		machine,
		segments,
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