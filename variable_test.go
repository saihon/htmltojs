package htmltojs

import (
	"testing"
)

func TestToString(t *testing.T) {
	v := Variable{Prefix: "_", name: []byte{'a'}}
	expect := "_a"
	actual := v.toString()
	if actual != expect {
		t.Errorf("\ngot : %v, want: %v\n", actual, expect)
	}
}

func TestGenName(t *testing.T) {
	v := Variable{}

	for i := 0; i < 25; i++ {
		actual := v.genName()
		expect := string(i + _A)
		if actual != expect {
			t.Fatalf("\ngot : %#v, want: %#v\n", actual, expect)
		}
	}

	actual := v.genName()
	expect := "z"
	if actual != expect {
		t.Fatalf("\ngot : %#v, want: %#v\n", actual, expect)
	}

	l := 730
	data := map[int]string{
		0:   "aa",
		1:   "ab",
		25:  "az",
		26:  "ba",
		51:  "bz",
		52:  "ca",
		77:  "cz",
		675: "zz",
		676: "aaa",
		701: "aaz",
		702: "aba",
		727: "abz",
		728: "aca",
		729: "acb",
		730: "acc",
	}
	for i := 0; i < l; i++ {
		actual := v.genName()
		expect, ok := data[i]
		if ok {
			if actual != expect {
				t.Fatalf("\ngot : %#v, want: %#v\n", actual, expect)
			}
		}
	}
}
