package arm

import "github.com/polyverse/ropoly-cmd/lib/types/Gadget"

// Adapted from:
// 	https://github.com/JonathanSalwan/ROPgadget/blob/master/ropgadget/gadgets.py
// https://github.com/polyverse/EnVisen/blob/master/internaljs/instruction_gadget_worker.js
var GadgetSpecs = []*Gadget.EndSpec{
	// SYS Gadgets
	{Gadget.MustCompile("\x00-\xff]{3}\xef"), 4, 4},

	// JOP Gadgets
	{Gadget.MustCompile("[\x10-\x19\x1e]{1}\xff\x2f\xe1"), 4, 4},
	{Gadget.MustCompile("[\x30-\x39\x3e]{1}\xff\x2f\xe1"), 4, 4},
	{Gadget.MustCompile("[\x00-\xff][\x80-\xff][\x10-\x1e\x30-\x3e\x50-\x5e\x70-\x7e\x90-\x9e\xb0-\xbe\xd0-\xde\xf0-\xfe][\xe8\xe9]"), 4, 4},

	// No ROP Gadgets
}