layout: true
class: center, middle, footer

---

![:scaleImg 100%](/static/internal/img/presla-logo-black.svg "PresLa")

---

## Url parameters

.single-center-column-fit.leftalign[
To load your presentation, add the file name without extension to the url.

For example to load the presentation `foo.md` from your `md` folder, use:

```
http://localhost:8080/foo
```

> If you've altered the server port, use yours.

]

---

.left-column.leftalign[
## Creating new Presentations
]
.right-column.leftalign[
- Edit the `config.toml` and create some markdown files according to the `README.md`
- At the end of the presentation you'll find a list of your presentation files
- Fill in markdown according to the [wiki](https://github.com/gnab/remark/wiki)
]

---

class: top

.single-center-column-fit[

## Here, execute some php code!

php needs to be installed and the examples need a `tmp` dir in your root.

<div class="editor" data-filename="/tmp/main.php" data-executor="php"><?php

echo time()."\n";

sleep(1);

echo "foo\n";

sleep(3);

echo "bar\n";
</div>
]

---

class: top

.single-center-column-fit[

## Different theme? Different theme

You can set this via data attributes

<div class="editor" data-filename="/tmp/main.php" data-executor="php" data-theme="solarized_light"><?php

echo time()."\n";

sleep(1);

echo "foo\n";

sleep(3);

echo "bar\n";
</div>
]

---

class: top

.single-center-column-fit[

## Other languages? Sure!

You can easily write own executors with javascript

<div class="editor" data-filename="/tmp/main.go" data-executor="go">package main

import (
    "fmt"
    "os/user"
)

func main() {
    fmt.Println("Hello go!")
    usr, _ := user.Current()
    fmt.Println("Hello " + usr.Name + "!")
}
</div>
]

---

# Customizing

For creating your own theme, static files and templates, take a look at the [repository](https://git.3stadt.com/3stadt/presla).

---

## available presentations

.single-center-column-fit.leftalign[

{{range .Presentations}}
- [{{.}}](/{{.}})
{{end}}

]