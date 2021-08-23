package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "unsafe"
)

func init() {
    returnsCmd.PersistentFlags().BoolVarP(&inputIsPid, "pid", "p", false, "Input argument represents the PID of a running process, rather than a path to a binary file.")

    rootCmd.AddCommand(returnsCmd)
}

var returnsCmd = &cobra.Command{
    Use:        "returns <path/to/binary (or PID with --pid)>",
    Short:      "Finds the C3 return addresses in a binary file or running program",
    Args:       cobra.ExactArgs(1),
    RunE:       func(cmd *cobra.Command, args []string) error {
        f := FORM_FILEPATH
        if inputIsPid {
            f = FORM_PID
        }

        bindump, err := getInputAsBinDump(args[0], f)
        if err != nil {
       		return err
       	}

       	retsSet := make(map[uint32]bool)
       	for _, segment := range bindump.Segments {
       	    for address, b := range segment.Contents {
       	        if b == 0xc3 {
       	            fileAddress = address % 0x200000 // TODO: quick hack that could be wrong in some circumstances
       	            retsSet[uint32(address)] = true
       	        }
       	    }
       	}
       	retsBytes := make([]byte, 0, len(retsSet))
       	for ret := range retsSet {
       	    retsBytes = append(retsBytes, (*[4]byte)(unsafe.Pointer(&ret))[:]...)
       	}

        fmt.Print(string(retsBytes))
       	return nil
    },
}