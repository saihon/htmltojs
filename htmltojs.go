package htmltojs

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	_AttrVarNameKey = "variable-name"
)

func toCamelCase(s, sep string) string {
	a := strings.Split(s, sep)
	for i := 1; i < len(a); i++ {
		a[i] = strings.Title(a[i])
	}
	return strings.Join(a, "")
}

// HTMLtoJS is buffers and some settings implementation
type HTMLtoJS struct {
	// buffering converted javascript string
	Buffer bytes.Buffer

	Variable Variable

	Ignores  Filter
	Excludes Filter
	Includes Filter

	// A string of elements to append if there is no parent node.
	// such as "document.body", "document.documentElement"
	DefaultParent string
}

// New returns new HTMLtoJS with default settings
func New() *HTMLtoJS {
	return &HTMLtoJS{
		Variable: Variable{
			Prefix: "_",
		},
		Ignores:       Ignores{&Data{atom.Head}},
		Excludes:      Excludes{&Data{atom.Html, atom.Head, atom.Body, atom.Script}},
		Includes:      Includes{&Data{}},
		DefaultParent: "document.body",
	}
}

func (h *HTMLtoJS) writeFormat(format string, args ...interface{}) {
	h.Buffer.WriteString(fmt.Sprintf(format, args...))
}

func (h *HTMLtoJS) writeLineBreak() {
	h.Buffer.WriteString("\n")
}

var pattern = regexp.MustCompile(`[\r\n\t]|\s+|\\"|"`)

func replaceFunc(s string) string {
	switch s {
	case "\n", "\r":
		return ""
	case "\\\"", "\"":
		return "\\\""
	case "\t":
		fallthrough
	default:
		return " "
	}
}

func escapeText(s string) string {
	return pattern.ReplaceAllStringFunc(s, replaceFunc)
}

func (h *HTMLtoJS) writeCreateTextNode(varName, text string, node *html.Node) {
	// write text node without using createTextNode if tag name is script
	if node.Parent.DataAtom == atom.Script {
		h.Buffer.WriteString(text)
		h.writeLineBreak()
		return
	}

	parentVarName := h.getParentVarName(node)
	if parentVarName == h.DefaultParent {
		return
	}
	text = escapeText(text)
	if text == "" {
		return
	}
	h.writeFormat("var %s = document.createTextNode(\"%s\");", varName, text)
	h.writeLineBreak()
	h.writeFormat("%s.appendChild(%s);", parentVarName, varName)
	h.writeLineBreak()
}

func (h *HTMLtoJS) writeAttrDataset(varName string, attr html.Attribute) {
	key := toCamelCase(attr.Key[5:], "-")
	h.writeFormat("%s.dataset.%s = \"%s\";", varName, key, escapeText(attr.Val))
	h.writeLineBreak()
}

func (h *HTMLtoJS) writeAttrStyles(varName, style string) {
	styles := strings.Split(style, ";")
	for _, v := range styles {
		a := strings.Split(strings.TrimSpace(v), ":")
		if len(a) > 1 {
			property := toCamelCase(strings.TrimSpace(a[0]), "-")
			value := strings.TrimSpace(a[1])
			h.writeFormat("%s.style.%s = \"%s\";", varName, property, value)
			h.writeLineBreak()
		}
	}
}

func (h *HTMLtoJS) writeClassName(varName, className string) {
	h.writeFormat("%s.className = \"%s\";", varName, className)
	h.writeLineBreak()
}

var events = [61]string{
	"onabort",
	"onafterprint",
	"onbeforeprint",
	"onbeforeunload",
	"onblur",
	"oncancel",
	"oncanplay",
	"oncanplaythrough",
	"onchange",
	"onclick",
	"oncuechange",
	"ondbclick",
	"ondurationchange",
	"onemptied",
	"onended",
	"onerror",
	"onfocus",
	"onhashchange",
	"oninput",
	"oninvalid",
	"onkeydown",
	"onkeypress",
	"onkeyup",
	"onload",
	"onloadeddata",
	"onloadedmetadata",
	"onloadstart",
	"onmessage",
	"onmousedown",
	"onmouseenter",
	"onmouseleave",
	"onmousemove",
	"onmouseout",
	"onmouseover",
	"onmouseup",
	"onmousewheel",
	"onoffline",
	"ononline",
	"onpagehide",
	"onpageshow",
	"onpause",
	"onplay",
	"onplaying",
	"onpopstate",
	"onprogress",
	"onratechange",
	"onresize",
	"onreset",
	"onscroll",
	"onseeked",
	"onseeking",
	"onselect",
	"onshow",
	"onstalled",
	"onstorage",
	"onsubmit",
	"onsuspend",
	"ontimeupdate",
	"onunload",
	"onvolumechange",
	"onwaiting",
}

func isEvent(key string) bool {
	if !strings.HasPrefix(key, "on") {
		return false
	}
	key = strings.ToLower(key)
	for _, event := range events {
		if key == event {
			return true
		}
	}
	return false
}

func (h *HTMLtoJS) writeAddEventListener(varName, key, value string) {
	h.writeFormat("%s.addEventListener(\"%s\", function() {%s}, false);", varName, key[2:], value)
	h.writeLineBreak()
}

func (h *HTMLtoJS) writeAttribute(varName, key, value string) {
	key = toCamelCase(key, "-")
	if value == "" {
		h.writeFormat("%s.%s = true;", varName, key)
	} else {
		h.writeFormat("%s.%s = \"%s\";", varName, key, value)
	}
	h.writeLineBreak()
}

func (h *HTMLtoJS) writeAttributes(varName string, node *html.Node) {
loop:
	for _, attr := range node.Attr {
		switch attr.Key {
		case _AttrVarNameKey:
			continue loop
		case "style":
			h.writeAttrStyles(varName, attr.Val)
		case "class":
			h.writeClassName(varName, attr.Val)
		default:
			switch {
			case strings.HasPrefix(attr.Key, "data-"):
				h.writeAttrDataset(varName, attr)
			default:
				if isEvent(attr.Key) {
					h.writeAddEventListener(varName, attr.Key, attr.Val)
				} else {
					h.writeAttribute(varName, attr.Key, attr.Val)
				}
			}
		}
	}
}

func (h *HTMLtoJS) writeCreateElement(varName string, node *html.Node) {
	h.writeFormat("var %s = document.createElement(\"%s\");", varName, node.Data)
	h.writeLineBreak()
}

func (h *HTMLtoJS) getParentVarName(node *html.Node) string {
	varName := h.DefaultParent
	if node.Parent == nil {
		return varName
	}
	for _, a := range node.Parent.Attr {
		if a.Key == _AttrVarNameKey {
			varName = a.Val
			break
		}
	}
	return varName
}

func (h *HTMLtoJS) writeAppendChild(varName string, node *html.Node) {
	parentVarName := h.getParentVarName(node)
	h.writeFormat("%s.appendChild(%s);", parentVarName, varName)
	h.writeLineBreak()
}

func (h *HTMLtoJS) excludeMatch(n *html.Node) bool {
	if h.Includes.len() > 0 {
		return true
	}
	return h.Excludes.match(n)
}

func (h *HTMLtoJS) setVarNameToAttribute(varName string, node *html.Node) {
	node.Attr = append(node.Attr, html.Attribute{Key: _AttrVarNameKey, Val: varName})
}

func (h *HTMLtoJS) parseNode(node *html.Node) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			if text := strings.TrimSpace(c.Data); text != "" {
				varName := h.Variable.genName()
				h.writeCreateTextNode(varName, text, c)
			}

		} else if c.Type == html.ElementNode {
			if h.Includes.match(c) || !h.excludeMatch(c) {
				varName := h.Variable.genName()
				h.setVarNameToAttribute(varName, c)
				h.writeCreateElement(varName, c)
				h.writeAttributes(varName, c)
				h.writeAppendChild(varName, c)
			}

			if !h.Ignores.match(c) {
				h.parseNode(c)
			}
		}
	}
}

// Parse parses the HTML data and buffers results.
func (h *HTMLtoJS) Parse(r io.Reader) error {
	node, err := html.Parse(r)
	if err != nil {
		return err
	}
	h.parseNode(node)
	return nil
}

// String returns string data from buffer.
func (h *HTMLtoJS) String() string {
	return h.Buffer.String()
}

// WriteTo writes data to the specified writer.
func (h *HTMLtoJS) WriteTo(w io.Writer) (int64, error) {
	return h.Buffer.WriteTo(w)
}
