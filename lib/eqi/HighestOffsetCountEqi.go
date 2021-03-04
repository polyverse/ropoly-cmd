package eqi

import (
	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
)

func HighestOffsetCountEqi(f1, f2 Fingerprint.Fingerprint) float64 {
	gadgetInstanceCount := 0
	gadgetInstanceCountsByDisplacement := make(map[int64]uint)

	for key, f1Addresses := range f1.Contents {
		gadgetInstanceCount += len(f1Addresses)
		f2Addresses := f2.Contents[key]
		for _, f1Address := range f1Addresses {
			for _, f2Address := range f2Addresses {
				displacement := int64(f2Address) - int64(f1Address)
				gadgetInstanceCountsByDisplacement[displacement] += 1
			}
		}
	}

	highestOffsetCount := uint(0)
	for _, offsetCount := range gadgetInstanceCountsByDisplacement {
		if offsetCount > highestOffsetCount {
			highestOffsetCount = offsetCount
		}
	}
	
	return float64(100.0) * (float64(1.0) - (float64(highestOffsetCount) / float64(gadgetInstanceCount)))
}