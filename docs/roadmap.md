# Roadmap

The features below are a loose plan of things that need to happen before presla reaches its full potential.

Expect changes as this is a spare time project.

## v0.0.11

- Template download from [presla-templates](https://github.com/3stadt/presla-templates) via console command
- Configuration via webbrowser

## v0.0.12

- More executors
- One alternative template
- More beautiful default template
- Better documentation of template creation
- Default and alternative templates downloadable via zip from default presentation // generally better docs on custom templates
- Better [decktape](https://github.com/astefanutti/decktape) integration or documentation for PDF generation
- Some sort of [slideshare](https://slideshare.net)/[speakerdeck](https://speakerdeck.com/) integration of documentation

---

## Released versions

## v0.0.4

- Unit tests
- Better support for directory paths in windows

## v0.0.5

- More functionality exposed to executors:
  - The current OS
  - Check if software X is installed
  - Check if a docker image is installed, install or error if not
    - See [docker-registry-client](https://github.com/heroku/docker-registry-client) and [go-dockerclient](https://github.com/fsouza/go-dockerclient)
  - Browser log output directly from executors
  - execSilent function that sends no log output to browser

## v0.0.6

- Slimmed down ace editor
- Auto update mechanism
- Multiple executors per editor  

## v0.0.7

Minor fixes

- Footer text left aligned
- Mode python included in default presentation

## v0.0.8

- When cloning a presentation via `c` in the browser, the editors are now synced
  - Useful when presentation view + beamer view
- Editor/log view height & width can now be set via data attributes
  - Usage: `data-editorwidth="60rem" data-editorheight="20rem" data-logwidth="59rem" data-logheight="5rem"`

## v0.0.9

- optional command line underneath the presentation for sending arguments to the code

## v0.0.10

- Add [mermaidjs](https://mermaidjs.github.io/) diagrams
- Add possibility to load code files from markdown folder
  - Usage: Add an attribute to your editor definition, example: `data-contenturl="/some/path/file.php"`
