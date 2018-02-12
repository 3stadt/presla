# Presla

Here you can find an useful way to create your presentation using presla. Written in go, it can be used easily on every operating system.


## Getting Started

### Prerequisites

In general, there don't have to be any programs pre-installed. 

### Installing

At first, you just need to download the file matching your operating system [here](https://git.3stadt.com/3stadt/presla/releases). In detail, you can choose as following:
* presla-0.0.1-linux-amd64 for Linux
* presla-0.0.1-darwin-amd64 for Mac
* presla-0.0.1-windows-amd64.exe for Windows

Save it in a directory as you wish. Keep in mind that you may need to adjust the permissions of this file to execute it properly. Now, you should be able to access presla via [http://localhost:8080
](http://localhost:8080), containing no presentations yet. There you go! :) ... Well, we need to get some things done first.

### Doing first configuration steps

As you execute this file the first time, you may notice that a file called "presla.toml" was created in your folder. Here you can configure the path to your markup files and which themes you're using.

You can leave it in this directory, it will already work. However, you can save it in the .config folder of your user or in your users' folder itself if you're running Linux. Presla will search for this presla.toml in the following order:

```
No config file at presla.toml ...
No config file at /home/user/.presla.toml ...
Using config file: /home/user/.config/presla.toml
```

In the case of this example, the config.toml is located in the .config folder. 

**Keep in mind** that presla doesn't reacts to chances at runtime. To apply your changes in presla.toml, you need to restart it.

### Defining the path to your markup-files

In Presla, you create the content of your presentation using markup-files (.md). You can save them in a location as you wish, but you need to make the corresponding path public to presla by yourself. That can be done in the presla.toml using the following line:


```
## The path to your markdown files.
## One markdown file holds one presentation
MarkdownPath="/home/user/Documents/Talks/"
```

Please define your path as a MarkdownPath here. If you already have some markdown-files in there, they should be listed as presentations now.

As you know, the server for using presla is running on http://localhost:8080 by default. If you wish to adjust that, you can edit the following line as you please:

 ```
## The port to bind on. You should use localhost as host
ListenOn="localhost:8080"
 ```

## Working with Presla

If you wish, you can use Presla right away. However, you might need further styles as provided in the default. Fortunately, you're able to use your own themes for your presentations.

### Including your theme

To include your theme, please edit the following lines in presla.toml:

```
## Optional: Path to your own template.
## Needs the index.html holding remarkjs, an info.md as starting point and footer-text.svg 
TemplatePath="/home/user/presla/Themes/your-presla-theme/templates"

## Optional: path to the templates static files
## Holds css, js, fonts and images used in your template
StaticFiles="/home/user/presla/Themes/your-presla-theme/static"
```

Just change these paths to the ones holding your theme's files and you're good to go.

### Using individual themes for different presentations

You might wish to have individual themes depending your presentation. Therefore you just need to look at the following paragraph of your presla.toml:

```
## Optional, can be used multiple times
## This way you can specify a template used for only one presentation
# [[Presentations]]
# PresentationName="my_presentation"
# TemplatePath="/home/user/Documents/presla-theme-my-presentation/static"
# StaticFiles="/home/user/Documents/presla-theme-my-presentation/templates"
```

You need to define the name of your presentation and set the path in which your individual theme and its static files are located. This way you can assign an individual theme to the corresponding presentation.


## Changing the footer text

In order to change the footer text, you just need to set the following line of your presla.toml as you wish:

```
## Whatever you want to show as text when including /svg/footer-text.svg
## By default shown in the lower right corner
FooterText="@yourname"
```

## How Presla works in detail

To be continued~

## Authors

* @3stadt
* @leichteckig

## License

Please see LICENSE.md for details.

[![license](https://i.creativecommons.org/l/by-sa/4.0/80x15.png)](http://creativecommons.org/licenses/by-sa/4.0/)

## Acknowledgments

To be continued~





