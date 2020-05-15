# Overview

[![GoDoc](https://godoc.org/github.com/ysmood/byframe?status.svg)](https://pkg.go.dev/github.com/ysmood/byframe?tab=doc)

It's a low overhead length header format with dynamic header length.
So we don't waste resources on the header itself when framing data.
The algorithm is based on LEB128.

This lib also contains functions to encode and decode the header,
so you have the full flexibility to decide how to use it,
such as streaming TCP frames, indexing database, etc.

## Format

Each frame has two parts: the header and body.

```txt
|     frame     |
| header | body |
```

### Header

Each byte (8 bits) in the header has two parts, "continue" and "fraction":

```txt
bit index |    0     | 1 2 3 4 5 6 7 |
sections  | continue |   fraction    |
```

If the "continue" is 0, the header ends.
If the "continue" is 1, then the followed byte should also be part of the header.

Sum all the fractions together, we will get the size of the message.

For example:

```txt
|                         frame                              |
|                      header                         | body |
| continue |   fraction    | continue |   fraction    |      |
|    0     | 1 0 0 0 0 0 0 |    1     | 1 1 0 1 0 0 0 | ...  |
```

So the size of the body is 0b1101000,1000000 bytes.
