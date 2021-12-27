package pkg

func CompareBytes(in, out []byte) float64 {
	var (
		errorsCount int
	)

	for idx := range in {
		if in[idx] != out[idx] {
			errorsCount++
		}
	}

	return float64(((len(in) - errorsCount) / len(in)) * 100)
}
