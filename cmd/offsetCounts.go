package cmd

import (
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
    "github.com/spf13/cobra"
    "strconv"
)

var depth int

func init() {
    offsetCountsCmd.PersistentFlags().IntVarP(&depth, "depth", "d", -1, "maximum number of offset counts to report")

    rootCmd.AddCommand(offsetCountsCmd)
}

var offsetCountsCmd = &cobra.Command {
    Use:    "offset-counts <path> <path>",
    Short:  "Compares two fingerprints and reports the greatest numbers of shared displacements gadgets move by",
    Args:   cobra.ExactArgs(2),
    RunE:       func(cmd *cobra.Command, args []string) error {
        f1, err := Fingerprint.OpenFingerprint(args[0])
        if err != nil {
           	return err
        }

        f2, err := Fingerprint.OpenFingerprint(args[1])
        if err != nil {
            return err
       	}

        blacklistedDisplacements := make(map[int64]bool)
        keepIterating := true

        for keepIterating && (depth < 0 || (len(blacklistedDisplacements) < depth)) {
            gadgetInstanceCountsByDisplacement := make(map[int64]uint)

            for key, f1Addresses := range f1.Contents {
                f2Addresses := f2.Contents[key]
                for _, f1Address := range f1Addresses {
                    considerThisGadget := true
                    displacements := make(map[int64]bool)
                    for _, f2Address := range f2Addresses {
                        displacement := int64(f2Address) - int64(f1Address)
                        if blacklistedDisplacements[displacement] {
                            considerThisGadget = false
                            break
                        }
                        displacements[displacement] = true
                    }
                    if considerThisGadget {
                        for displacement := range displacements {
                            gadgetInstanceCountsByDisplacement[displacement] += 1
                        }
                    }
                }
            }

            highestOffsetCount := uint(0)
            var highestOffsetCountOffset int64
            for offset, offsetCount := range gadgetInstanceCountsByDisplacement {
                if offsetCount > highestOffsetCount {
                    highestOffsetCount = offsetCount
                    highestOffsetCountOffset = offset
                }
            }
            if highestOffsetCount == 0 {
                keepIterating = false
            }
            blacklistedDisplacements[highestOffsetCountOffset] = true
            countString := strconv.FormatUint(uint64(highestOffsetCount), 10)
            fmt.Println(countString)
        }

        return nil
    },
}