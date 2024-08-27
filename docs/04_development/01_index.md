# Development

## Versioning

In the future this package will follow the [semantic versioning](https://semver.org) specification.

Up until version 1.0 though, breaking changes might be introduced at any time (minor version bumps).

## Compatibility

This go project makes heavy use of generics. As this feature has been introduced with go 1.18, that version is the
earliest to be supported by this library.

In general, the two latest (minor) versions of go - and within those, only the latest patch - will be supported
officially. This means that older versions might still work, but could also break at any time, with any new
release and without further notice.

Deprecating an "outdated" go version does not yield a new major version of this library. There will be no support for
older versions whatsoever. This rather hard handling is intended, because it is the official handling for the go
language itself. For further information, please refer to the
[official documentation](https://go.dev/doc/devel/release#policy) or [endoflife.date](https://endoflife.date/go).
