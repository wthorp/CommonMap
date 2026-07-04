# CommonMap
CommonMap turns RPF map data into a deployable web map service.

The goal is simple: take existing NGA-style RPF holdings and make them usable in real GIS tools without copying, tiling, or converting the source data first. CommonMap indexes the input, preserves the original data and metadata, and serves the result as a WMS that can be used in third-party applications.

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

CommonMap is built around very fast indexing of large RPF holdings. This is a direct-indexing workflow, not a copy-and-convert pipeline:

- indexes the source files in place
- avoids building a separate raster cache before use
- scales to very large collections quickly enough to make ad hoc deployment practical

CommonMap was built to make huge RPF collections operational fast. Point it at the source data, build the index quickly, and start serving map content without reprocessing the archive into another format first.

## Typical Use

CommonMap is intended to sit between a directory of RPF holdings and the client application that needs map context.

1. Point CommonMap at the source data.
2. Let it index the files and build its map artifacts.
3. Use the exposed WMS from applications such as QGIS, ArcGIS, or other GIS clients.

The result is a shareable background map service built from the source data you already have.

## Related Government Use

CommonMap was built for deployable geospatial workflows where users need reliable background map context inside existing mission and analysis tools.

For related context, Vantor's public GEGD Pro material is here:

- [Vantor: NGA award to operate and enhance GEGD Pro](https://vantor.com/blog/nga-awards-vantor-dollar70m-option-year-contract-to-operate-and-enhance-governments-primary-commercial-geoint-platform/)

GEGD Pro is Vantor's secure platform for government GEOINT delivery. It is not a public portal; access is limited to authorized users with credentials. Vantor's public contact/support page for the program is here:

- [Vantor GEGD Program Support](https://vantor.com/contact-us/)

GPL3 licensed.
Other commercial licenses available by request.
