# CommonMap
CommonMap turns RPF map data into a deployable web map service.

The goal is simple: take existing NGA-style RPF holdings and make them usable in normal GIS tools without copying, tiling, or converting the source data first. CommonMap indexes the input, preserves the original data and metadata, and serves the result as a WMS that can be used in third-party applications.

## What It Does

- Indexes RPF data directly instead of building a separate tile cache
- Preserves the source imagery's intrinsic metadata
- Serves the result as a WMS for downstream GIS clients
- Chooses the right data series for the current zoom level
- Falls back to coarser imagery, then vector context, when finer coverage is missing

## Why RPF

CommonMap is built around RPF because it is compact, metadata-rich, and already widely available in the environments this tool was built for. That makes it a good fit for:

- deployable/offline use
- fast copying and low storage overhead
- map context inside mission or analysis tools that already speak WMS

Compared to prebuilt tile caches, the intended advantage is that the original data stays intact, the metadata is preserved, and the deployment footprint stays smaller.

## Fast Indexing

CommonMap is designed around very fast indexing of large RPF holdings. The original demo positions it as a direct-indexing workflow rather than a copy-and-convert pipeline:

- indexes the source files in place
- avoids building a separate raster cache before use
- scales to very large collections quickly enough to make ad hoc deployment practical

In the demo, CommonMap is described as indexing on the order of tens of thousands of files per second and handling tens of millions of files in minutes on representative hardware. The important part for this project is the workflow: point it at the source data, build the index fast, and start serving map content without reprocessing the entire archive into another format.

## Typical Use

CommonMap is intended to sit between a directory of RPF holdings and the client application that needs map context.

1. Point CommonMap at the source data.
2. Let it index the files and build its map artifacts.
3. Use the exposed WMS from applications such as QGIS, ArcGIS, or other GIS clients.

The result is a shareable background map service built from the source data you already have.

## Related Government Use

The original CommonMap pitch was aimed at deployable geospatial workflows where users need background map context inside existing mission or analysis tools.

That is adjacent to platforms such as Vantor's GEGD Pro portal, but not the same thing. Public Vantor materials describe GEGD Pro as a secure, web-based GEOINT platform built for NGA and used across U.S. government organizations rather than a public portal. In practice, that means access is limited to authorized users with credentials rather than open public access.

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
