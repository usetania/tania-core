package mathhelper

// EPSILON is used because float equality and zeroness is unpredictable
const EPSILON = 0.0001

// IsEqual is used to calculate
func IsEqual(a, b float32) bool {
	if a-b < 0 {
		return (a-b)*-1 <= EPSILON
	}

	return a-b <= EPSILON
}
