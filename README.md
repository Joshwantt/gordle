# gordle
Lightweight Implementation of Wordle with Go and Astro

## Dependencies

- [Go](https://go.dev/dl/) 1.27.1
- [Node.js](https://nodejs.org/) 24.21.0
- [pnpm](https://pnpm.io/installation) 12.5.1
- [air](https://github.com/air-verse/air) for live rebuilds in development
- make

The Go, Node.js and pnpm versions are pinned in `mise.toml`. If you use [mise](https://mise.jdx.dev/), install all three with:

```sh
mise install
```

Install air:

```sh
go install github.com/air-verse/air@latest
```

Install the frontend packages:

```sh
pnpm -C web install
```

Go modules download automatically on the first build.

## Development

```sh
make dev
```

This builds the site and the server, then serves both at http://localhost:8080. Any change to a Go, Astro, TypeScript or CSS file triggers a rebuild and restart; refresh the browser to see it.

Run the tests:

```sh
make test
```

## Build

```sh
make build
./bin/gordle
```

The server serves the built site from `web/dist`. When running it from another directory, pass the path with `-static`, and change the port with `-address`:

```sh
./bin/gordle -static /path/to/web/dist -address :9000
```
