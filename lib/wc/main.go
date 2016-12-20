package wc

import (
	"bufio"
	"io"
	"unicode/utf8"
)

func WordCountInString(text string) (count, byte_count, invalid int) {
	b := []byte(text)

	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			invalid += size
		} else {
			count++
			byte_count += size
		}
		b = b[size:]
	}
	return
}

// ファイルを読み取って文字数、バイト数、不正バイト数を数えて返す。
func WordCountInFile(rd io.Reader) (count, byte_count, invalid int) {
	in := bufio.NewReader(rd)

	for {
		r, size, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if r == utf8.RuneError {
			invalid += size
		} else {
			byte_count += size
			// 正常な文字だけカウントする。
			count++
		}
	}
	return
}
