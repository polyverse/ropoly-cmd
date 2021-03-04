package eqi

import (
	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
	"math/rand"
)

func SharedOffsetExistsMonteCarloEqi(numGadgets uint, trials uint) func(Fingerprint.Fingerprint, Fingerprint.Fingerprint) float64 {
	return func(f1, f2 Fingerprint.Fingerprint) float64 {

		gadgetsToIndices := make(map[string]int)
		workingIndex := 0
		for key, addresses := range f1.Contents {
			gadgetsToIndices[key] = workingIndex
			workingIndex += len(addresses)
		}
		gadgetInstanceCount := workingIndex

		displacementFoundCount := uint(0)
		for i := uint(0); i < trials; i++ {

			gadgetInstanceIndices := make([]int, numGadgets, numGadgets)
			for j := uint(0); j < numGadgets; j++ {
				gadgetInstanceIndices[j] = rand.Intn(gadgetInstanceCount)
			}

			var sharedDisplacements map[int64]bool
			for k, gadgetInstanceIndex := range gadgetInstanceIndices {
				firstIteration := k == 0

				var gadgetKey string
				var f1Address uint64
				for key, keyIndex := range gadgetsToIndices {
					if keyIndex <= gadgetInstanceIndex && gadgetInstanceIndex < keyIndex + len(f1.Contents[key]) {
						gadgetKey = key
						f1Address = f1.Contents[key][gadgetInstanceIndex -keyIndex]
						break
					}
				}

				f2Addresses := f2.Contents[gadgetKey]
				displacements := make(map[int64]bool)
				for _, f2Address := range f2Addresses {
					displacement := int64(f2Address) - int64(f1Address)
					displacements[displacement] = true
				}

				if firstIteration {
					sharedDisplacements = displacements
				} else {
					for displacement, sharedSoFar := range sharedDisplacements {
						sharedDisplacements[displacement] = sharedSoFar && displacements[displacement]
					}
				}
			}

			sharedDisplacementExists := false
			for _, actuallyShared := range sharedDisplacements {
				if actuallyShared {
					sharedDisplacementExists = true
					break
				}
			}
			if sharedDisplacementExists {
				displacementFoundCount++
			}
		}

		noCommonDisplacementCount := trials - displacementFoundCount
		return float64(100.0) * (float64(noCommonDisplacementCount) / float64(trials))
	}
}
