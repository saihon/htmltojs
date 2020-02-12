## htmltojs

Convert HTML to JavaScript for command line

<br/>

## install

```
go get github.com/saihon/htmltojs/cmd/htmltojs
```

<br/>

## example

Pipe the output from cat command into htmltojs

```
$ cat index.html | htmltojs
```

Specifies the output file

```
$ cat index.html | htmltojs -o index.js
```

Specifies the input file

```
$ htmltojs -i index.html
```

Specifies the output file

```
$ htmltojs -i index.html -o index.js
```

Redirect

```
$ htmltojs << EOS
<ins class="adsbygoogle"
     style="display:block"
     data-ad-client="ca-pub-xxxxxxxx"
     data-ad-slot="xxxxxxxxx"
     data-ad-format="auto"
     data-full-width-responsive="true"></ins>
<script>
     (adsbygoogle = window.adsbygoogle || []).push({});
</script>
EOS

var _a = document.createElement("ins");
_a.className = "adsbygoogle";
_a.style.display = "block";
_a.dataset.adClient = "ca-pub-xxxxxxxx";
_a.dataset.adSlot = "xxxxxxxxx";
_a.dataset.adFormat = "auto";
_a.dataset.fullWidthResponsive = "true";
document.body.appendChild(_a);
(adsbygoogle = window.adsbygoogle || []).push({});

```

<br/>
<br/>
