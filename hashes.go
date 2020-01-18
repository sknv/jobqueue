package jobqueue

import (
	"hash/fnv"
)

// FNV calculates an integer FNV hash.
func FNV(str string) int {
	h := fnv.New32a()
	h.Write([]byte(str))
	return int(h.Sum32())
}
