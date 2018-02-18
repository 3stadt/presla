layout: true
class: center, middle, footer

---

[![:scaleImg 100%](/static/internal/img/presla-logo-black.svg "Presla")](http://presla.io/)

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

## Creating new Presentations

.single-center-column-fit.leftalign[
- Your config file is: `{{ .ConfigFile }}`
- Edit the `MarkdownPath` and create markdown file according to the `README.md`
- Put markdown into the file according to the [wiki](https://github.com/gnab/remark/wiki)
  - Restart presla so it picks up your newly created file
- At the end of this presentation you'll find a list of your current presentation files
]

---

class: top

.single-center-column-fit[

## Here, execute some php code!

php needs to be installed and the examples need a `tmp` dir in your root.

<div class="editor" data-filename="{{ .TempDir }}/main.php" data-executor="php"><?php

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

<div class="editor" data-filename="{{ .TempDir }}/main.php" data-executor="php" data-theme="solarized_light"><?php

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

<div class="editor" data-filename="{{ .TempDir }}/main.go" data-executor="go" data-mode="golang">package main

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