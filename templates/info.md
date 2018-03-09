layout: true
class: center, middle, top, footer

---

[![:scaleImg 100%](/static/internal/img/presla-logo-black.svg "Presla")](http://presla.io/)

---

## Creating new Presentations

.single-center-column-fit.leftalign[
- Your config file is: `{{ .ConfigFile }}`
- Edit the `MarkdownPath` and create markdown file according to the `README.md`
- Put markdown into the file according to the [wiki](https://github.com/gnab/remark/wiki)
  - Restart presla so it picks up your newly created file
- At the end of this presentation you'll find a list of your current presentation files
]

---

## You can execute code with presla!
.single-center-column-fit.leftalign[

- All temporary files for the examples are created in this directory when you press execute:
  - ***`{{ .TempDir }}`***
- To insert an editor/logview in your presentations, create a div like this:

```golang
<div class="editor" data-filename="/some/dir/file.go" data-executor="go" data-mode="golang">
package main

import "fmt"

func main() {
    fmt.Println("Hello presla!")
}
</div>
```
]
.single-center-column-fit.leftalign[
> [PHP](http://php.net/manual/install.php) needs to be installed to execute php code!   
> [Golang](https://golang.org/doc/install) needs to be installed to execute go code!
]

---

## Here, execute some php code!

.single-center-column-fit[

<div class="editor" data-filename="{{ .TempDir }}/main.php" data-executor="php" data-theme="solarized_dark"><?php

echo time()."\n";

sleep(1);

echo "foo\n";

sleep(3);

echo "bar\n";
</div>
]

---

## Need to pass some arguments? Sure!

.single-center-column-fit[

<div class="editor" data-showcmd="true" data-filename="{{ .TempDir }}/main.php" data-executor="php" data-theme="solarized_dark"><?php
var_dump($argv);
</div>
]

---

## Errors show up in red

.single-center-column-fit[

<div class="editor" data-filename="{{ .TempDir }}/main.php" data-executor="php" data-theme="solarized_dark"><?php

echo foobar(); // undefined function

</div>
]

---

.single-center-column-fit[

## Different theme? Different theme

You can set this via data attributes

<div class="editor" data-filename="{{ .TempDir }}/main.php" data-executor="php" data-theme="solarized_light"><?php

echo time()."\n";

sleep(1);

echo "foo\n";

sleep(3);

echo "bar\n";
</div>
]

---

.single-center-column-fit[

## Other languages? Sure!

You can easily write own executors with javascript

<div class="editor" data-filename="{{ .TempDir }}/main.go" data-executor="go" data-mode="golang" data-theme="solarized_dark">package main

import (
    "fmt"
    "os/user"
)

func main() {
    fmt.Println("Hello go!")
    usr, _ := user.Current()
    fmt.Println("Hello " + usr.Username + "!")
}
</div>
]

---

## Even multiple executors are possible

You could show differences between language versions

.single-center-column-fit[
<div class="editor" data-filename="{{ .TempDir }}/main.go" data-executors="go;go-1.6;php" data-mode="golang" data-theme="solarized_dark">package main

import (
    "fmt"
    "runtime"
)

func main() {
    fmt.Println("Version: " + runtime.Version())
}</div>
]

---

## Whatever is installed on your system.. :)

.single-center-column-fit[
<div class="editor" data-filename="{{ .TempDir }}/main.py" data-executor="python" data-mode="python" data-theme="gob"># Want some python?
def greet(name):
    print('Hello', name + '!')

greet('python')
greet('world')
greet('presla')
</div>
]

---

# Customizing

For creating your own theme, static files and templates, take a look at the [repository](https://git.3stadt.com/3stadt/presla).

---

## Available presentations

.single-center-column-fit.leftalign[

{{range .Presentations}}
- [{{.}}](/{{.}})
{{end}}

]