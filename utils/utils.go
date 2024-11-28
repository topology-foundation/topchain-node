package utils

import "math"

const EPOCH_SIZE = 10

// 7 days equivalent epochs
const DEAL_EXPIRY_CLAIM_WINDOW = 144

func ConvertBlockToEpoch(blockOffset, epochSize uint64) uint64 {
	epoch := math.Ceil(float64(blockOffset) / float64(epochSize))
	return uint64(epoch)
}
