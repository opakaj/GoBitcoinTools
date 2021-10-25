package gobitcointools

import (
	"encoding/hex"
	"reflect"
	"strings"
)

func bintob58check(inp string, magicbyte byte) string {
	inp_fmtd := string(rune((int(magicbyte)))) + inp
	leadingzbytes := 0
	for _, x := range inp_fmtd {
		if x != 0 {
			break
		}
		leadingzbytes += 1
	}
	checksum := BinDblSha256(inp_fmtd)[:4]
	return strings.Repeat("1", leadingzbytes) + ChangeBase(inp_fmtd+string(checksum), 256, 58, 0)
}

func safeFromHex(s string) []byte {
	bs, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	b := []byte(bs)
	return b
}

func fromIntRepresentationToBytes(a int) []byte {
	b := []byte(string(rune(a)))
	return b
}

func fromStringToBytes(a interface{}) []byte {
	if func() bool {
		switch a.(type) {
		case byte:
			return true
		}
		return false
	}() {
		return a.([]byte)
	} else {
		b := []byte(a.(string))
		return b
	}
}

func fromIntToByte(a int) byte {
	return byte(a)
}

func fromByteToInt(a byte) int {
	return int(a)
}

func Encode(val, base, minlen int) string {
	base, minlen = int(base), int(minlen)
	codeString := GetCodeString(base)
	resultBytes := make([]byte, 32)
	for val > 0 {
		curcode := codeString[val%base]
		resultBytes = append(resultBytes, byte(int(curcode)))
		val /= base
	}

	padSize := minlen - len(resultBytes)

	paddingElement := func() []byte {
		if reflect.DeepEqual(base, 256) {
			return []byte("\u0000")
		}
		return func() []byte {
			if reflect.DeepEqual(base, 58) {
				return []byte("1")
			}
			return []byte("0")
		}()
	}()
	if padSize > 0 {
		resultBytes = func(repeated []byte, n int) (result []byte) {
			for i := 0; i < n; i++ {
				result = append(result, repeated...)
			}
			return result
		}(paddingElement, padSize)
		resultBytes = append(resultBytes, resultBytes...)
	}
	resultString := strings.Join(func() (elts []string) {
		for _, y := range resultBytes {
			elts = append(elts, string(rune(y)))
		}
		return
	}(), "")
	result := func() string {
		if reflect.DeepEqual(base, 256) {
			return string(resultBytes)
		}
		return resultString
	}()
	return result
}

func Decode(str interface{}, base int) int {
	var extract func(d interface{}, cs string) interface{}
	if reflect.DeepEqual(base, 256) && func() bool {
		switch str.(type) {
		case string:
			return true
		}
		return false
	}() {
		str = []byte(str.(string))
	}

	base = int(base)
	codeString := GetCodeString(base)
	result := 0
	if base == 256 {
		extract = func(d interface{}, cs string) interface{} {
			return d
		}
	} else {
		extract = func(d interface{}, cs string) interface{} {
			strings.Index(cs, func() string {
				if func() bool {
					switch d.(type) {
					case string:
						return true
					}
					return false
				}() {
					return d.(string)
				}
				return string(d.(rune))
			}())
			return d.(string)
		}
	}
	if base == 16 {
		str = strings.ToLower(str.(string))
	}
	for len([]rune(str.(string))) > 0 {
		result *= base
		result += extract(int(str.(string)[0]), codeString).(int)
		str = str.(string)[1:]
	}
	return result
}
