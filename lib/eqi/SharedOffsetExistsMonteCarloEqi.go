package eqi

import (
	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
	"math/rand"
)

type InstanceId struct {
	Key uint64
	Address uint64
}

func SharedOffsetExistsMonteCarloEqi(numGadgets uint, trials uint) func(Fingerprint.Fingerprint, Fingerprint.Fingerprint) float64 {
	return func(f1, f2 Fingerprint.Fingerprint) float64 {

		indexedInstances := make([]InstanceId, 0, 0)
		for key, f1Addresses := range f1.Contents {
			for _, f1Address := range f1Addresses {
				instanceId := InstanceId{
					Key: key,
					Address: f1Address,
				}
				indexedInstances = append(indexedInstances, instanceId)
			}
		}

		noCommonDisplacementCount := 0
		for trial := uint(0); trial < trials; trial++ {

			selectedInstances := make([]InstanceId, numGadgets, numGadgets)
			for i := uint(0); i < numGadgets; i++ {
				instanceId := indexedInstances[rand.Intn(len(indexedInstances))]
				selectedInstances[i] = instanceId
			}

            commonDisplacements := displacementSetForF1Instance(f1, f2, selectedInstances[0])

			for i := uint(1); (len(commonDisplacements) != 0) && (i < numGadgets); i++ {
                nextInstanceDisplacements := displacementSetForF1Instance(f1, f2, selectedInstances[i])
				for displacement := range commonDisplacements {
					if !nextInstanceDisplacements[displacement] {
						delete(commonDisplacements, displacement)
					}
				}
			}

			if len(commonDisplacements) == 0 {
				noCommonDisplacementCount++
			}
		}

		return float64(100.0) * (float64(noCommonDisplacementCount) / float64(trials))
	}
}

func displacementSetForF1Instance(f1, f2 Fingerprint.Fingerprint, instanceId InstanceId) map[int64]bool {
    displacements := make(map[int64]bool)
    f2Addresses := f2.Contents[instanceId.Key]
    for _, f2Address := range f2Addresses {
        displacements[int64(f2Address) - int64(instanceId.Address)] = true
    }
    return displacements
}