package Instruction

type DecoderFunc func([]byte) (*Instruction, error)
