package htmltojs

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestFilter(t *testing.T) {
	expect := []atom.Atom{atom.H1, atom.Div, atom.P}

	d := new(Data)
	d.Set(expect...)
	actual := []atom.Atom(*d)
	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("\nSet - got : %#v, want: %#v\n", actual, expect)
	}

	for i := len(expect) - 1; i >= 0; i-- {
		d.Del(expect[i])
		if d.len() != i {
			t.Fatalf("\nDel - got : %d, want: %d\n", d.len(), i)
		}
	}

	for _, v := range expect {
		d.Add(v)
	}
	actual = []atom.Atom(*d)
	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("\nAdd - got : %#v, want: %#v\n", actual, expect)
	}

	data := []struct {
		node   *html.Node
		expect bool
	}{
		{&html.Node{DataAtom: atom.H1}, true},
		{&html.Node{DataAtom: atom.Div}, true},
		{&html.Node{DataAtom: atom.P}, true},
		{&html.Node{DataAtom: atom.Span}, false},
		{&html.Node{DataAtom: atom.A}, false},
		{&html.Node{DataAtom: atom.H3}, false},
	}

	for i, v := range data {
		actual := d.match(v.node)
		if actual != v.expect {
			t.Fatalf("\nMatch %d - got : %#v, want: %#v\n", i, actual, v.expect)
		}
	}
}
