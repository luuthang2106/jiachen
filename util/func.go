package util

func ToPointer[V comparable](val V) *V {
	return &val
}

func Contains[K comparable](arr []K, v K) bool {
	for _, ele := range arr {
		if ele == v {
			return true
		}
	}
	return false
}
