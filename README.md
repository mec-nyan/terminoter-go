# Terminoter

[![go version](https://img.shields.io/badge/go-1.26.2-blue?style=for-the-badge)](https://go.dev/)
[![bubbletea](https://img.shields.io/badge/tea-2.0.6-orchid?style=for-the-badge)](https://github.com/charmbracelet/bubbletea)
[![lipgloss](https://img.shields.io/badge/lipgloss-2.0.3-mediumpurple?style=for-the-badge)](https://github.com/charmbracelet/lipgloss)


A simple and cute note-taking app for the terminal.

>[!NOTE]
> **AI Disclosure:**
> None of the contents of this repo was **AI** generated (not that there's anything wrong with that).


Demo:
![simple_demo](./demos/demo.gif)

## About

***Terminoter*** is a note-taking app to help you keep all your notes at hand when working on the
terminal/console.  It tries to be *nice to look at* without overcomplicating the **UI**.  Overall,
I try to keep things *as simple as possible, but not simpler*.

**Terminoter** is based on [**The Elm architecture** (TEA)](https://guide.elm-lang.org/)
and powered by the awesome
[**Charm**](https://charm.land/)'s libraries
[**Bubbletea**](https://github.com/charmbracelet/bubbletea) and
[**Lipgloss**](https://github.com/charmbracelet/lipgloss).

## Installation

```sh
go install github.com/mec-nyan/terminoter-go@latest
```

## Roadmap

- [x] Add/remove notes.
- [x] Support all kind of text (emoji, lists, etc).
- [x] Save to default location (usually `~/.local/share/terminoter/...`).
  - [x] Optionally read from/write to custom file.
- [ ] Different layouts according to terminal size.
  - [ ] Handle the case where the notes occupy more than the available space.
- [ ] Change layout.
- [ ] Collapse notes, showing only first line or title.
- [ ] Delete everything.
- [ ] Atomic writes.
- [ ] Edit notes.
  - [ ] Minimal edition facilities (i.e. emacs/vim binidngs).
    - [ ] Delete word (forwards/backwards).
    - [ ] Delete line (id).
    - [ ] Delete list, paragraph, etc.
    - [ ] Switch case.
    - [ ] ...
  - [ ] Minimal `terminal` markup (i.e. add `bold, italic, underline`, etc).
- [ ] Change notes' order.
- [ ] Save metadata (i.e. date/time).
- [ ] Backup file.
- [ ] Emoji chooser.
  - [ ] Icon/symbol chooser (Unicode/Nerdfont/other icons).
- [ ] Fix the demo/show emojis.
- [ ] Rewrite it in Rust :crab:.
