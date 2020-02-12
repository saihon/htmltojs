[![Actions Status](https://github.com/saihon/htmltojs/workflows/Go/badge.svg)](https://github.com/saihon/htmltojs/actions?query=workflow%3AGo)

## htmltojs

Convert HTML to JavaScript

<br/>

## install

```
go get github.com/saihon/htmltojs
```

<br/>

## example

```go

package main

import (
	"log"
	"os"
	"strings"

	htmltojs "github.com/saihon/htmltojs"
)

func main() {
	r := strings.NewReader(`
<div id="content">
	text node 1
	<ul style="list-style-type: none;">
	  <li class="list">foo</li>
	  <li class="list">bar</li>
	  <li class="list">baz</li>
	</ul>
	text node 2
</div>
<script>
(function () {
  console.log('hello world');
})();
</script>`)

	h := htmltojs.New()

	if err := h.Parse(r); err != nil {
		log.Fatal(err)
	}

	h.WriteTo(os.Stdout)
}

```

Output

```javascript

var _a = document.createElement("div");
_a.id = "content";
document.body.appendChild(_a);
var _b = document.createTextNode("text node 1");
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
var _j = document.createTextNode("text node 2");
_a.appendChild(_j);
(function () {
  console.log('hello world');
})();

```

<br/>
<br/>
