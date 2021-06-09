package cmd

import (
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
    "github.com/spf13/cobra"
    "strconv"
)

func init() {
    rootCmd.AddCommand(countCmd)
}

var countCmd = &cobra.Command{
    Use:        "count <path>",
    Short:      "Counts the gadgets in a fingerprint",
    Args:       cobra.ExactArgs(1),
    RunE:       func(cmd *cobra.Command, args []string) error {

        fingerprint, err := Fingerprint.OpenFingerprint(args[0])
            if err != nil {
                return err
            }

        count := 0
        for _, addresses := range fingerprint.Contents {
            count += len(addresses)
        }

        countString := strconv.FormatInt(int64(count), 10)
        fmt.Print(countString)
        return nil
    },
}