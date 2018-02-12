# Presla

Presla (Presentation Lambda) is an application that runs on your computer. It creates a local webserver that is equipped with [remarkjs](https://remarkjs.com/), allowing you to create your own presentations in markdown.

Additionaly, it is configured so it can execute Code examples directly from within your presentation, using [ace editor](https://ace.c9.io/) and code executors written in JS. Everything can be completely customized.


## Getting Started

### Prerequisites

You'll need a x64 based Mac, Linux or Windows System and a Web Browser. Preferably Google Chrome.

### Installing

Download the file matching your operating system [here](https://git.3stadt.com/3stadt/presla/releases). In detail, you can choose as following:

* presla-x.x.x-linux-amd64 for Linux
* presla-x.x.x-darwin-amd64 for Mac
* presla-x.x.x-windows-amd64.exe for Windows

Save it in any directory. Give it executables permissions, `chmod +x <file>` on Linux/Mac.

Execute the file via console. Although double clicking is fine, this will close all log output when finished.

Access Presla via [http://localhost:8080](http://localhost:8080).

At this point you'll only see the included info presentation, since you yet need to create your presentations. 

### Configuration

Presla is configured via a [toml](https://github.com/toml-lang/toml/blob/master/README.md) file. At first start, a `presla.toml` file is created in the current directory, you may want to copy it over to another directory. The possible directories are:

- `<current dir>/presla.toml`
- `<homde dir>/.presla.toml`
- `<home dir>/.config/presla.toml`

**Please note:** Presla doesn't react to config changes at runtime. To apply your changes in presla.toml, you need to restart the application.

The default configuration file looks like this:

```toml
## The path to your markdown files.
## One markdown file holds one presentation
MarkdownPath="./"

## Whatever you want to show as text when including /svg/footer-text.svg
## By default shown in the lower right corner
FooterText="please edit presla.toml"

## The port to bind on. You should use localhost as host
ListenOn="localhost:8080"

## Optional: Path to your own template.
## Needs the index.html holding remarkjs, an info.md as starting point and footer-text.svg 
# TemplatePath="/home/user/Documents/presla-theme/templates"

## Optional: path to the templates static files
## Holds css, js, fonts and images used in your template
# StaticFiles="/home/user/Documents/presla-theme/static"

## Optional, define your own Executors for running code from the presentation
# CustomExecutors="/home/user/Documents/presla-executors"

## Optional, can be used multiple times
## This way you can specify a template used for only one presentation
# [[Presentations]]
# PresentationName="my_presentation"
# TemplatePath="/home/user/Documents/presla-theme-my-presentation/static"
# StaticFiles="/home/user/Documents/presla-theme-my-presentation/templates"
```

### Defining the path to your markup-files

In Presla, you create the content of your presentation using markdown files with a `.md` extension. Create them at a custom location and add that location to your config file, e.g.:


```
MarkdownPath="/home/user/some/path/talks/"
```

If you already have markdown files in the directory, they should be listed as presentations on the last slide of the default presentation.

## Working with Presla

Presla is ready to use, all styles an javascript is included in the executable. However, you might want to add your own styles.

### Including your theme

To include your theme, please edit the following definitions in presla.toml:

```
TemplatePath="/home/user/presla/Themes/your-presla-theme/templates"
StaticFiles="/home/user/presla/Themes/your-presla-theme/static"
```

Your own theme sould have the following layout:

```
|-- static
|   |-- css
|   |   `-- theme.css
|   |-- favicon.ico
|   |-- img
|   |   `-- logo.svg
|   `-- js
|       |-- ace
|       |   |-- ... # the ace editor files
|       |-- editor-loader.js
|       |-- remark-latest.min.js
|       `-- remark-loader.js
`-- templates
    |-- footer-text.svg
    |-- index.html
    `-- info.md

```

You can start by copying the directories mentioned above from the current [master branch](https://git.3stadt.com/3stadt/presla).

### Using individual themes for different presentations

If you need individual themes for different presentations you are able to configure that in `presla.toml`:

```
[[Presentations]]
PresentationName="my_presentation" # The name of your markdown file without extension
TemplatePath="/home/user/Documents/presla-theme-my-presentation/static"
StaticFiles="/home/user/Documents/presla-theme-my-presentation/templates"
```

All of the above is optional and can be defined multiple times.

## Changing the footer text

In order to change the footer text, you just need to set the following line of your presla.toml:

```
FooterText="some text"
```

## Authors

* @3stadt
* @leichteckig

## License

Please see LICENSE.md for details.

[![license](https://i.creativecommons.org/l/by-sa/4.0/80x15.png)](http://creativecommons.org/licenses/by-sa/4.0/)