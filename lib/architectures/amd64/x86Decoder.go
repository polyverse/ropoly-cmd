package amd64

import (
	"github.com/polyverse/ropoly-cmd/lib/types/Gadget"
	"github.com/polyverse/ropoly-cmd/lib/types/Instruction"
	"fmt"
	"golang.org/x/arch/x86/x86asm"
)

func InstructionDecoder(opcodes []byte) (instruction *Instruction.Instruction, err error) {
	var inst x86asm.Inst

	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Unable to decode instruction due to disassembler panic: %v", x)
		}
	}()

	inst, err = x86asm.Decode(opcodes, 64)
	if err != nil {
		return
	}

	instruction = &Instruction.Instruction{
		Octets: opcodes[0:inst.Len],
		DisAsm: inst.String(),
	}
	return
}

func GadgetDecoder(opcodes []byte) (Gadget.Gadget, error) {
	gadget := Gadget.Gadget{}

	for len(opcodes) > 0 {
		instr, err := InstructionDecoder(opcodes)
		if err != nil {
			return nil, err
		}
		gadget = append(gadget, instr)
		gadlen := len(instr.Octets)
		if len(opcodes) <= gadlen {
			break
		}

		opcodes = opcodes[gadlen:]
	}
	return gadget, nil
}
