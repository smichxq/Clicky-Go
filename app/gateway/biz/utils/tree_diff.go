package utils

import (
	"sort"

	"github.com/cespare/xxhash"
)

// Calculate key hash for a map[string]uint64
//
// Notice: no need to sort the keys of the map
func HashKeys(set map[string]uint64) uint64 {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := xxhash.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte{0})
	}
	return h.Sum64()
}
