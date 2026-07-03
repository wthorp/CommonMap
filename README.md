# CommonMap
Fast RPF indexing and MapServer

Not very cleaned up but moved to go modules.

## Development

Bootstrap a fresh clone:

```sh
make setup
```

This installs the local tool binaries into `./bin` and installs the repo-managed Git hooks into `.git/hooks`.

If you only need one half of setup:

```sh
make tools
make hooks
```

Common local commands:

```sh
make fmt
make check
make vuln
```

Git hook behavior:

- `pre-commit`: runs `gofmt` on staged Go files and re-stages fixes, then runs `go test ./...`
- `pre-push`: runs `make check`

GPL3 licensed.
Other commercial licenses available by request.
