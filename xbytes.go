package xbytes

import (
	"unicode/utf8"
)

type ASCIISet [8]uint32

func (as *ASCIISet) contains(c byte) bool {
	return (as[c>>5] & (1 << uint(c&31))) != 0
}

func makeASCIISet(chars string) (as ASCIISet, ok bool) {
	for i := 0; i < len(chars); i++ {
		c := chars[i]
		if c >= utf8.RuneSelf {
			return as, false
		}
		as[c>>5] |= 1 << uint(c&31)
	}
	return as, true
}

func MustMakeASCIISet(chars string) ASCIISet {
	as, ok := makeASCIISet(chars)
	if !ok {
		panic("non ascii")
	}
	return as
}

func TrimRightByte(s []byte, b rune) []byte {
	i := lastIndexFunc(s, b)
	if i >= 0 && s[i] >= utf8.RuneSelf {
		_, wid := utf8.DecodeRune(s[i:])
		i += wid
	} else {
		i++
	}
	return s[0:i]
}

func TrimRightASCIISet(s []byte, b ASCIISet) []byte {
	i := lastIndexFuncASCIISet(s, b)
	if i >= 0 && s[i] >= utf8.RuneSelf {
		_, wid := utf8.DecodeRune(s[i:])
		i += wid
	} else {
		i++
	}
	return s[0:i]
}

func TrimLeftASCIISet(s []byte, f ASCIISet) []byte {
	i := indexFuncASCIISet(s, f)
	if i == -1 {
		return nil
	}
	return s[i:]
}

func TrimASCIISet(s []byte, b ASCIISet) []byte {
	return TrimLeftASCIISet(TrimRightASCIISet(s, b), b)
}

func TrimLeftByte(s []byte, f rune) []byte {
	i := indexFunc(s, f)
	if i == -1 {
		return nil
	}
	return s[i:]
}

func TrimByte(s []byte, b rune) []byte {
	return TrimLeftByte(TrimRightByte(s, b), b)
}

func indexFunc(s []byte, f rune) int {
	start := 0
	for start < len(s) {
		wid := 1
		r := rune(s[start])
		if r >= utf8.RuneSelf {
			r, wid = utf8.DecodeRune(s[start:])
		}
		if f != r {
			return start
		}
		start += wid
	}
	return -1
}

func indexFuncASCIISet(s []byte, f ASCIISet) int {
	start := 0
	for start < len(s) {
		wid := 1
		r := rune(s[start])
		if r >= utf8.RuneSelf {
			r, wid = utf8.DecodeRune(s[start:])
		}
		if !(r < utf8.RuneSelf && f.contains(byte(r))) {
			return start
		}
		start += wid
	}
	return -1
}

func lastIndexFunc(s []byte, f rune) int {
	for i := len(s); i > 0; {
		r, size := rune(s[i-1]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeLastRune(s[0:i])
		}
		i -= size
		if f != r {
			return i
		}
	}
	return -1
}

func lastIndexFuncASCIISet(s []byte, f ASCIISet) int {
	for i := len(s); i > 0; {
		r, size := rune(s[i-1]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeLastRune(s[0:i])
		}
		i -= size
		if !(r < utf8.RuneSelf && f.contains(byte(r))) {
			return i
		}
	}
	return -1
}
