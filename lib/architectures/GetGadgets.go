package architectures

import (
	"debug/elf"
	"github.com/polyverse/ropoly-cmd/lib/types/BinDump"
	"github.com/polyverse/ropoly-cmd/lib/types/Gadget"
	"errors"
	"strconv"
)

func getGadgets(segment BinDump.Segment, gadgetSpecs []*Gadget.EndSpec, decodeGadget Gadget.DecoderFunc, minLength uint64, maxLength uint64) ([]*Gadget.GadgetInstance, error) {
	vaddrsToGadgetIDs := make(map[uint64]map[int]string)
	gadgetInstances := []*Gadget.GadgetInstance{}
	for _, gadgetSpec := range gadgetSpecs {
		for match, err := gadgetSpec.Opcode.FindBytesMatchStartingAt(segment.Contents, 0); match != nil; match, err = gadgetSpec.Opcode.FindNextOverlappingMatch(match) {
			if err != nil {
				return nil, err
			}

			for i := uint64(0); i < maxLength; i++ {
				if (segment.Addr+uint64(match.Index)-(i*gadgetSpec.Align))%gadgetSpec.Align == 0 {

					// Get the probable gadget at alignment
					start := uint64(match.Index) - (i * gadgetSpec.Align)
					end := match.Index + gadgetSpec.Size
					if start >= uint64(end) || end >= len(segment.Contents) {
						continue
					}
					opcode := segment.Contents[start:end]

					// Disassemble it
					gadget, err := decodeGadget(opcode)

					if err != nil {
						//return nil, err // TODO: log soft error
						continue
					}

					if uint64(len(gadget)) < minLength {
						continue
					}

					// Ensure the byte sequence matches the regex we're parsing
					if match, err := gadgetSpec.Opcode.FindBytesMatchStartingAt(gadget.Bytes(), 0); err != nil || match == nil {
						continue
					}

					vaddr := segment.Addr + uint64(match.Index) - (i * gadgetSpec.Align)
					if vaddrsToGadgetIDs[vaddr] == nil {
						vaddrsToGadgetIDs[vaddr] = make(map[int]string)
					}
					if vaddrsToGadgetIDs[vaddr][len(gadget.Bytes())] == "" {
						vaddrsToGadgetIDs[vaddr][len(gadget.Bytes())] = gadget.Bytes().String()
						gadgetInstance := &Gadget.GadgetInstance{
							Address: vaddr,
							Gadget: gadget,
						}
						gadgetInstances = append(gadgetInstances, gadgetInstance)
					} else {
						if vaddrsToGadgetIDs[vaddr][len(gadget.Bytes())] != gadget.Bytes().String() {
							panic("Different gadgets found at same address:\n" + vaddrsToGadgetIDs[vaddr][len(gadget.Bytes())] + "\n" + gadget.Bytes().String())
						}
					}
				}
			}
		}
	}
	return gadgetInstances, nil
}

func GetGadgets(binary BinDump.BinDump, spec Gadget.UserSpec) ([]*Gadget.GadgetInstance, error) {
	var arch Architecture
	switch binary.Machine {
	case elf.EM_386:
		arch = X86
	case elf.EM_960:
		arch = X86
	case elf.EM_ARM:
		arch = ARM
	case elf.EM_X86_64:
		arch = X86
	case elf.EM_IA_64:
		arch = X86
	case elf.EM_8051:
		arch = X86
	case elf.EM_L10M:
		arch = X86
	case elf.EM_K10M:
		arch = X86
	case elf.EM_AARCH64:
		arch = ARM
	default:
		return nil, errors.New("unsupported architecture " + strconv.FormatUint(uint64(binary.Machine), 10))
	}

	allInstances := []*Gadget.GadgetInstance{}
	for _, segment := range binary.Segments {
		instances, err := getGadgets(segment, GadgetSpecLists[arch], GadgetDecoderFuncs[arch], uint64(spec.MinSize + 1), uint64(spec.MaxSize + 1))
		if err != nil {
			return nil, err
		}
		allInstances = append(allInstances, instances...)
	}
	return allInstances, nil
}