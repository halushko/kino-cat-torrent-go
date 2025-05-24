package helpers

func Byte2Gb(b int64) float64 {
	return float64(b) / (1024.0 * 1024.0 * 1024.0)
}
