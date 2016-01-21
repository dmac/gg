# gg

Package gg provides a common interface to OpenGL backends.
It enables compiling shared graphics code for different targets, including desktop and web.

## Stability

This project is not yet stable and may contain breaking API changes.

## Currently supported backends

- OpenGL 2.1
- WebGL

## TODO

- Usage docs
- Create a helper package (ggh) and consider moving higher-level conveniences like Sprites into it.
- Rewrite tetris example using new API
- Can usage of interface{} be improved?
  - TexImage2D data argument
  - gg types (Buffer, Shader, etc.)
- Make example web page with playable game
