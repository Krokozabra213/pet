package pkg

import "fmt"

func RenameKeyMap[K comparable, V any](m map[K]V, oldKey, newKey K) error {
	if oldKey == newKey {
		return nil
	}
	if _, exist := m[newKey]; exist {
		return fmt.Errorf("ключ %v уже существует", newKey)
	}
	if v, exist := m[oldKey]; exist {
		delete(m, oldKey)
		m[newKey] = v
	}
	return nil
}

func Hash(s string, max int) int {
	return int(s[0]) % max
}
