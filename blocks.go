package gobitcointools

import (
	"fmt"
	"reflect"
	"strings"
)

// ReverseString - output the reverse string of a given string s
func ReverseString(s string) string {

	strLen := len(s)

	// The reverse of a empty string is a empty string
	if strLen == 0 {
		return s
	}

	// Same above
	if strLen == 1 {
		return s
	}

	// Convert s into unicode points
	r := []rune(s)

	// Last index
	rLen := len(r) - 1

	// String new home
	rev := []string{}

	for i := rLen; i >= 0; i-- {
		rev = append(rev, string(r[i]))
	}

	return strings.Join(rev, "")
}

func serialize_header(inp string) []byte {
	o := encode(int(inp["version"]), 256, 4)[:] + string(inp["prevhash"])[:] + string(inp["merkle_root"])[:] + encode(inp["timestamp"], 256, 4)[:] + encode(inp["bits"], 256, 4)[:] + encode(inp["nonce"], 256, 4)[:]
	h := []byte(bin_sha256(bin_sha256(o))[:])
	if !reflect.DeepEqual(h, inp["hash"]) {
		panic(fmt.Errorf("AssertionError: %v", [2]interface{}{sha256(o), inp["hash"]}))
	}
	return []byte(o)
}
