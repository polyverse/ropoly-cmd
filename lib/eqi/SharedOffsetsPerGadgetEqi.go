package eqi

import (
	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
)

func SharedOffsetsPerGadgetEqi(f1, f2 Fingerprint.Fingerprint) float64 {
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

    sharedCountSum := uint(0)
    for key, f1Addresses := range f1.Contents {
        f2Addresses := f2.Contents[key]
        for _, f1Address := range f1Addresses {
            sharedCount := uint(0)
            for _, f2Address := range f2Addresses {
                displacement := int64(f2Address) - int64(f1Address)
                sharedCountCandidate := gadgetInstanceCountsByDisplacement[displacement]
                if sharedCountCandidate > sharedCount {
                    sharedCount = sharedCountCandidate
                }
            }
            sharedCountSum += sharedCount
        }
    }
	sharedProportionSum := float64(sharedCountSum) / float64(gadgetInstanceCount)
	sharedProportionMean := sharedProportionSum / float64(gadgetInstanceCount)
	return (float64(1.0) - sharedProportionMean) * float64(100.0)
}