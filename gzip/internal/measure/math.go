package measure

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bucketSize(samples, buckets int) int {
	size := samples / buckets

	if samples%buckets > 0 {
		size++
	}

	return size
}
