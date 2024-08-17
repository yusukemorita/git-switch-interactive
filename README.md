# README

A simple tool to interactively switch git branches.

## Installation

```
go install github.com/yusukemorita/git-switch-interactive@latest
```

or build locally

```
go build -o git-switch-interactive; mv ./git-switch-interactive ~/go/bin/
```

## Releasing a new version

```
git tag v0.2.0
git push origin v0.1.2
```

## Thanks

- [Writing an interactive CLI menu in Golang](https://medium.com/@nexidian/writing-an-interactive-cli-menu-in-golang-d6438b175fb6) was a huge help, thanks @Nexidian !
