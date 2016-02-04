# gg

Package gg provides a common interface to OpenGL backends.
It enables compiling shared graphics code for different targets, including desktop and web.

## Stability

This project is not yet stable and may contain breaking API changes.
The full OpenGL API is not yet implemented by the `gg.Backend` interface.

## Currently supported backends

- OpenGL 2.1
- WebGL

## Examples

The examples target two platforms: native (OpenGL 2.1) and web (WebGL). To build for each platform:

```
# native
$ go build

# web
$ gopherjs build
$ python -m SimpleHTTPServer
```
