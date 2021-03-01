package thumb

import (
	"github.com/polyverse/ropoly-cmd/lib/types/Gadget"
	"github.com/polyverse/ropoly-cmd/lib/types/Instruction"
	"strconv"
)

func GadgetDecoder(opcodes []byte) (Gadget.Gadget, error) {
	gadget := Gadget.Gadget{}

	for len(opcodes) > 0 {
		instr := Instruction.Instruction {
			Octets: opcodes[0:2],
			DisAsm: formatByte(opcodes[0]) + formatByte(opcodes[1]),
		}
		gadget = append(gadget, &instr)
		gadlen := len(instr.Octets)
		if len(opcodes) <= gadlen {
			break
		}

		opcodes = opcodes[gadlen:]
	}
	return gadget, nil
}

func formatByte(b byte) string {
	ret := strconv.FormatUint(uint64(b), 16)
	if len(ret) == 1 {
		ret = "0" + ret
	}
	return ret
}