package wc

// テストを実行するとき
// $ go test github.com/oyakata/kanjo/lib/wc/

import (
	"strings"
	"fmt"
	"testing"
)

type T struct {
	Input                     string
	Count, ByteCount, Invalid int
}

var valid, broken, exp string

func init() {
	valid = string([]byte{
		240, 169, 184, 189, // 𩸽
		227, 129, 174, // の
		227, 129, 178, // ひ
		227, 130, 137, // ら
		227, 129, 141, // き
		240, 159, 153, 134, // face with ok gesture
		240, 159, 153, 134, // face with ok gesture
		240, 159, 153, 134, // face with ok gesture
	})

	broken = string([]byte{
		169, 184, 189, // [NG] 先頭欠損
		227, 129, 174, // の
		129, 178, // [NG] 先頭欠損
		227, 137, // [NG] 2桁目欠損
		227, 129, // [NG] 3桁目欠損
		240, 153, 134, // [NG] 2桁目欠損
		240, 159, 134, // [NG] 3桁目欠損
		240, 159, 153, // [NG] 末尾欠損
	})

	part := "𩸽のひらきを居酒屋で注文して、1時間経つがまだ来ない。𠈻な客が店員を引き止めてなじるからだ。"
	// 46 * 3 + 7 + 2 = 147文字
	// 1行: (4byte * 2文字) + (3byte * 43文字) + (1byte * 1文字) = 138byte
	// 全体: 140 * 3 + 7 + 6 = 427byte
	exp = fmt.Sprintf("%v\n\n%v\n\n\n%v\n\n以上", part, part, part)
}

func TestWordCountInString(t *testing.T) {
	cases := []T{
		{valid, 8, 28, 0},
		{broken, 1, 3, 18},
		{"", 0, 0, 0},
		{exp, 147, 427, 0},
	}

	for _, tc := range cases {
		x, y, z := WordCountInString(tc.Input)
		if x != tc.Count || y != tc.ByteCount || z != tc.Invalid {
			t.Errorf("WordCountInString=%v% v% v, want=%v% v% v",
				x, y, z,
				tc.Count, tc.ByteCount, tc.Invalid,
			)
		}
	}
}

func TestWordCountInFile(t *testing.T) {
	cases := []T{
		{valid, 8, 28, 0},
		{broken, 1, 3, 18},
		{"", 0, 0, 0},
		{exp, 147, 427, 0},
	}

	for _, tc := range cases {
		x, y, z := WordCountInFile(strings.NewReader(tc.Input))
		if x != tc.Count || y != tc.ByteCount || z != tc.Invalid {
			t.Errorf("WordCountInString=%v% v% v, want=%v% v% v",
				x, y, z,
				tc.Count, tc.ByteCount, tc.Invalid,
			)
		}
	}
}
