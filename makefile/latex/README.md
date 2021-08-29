# Latex

Makefile for compiling Latex files with two different approaches.

## Docker

1. Build a Ubuntu based Docker image with Latex tools (`make docker-image`)
2. Run the Docker image, mount the current working directory inside the image, and execute `build.sh` script inside the container (`make docker-compile`)

## Tectonic

> A modernized, complete, self-contained TeX/LaTeX engine, powered by XeTeX and TeXLive.
> [Tectonic](https://github.com/tectonic-typesetting/tectonic)

Only the Tectonic executable is required. The rest is managed by Tectonic. Additionally, there is a (un)watch command to recompile Tex files on each change (requires [entr](https://eradman.com/entrproject/)).
