# Javascript Object Decoding

A small fork of Go's encoding/json (Go 1.25) with minimal tweaks to accept JavaScript-style object keys during decoding:

- Unquoted identifier keys are allowed: `{a:1, b:2}`
- Single-quoted keys are allowed: `{'a':1, 'b':2}`

Everything else behaves like encoding/json, with these differences:

- String values may be single-quoted or double-quoted. Numbers/booleans/null unchanged.
 - The JavaScript literal `undefined` is accepted on decode and treated like null.
- Struct field tags use `js:"..."` instead of `json:"..."`.
- Marshal/Encoder output remains standard JSON.

## Install

```
go get github.com/nukilabs/js
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/nukilabs/js"
)

func main() {
    var m map[string]any
    _ = js.Unmarshal([]byte("{a:1, 'b':2, \"c\":3}"), &m)
    fmt.Println(m) // map[a:1 b:2 c:3]
}
```

## Notes

- Token streaming API also supports single-quoted and bare identifier object keys.

- This project copies Go's encoding/json under its BSD license (see LICENSE) and applies minimal changes in the scanner and decoder.
