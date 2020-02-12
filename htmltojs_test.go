package htmltojs

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestToCamelCase(t *testing.T) {
	data := map[string]string{
		"foo":         "foo",
		"foo-bar":     "fooBar",
		"foo-bar-baz": "fooBarBaz",
	}

	for s, expect := range data {
		actual := toCamelCase(s, "-")
		if actual != expect {
			t.Errorf("\ngot : %v, want: %v\n", actual, expect)
		}
	}
}

func TestDefaultValues(t *testing.T) {
	h := New()

	expect := "document.body"
	actual := h.DefaultParent
	if actual != expect {
		t.Errorf("\nDefaultParent\ngot : %#v, want: %#v\n", actual, expect)
	}

	expect = "_"
	actual = h.Variable.Prefix
	if actual != expect {
		t.Errorf("\nVariable.prefix\ngot : %#v, want: %#v\n", actual, expect)
	}
}

func TestWriteFormat(t *testing.T) {
	h := new(HTMLtoJS)
	expect := "hello world"
	h.writeFormat("%s", expect)
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot : %v, want: %v\n", actual, expect)
	}
}

func TestWriteLineBreak(t *testing.T) {
	h := new(HTMLtoJS)
	expect := []byte("\n")
	h.writeLineBreak()
	actual := h.Buffer.Bytes()
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("\ngot : %#v, want: %#v\n", actual, expect)
	}
}

func TestEscapeText(t *testing.T) {
	s := `"hello"
	world
`
	expect := "\\\"hello\\\" world"
	actual := escapeText(s)
	if actual != expect {
		t.Errorf("\ngot : %v, want: %v\n", actual, expect)
	}
}

func TestGetParentVarName(t *testing.T) {
	h := new(HTMLtoJS)
	expect := "document.body"
	h.DefaultParent = expect
	actual := h.getParentVarName(&html.Node{})
	if actual != expect {
		t.Errorf("\ngot : %#v, want: %#v\n", actual, expect)
	}

	expect = "parent"
	parent := &html.Node{Attr: []html.Attribute{{Key: _AttrVarNameKey, Val: expect}}}
	node := &html.Node{Parent: parent}
	actual = h.getParentVarName(node)
	if actual != expect {
		t.Errorf("\ngot : %#v, want: %#v\n", actual, expect)
	}
}

func TestWriteCreateTextNode(t *testing.T) {
	varName := "b"

	h := new(HTMLtoJS)
	node := &html.Node{}
	node.Parent = &html.Node{Attr: []html.Attribute{{Key: _AttrVarNameKey, Val: "a"}}}
	text := "hello"
	h.writeCreateTextNode(varName, text, node)
	expect := `var b = document.createTextNode("hello");
a.appendChild(b);
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%#v\nwant:\n%#v\n", actual, expect)
	}

}

func TestWriteAttrDataset(t *testing.T) {
	h := new(HTMLtoJS)
	h.writeAttrDataset("a", html.Attribute{Key: "data-name", Val: "value"})
	expect := `a.dataset.name = "value";
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}

func TestWriteAttrStyles(t *testing.T) {
	h := new(HTMLtoJS)
	h.writeAttrStyles("a", "background-color: red; border: 1px solid blue; width: 100px; height: 100px;")
	expect := `a.style.backgroundColor = "red";
a.style.border = "1px solid blue";
a.style.width = "100px";
a.style.height = "100px";
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}

func TestWriteClassName(t *testing.T) {
	h := new(HTMLtoJS)
	h.writeClassName("a", "className")
	expect := `a.className = "className";
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}

func TestWriteAddEventListener(t *testing.T) {
	h := new(HTMLtoJS)
	h.writeAddEventListener("a", "onclick", "click();")
	expect := `a.addEventListener("click", function() {click();}, false);
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}

func TestWriteAttribute(t *testing.T) {
	h := new(HTMLtoJS)
	h.writeAttribute("a", "id", "value")
	h.writeAttribute("b", "checked", "")
	expect := `a.id = "value";
b.checked = true;
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}

func TestWriteAttributes(t *testing.T) {
	h := new(HTMLtoJS)
	node := &html.Node{Attr: []html.Attribute{
		{Key: "class", Val: "classname"},
		{Key: "type", Val: "text"},
		{Key: "disabled", Val: ""},
	}}
	h.writeAttributes("a", node)
	expect := `a.className = "classname";
a.type = "text";
a.disabled = true;
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}

func TestWriteCreateElement(t *testing.T) {
	h := new(HTMLtoJS)
	node := &html.Node{Data: "div"}
	h.writeCreateElement("a", node)
	expect := `var a = document.createElement("div");
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}

func TestAppendChild(t *testing.T) {
	h := new(HTMLtoJS)
	h.DefaultParent = "document.body"
	h.writeAppendChild("a", &html.Node{})
	expect := `document.body.appendChild(a);
`
	actual := h.Buffer.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}

func TestSetVarNameToAttribute(t *testing.T) {
	h := new(HTMLtoJS)
	expect := "a"
	node := &html.Node{}
	h.setVarNameToAttribute(expect, node)
	attr := node.Attr[0]
	if attr.Key != _AttrVarNameKey {
		t.Errorf("\ngot : %#v, want: %#v\n", attr.Key, _AttrVarNameKey)
	}

	if attr.Val != expect {
		t.Errorf("\ngot : %#v, want: %#v\n", attr.Val, expect)
	}
}

func TestParse(t *testing.T) {
	h := New()

	s := `
<div id="content">
content text node 1
<ul style="list-style-type: none;">
	<li class="list">foo</li>
	<li class="list">bar</li>
	<li class="list">baz</li>
</ul>
content text node 2
</div>
<script>
(function () {
	console.log('hello world');
})();
</script>
`

	if err := h.Parse(strings.NewReader(s)); err != nil {
		t.Fatalf("error: %s\n", err)
	}

	expect := `var _a = document.createElement("div");
_a.id = "content";
document.body.appendChild(_a);
var _b = document.createTextNode("content text node 1");
_a.appendChild(_b);
var _c = document.createElement("ul");
_c.style.listStyleType = "none";
_a.appendChild(_c);
var _d = document.createElement("li");
_d.className = "list";
_c.appendChild(_d);
var _e = document.createTextNode("foo");
_d.appendChild(_e);
var _f = document.createElement("li");
_f.className = "list";
_c.appendChild(_f);
var _g = document.createTextNode("bar");
_f.appendChild(_g);
var _h = document.createElement("li");
_h.className = "list";
_c.appendChild(_h);
var _i = document.createTextNode("baz");
_h.appendChild(_i);
var _j = document.createTextNode("content text node 2");
_a.appendChild(_j);
(function () {
	console.log('hello world');
})();
`
	actual := h.String()
	if actual != expect {
		t.Errorf("\ngot :\n%v\nwant:\n%v\n", actual, expect)
	}
}
