package utils

func ReverseBytes(a []byte) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

func ReversedBytes(a []byte) []byte {
	ReverseBytes(a)
	return a
}
