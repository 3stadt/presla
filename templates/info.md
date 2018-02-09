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
- Take a look at the `css/theme-shopware.css`
]

---

## available presentations

.single-center-column-fit.leftalign[

{{range .Presentations}}
- [{{.}}](/{{.}})
{{end}}

]