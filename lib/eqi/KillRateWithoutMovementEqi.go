package eqi

// import (
// 	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
// )
//
// func KillRateWithoutMovementEqi(f1, f2 Fingerprint.Fingerprint) float64 {
// 	gadgetInstanceCount := 0
// 	survivedCount := 0
//
// 	for key, f1Addresses := range f1.Contents {
// 		gadgetInstanceCount += len(f1Addresses)
// 		if len(f2.Contents[key]) > 0 {
// 			survivedCount += len(f1Addresses)
// 		}
// 	}
//
// 	killCount := gadgetInstanceCount - survivedCount
// 	killRate := float64(killCount) / float64(gadgetInstanceCount)
// 	return killRate * float64(100.0)
// }