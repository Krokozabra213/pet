package utils

import "math/bits"

// func RenameKeyMap[K comparable, V any](m map[K]V, oldKey, newKey K) error {
// 	if oldKey == newKey {
// 		return nil
// 	}
// 	if _, exist := m[newKey]; exist {
// 		return fmt.Errorf("ключ %v уже существует", newKey)
// 	}
// 	if v, exist := m[oldKey]; exist {
// 		delete(m, oldKey)
// 		m[newKey] = v
// 	}
// 	return nil
// }

// временная реализация
func Hash(s string, max int) int {
	return int(s[0]) % max
}

func SimpleUint64Hash(key uint64, seed uint64) uint64 {
	return key*0x9e3779b97f4a7c15 ^ seed
}

func IsPowerOfTwo(n uint8) bool {
	return n != 0 && (n&(n-1)) == 0
}

func LogarithmFloor(n uint8) uint8 {
	if n == 0 {
		return 0
	}
	return uint8(bits.Len(uint(n)) - 1)
}

func PowInt(x, y uint8) uint8 {

	result := x
	for i := uint8(1); i < y; i++ {
		result *= x
	}
	return result
}
