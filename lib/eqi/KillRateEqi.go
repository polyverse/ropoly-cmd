package eqi

// import (
// 	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
// )
//
// func KillRateEqi(f1, f2 Fingerprint.Fingerprint) float64 {
// 	gadgetInstanceCount := 0
// 	survivedCount := 0
//
// 	for key, f1Addresses := range f1.Contents {
// 		gadgetInstanceCount += len(f1Addresses)
// 		f2Addresses := f2.Contents[key]
// 		index1 := 0
// 		index2 := 0
// 		for index1 < len(f1Addresses) && index2 < len(f2Addresses) {
// 			f1Address := f1Addresses[index1]
// 			f2Address := f2Addresses[index2]
// 			if f1Address < f2Address {
// 				index1++
// 			} else if f1Address == f2Address {
// 				survivedCount++
// 				index1++
// 				index2++
// 			} else {
// 				index2++
// 			}
// 		}
// 	}
//
// 	killCount := gadgetInstanceCount - survivedCount
// 	killRate := float64(killCount) / float64(gadgetInstanceCount)
// 	return killRate * float64(100.0)
// }