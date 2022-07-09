# si(storage interface)
`si` is a collection of wrappers that aims to ease reading/writing data from/to repositories. It is mostly a client side library and the following repositories or communication protocols will be supported from standard or non-standard packages.

- file
- tcp
- sql
- http
- websocket ([gorilla websocket](https://github.com/gorilla/websocket))
- kafka [sarama](https://github.com/Shopify/sarama)
- elasticsearch [go-elasticsearch](https://github.com/elastic/go-elasticsearch)
- ftp

## Installation
```bash
go get -u github.com/go-wonk/si
```

## Quick Start


## Versions
### v0.1.1
- `siwrap` package has been renamed to `sisql`.