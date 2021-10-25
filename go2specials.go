package gobitcointools

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var CodeStrings = map[int]string{
	2:  "01",
	10: "0123456789",
	16: "0123456789abcdef",
	32: "abcdefghijklmnopqrstuvwxyz234567",
	58: "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz",
	256: strings.Join(func() (elts []string) {
		for x := 0; x < 256; x++ {
			elts = append(elts, string(rune(x)))
		}
		return
	}(), ""),
}

func BinDblSha256(s string) []byte {
	bytes_to_hash := FromStringToBytes(s)
	key1 := hmac.New(sha256.New, bytes_to_hash)
	bytes, _ := hex.DecodeString(hex.EncodeToString(key1.Sum(nil)))
	key1 = hmac.New(sha256.New, bytes)
	hex, _ := hex.DecodeString(hex.EncodeToString(key1.Sum(nil)))
	return hex
}

func Lpad(msg, symbol string, length int) string {
	if len([]rune(msg)) >= length {
		return msg
	}
	msg2 := strings.Repeat(symbol, (length - utf8.RuneCountInString(msg)))
	return msg2 + msg
}

func GetCodeString(base int) string {
	if _, ok := CodeStrings[base]; ok {
		return CodeStrings[base]
	} else {
		panic(fmt.Errorf("ValueError: %v", "Invalid base!"))
	}
}

func ChangeBase(str string, frm, to, minlen int) string {
	//minlen = 0
	if frm == to {
		return Lpad(str, GetCodeString(frm), minlen)
	}
	return encode(decode(str, frm), to, minlen)
}

func BinToB58check(inp string, magicbyte byte) string {
	inpFmtd := string(rune(int(magicbyte))) + inp
	var leadingzbytes int
	var checksum []byte

	re, err := regexp.Compile("^\u0000*")
	if err != nil {
		log.Fatal(err)
	}
	for _, word := range inpFmtd {
		found := re.MatchString(string(word))
		if found {
			leadingzbytes = len(string(word))
			checksum = BinDblSha256(inpFmtd)[:4]
		}
	}
	//msg2 := strings.Repeat(symbol, (length - utf8.RuneCountInString(msg)))
	return strings.Repeat("1", int(leadingzbytes)) + ChangeBase(inpFmtd+string(checksum), 256, 58, 0)
}

func BytesToHexString(b []byte) string {
	encodedString := hex.EncodeToString(b)
	return encodedString
}

func SafeFromHex(s string) string {
	bs, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func FromIntRepresentationToBytes(a int) string {
	return fmt.Sprintf("%v", a)
}

func FromIntToByte(a int) rune {
	/*z_bytes := make([]byte, 32)
	binary.BigEndian.PutUint64(z_bytes, uint64(a))
	return z_bytes*/

	return rune(a)
}

func FromByteToInt(a byte) int {
	return int(a)
}

func FromBytesToString(s byte) string {
	return string(s)
}

func FromStringToBytes(a string) []byte {
	return []byte(a)
}

func SafeHexlify(a int) string {
	s := fmt.Sprintf("%x", a)
	ui, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return "error"
	}
	return fmt.Sprintf("%x", ui)
}

// Max returns the larger of x or y.
func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func encode(val, base, minlen int) string {
	//minlen = 0
	base, minlen = int(base), int(minlen)
	CodeString := GetCodeString(base)
	result := " "
	for val > 0 {
		result = string(CodeString[val%base]) + result
		val /= base
	}
	runes := []rune(result)
	return strings.Repeat(string(CodeString[0]), int(Max(int64((minlen-len(runes))), 0))) + result
}

func decode(str string, base int) int {
	base = int(base)
	CodeString := GetCodeString(base)
	result := 0
	if base == 16 {
		str = strings.ToLower(str)
		str = string(str)
	}
	runes := []rune(str)
	for len(runes) > 0 {
		result *= base
		result += strings.Index(CodeString, string(str[0]))
		str = str[1:]
	}
	return result
}

//There are many ways to do this
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(x int) string {
	b := make([]byte, x)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
