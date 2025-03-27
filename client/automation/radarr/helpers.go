package radarr

// Helper function to convert []int64 to []int32
func convertInt64SliceToInt32(in []int64) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}
