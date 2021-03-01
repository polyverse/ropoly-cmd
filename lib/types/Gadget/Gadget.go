package Gadget

import (
	"github.com/polyverse/ropoly-cmd/lib/types/Instruction"
	"bytes"
	"strconv"
)

type Gadget []*Instruction.Instruction

type UserSpec struct {
	MinSize uint
	MaxSize uint
}

func (s UserSpec) String() string {
	return "MinSize=" + strconv.FormatUint(uint64(s.MinSize), 10) + ", MaxSize=" + strconv.FormatUint(uint64(s.MaxSize), 10)
}

type GadgetInstance struct {
	Address uint64
	Gadget Gadget
}

func (g Gadget) Bytes() Instruction.Octets {
	buffer := bytes.NewBuffer(nil)
	for _, instr := range g {
		buffer.Write(instr.Octets)
	}
	return buffer.Bytes()
}