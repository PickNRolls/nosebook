package boolean

func Xor(left bool, right bool) bool {
	return left && !right || !left && right
}
