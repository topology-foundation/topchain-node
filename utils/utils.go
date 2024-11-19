package utils

const EPOCH_SIZE = 10

func ConvertBlockToEpoch(block int64) uint64 {
	epoch := (block + EPOCH_SIZE - 1) / EPOCH_SIZE
	return uint64(epoch)
}
