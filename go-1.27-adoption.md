# Go 1.27 Adoption

Assessment of what Go 1.27 enables for SOM, with a focus on generic methods.
Verified against a local `go1.27.0` toolchain (spec, `api/go1.27.txt`, and a compiled probe).

## 1. Generic methods

### What actually landed

The spec (Method declarations) now states:

> If the method declaration specifies type parameters (possibly in addition to type parameters
> declared by the receiver specification), the method name denotes a *generic method* [Go 1.27].
> Like a generic function, a generic method must be instantiated before it can be called or used
> as a value.

Verified by compiling with `go1.27.0`:

- Generic method on a non-generic type: `func (r *Res) Scan[T any]() ([]T, error)` — works.
- Generic method on a generic type: `func (b builder[M]) Select[P any](f func(*M) P) []P` — works.
- Inference from arguments works; explicit instantiation (`r.Scan[string]()`) works.

### The hard limit

> Interface methods cannot declare type parameters, but they may use type parameters from the
> interface declaration.

SOM generates `repo.Client`, `WeatherRepo`, `LocationRepo`, … as **interfaces**
(`repo/som.interfaces.go`, `repo/node.*.go`). No generic method can be added to any of them.

Generic methods are therefore only available on concrete types:

- `query.Builder[M]` / `builder[M]` / `SearchBuilder[M]` / `PageBuilder[M]` (structs)
- `som.RawResult`
- `query.Page[M]`, `query.Entry[M]`
- `repo.dbConn` (internal)

This is workable — the repo interfaces return concrete builder structs — but it rules out
anything like `Client.Query[M]()`.

## 2. Recommended changes, in priority order

### 2.1 `RawResult.Scan[T]` — highest user-facing value

`core/embed/som/som.base.go.tmpl:422`:

```go
func (r *RawResult) Scan(dest any) error
func (r *RawResult) ScanOne(dest any) error
```

Untyped out-parameter, no compile-time checking, awkward "pointer to a slice" contract.

With generic methods:

```go
func (r *RawResult) Scan[T any]() ([]T, error)
func (r *RawResult) ScanOne[T any]() (T, error)  // ErrNotFound → zero T
```

Keep the existing methods for backwards compatibility. No design risk, pure ergonomics.

### 2.2 Typed projections and ad-hoc aggregates — a genuinely new capability

This is the only item where 1.27 unlocks something SOM cannot express today. A projection type
`P` is not derivable from the model type `M`, so today it must break the fluent chain into a
package-level function. With generic methods:

```go
func (b builder[M]) Select[P any](fields ...lib.KeyProvider[M]) Projection[M, P]
func (g GroupBuilder[M]) Into[R any](ctx context.Context) ([]R, error)
```

That would give ad-hoc `SELECT`-subset and `GROUP BY` results without a declared view.
Currently `define/aggregate` only serves `define.View`.

### 2.3 `*dbConn` internals — cosmetic, cheap

`core/embed/som/repo/som.client.go.tmpl:57+` contains the classic
"methods cannot be generic" workaround:

```go
func dbSelect[T any](ctx context.Context, db *dbConn, id *models.RecordID) (*T, error)
func dbCreate[T any](ctx context.Context, db *dbConn, id models.RecordID, data any) (*T, error)
func dbCreateNew[T any](...)
func dbUpdate[T any](...)
func dbInsert[T any](...)
func dbQueryOne[T any](...)
```

These become `db.Select[T](ctx, id)` etc. Internal only, no API impact.

### 2.4 Collapsing the `*Async` wrappers — not recommended

`core/embed/som/query/builder.go.tmpl` carries ~14 hand-written `XAsync` twins
(`CountAsync`, `AllAsync`, `FirstAsync`, …). A single generic method could replace them:

```go
func (b builder[M]) Async[R any](ctx context.Context, fn func(context.Context) (R, error)) *asyncResult[R]
```

But `b.Async(ctx, b.Count)` reads worse than `b.CountAsync(ctx)`, and the free function
`async(ctx, fn)` already exists today. API quality beats the line count here — skip.

## 3. What Go 1.27 does *not* fix

`modelInfo[M]` (`query/builder.go.tmpl:25`) and `RepoInfo[M]` (`repo/som.node.go.tmpl`) are
structs of per-model closures. They exist because the conversion type `C` must not leak into the
public `Builder[M]` / repo API. Generic methods do not help: `C` is fixed per model, not chosen
per call. The only alternative is `builder[M, C]`, which is a breaking API change and was already
possible before 1.27. Leave as is.

## 4. Other Go 1.27 items relevant to SOM

Full list of non-crypto/non-`net/http` API additions in `api/go1.27.txt`:
`bytes`, `hash/maphash`, `math/big`, `math/rand/v2`, `net/url`, `strings`, `testing/synctest`,
`unicode`, `uuid`.

### 4.1 Standard library `uuid` package — high value

Go 1.27 adds a stdlib `uuid`: `UUID [16]byte`, `New`, `NewV4`, `NewV7`, `Parse`, `MustParse`,
`Nil`, `Max`, `Compare`, `String`, and text marshalling.

SOM currently carries duplicated templates for each supported UUID library:

- `internal/types/uuid_google.go.tmpl`, `internal/types/uuid_gofrs.go.tmpl`
- `internal/distinct/distinct_uuid_google.go.tmpl`, `…_gofrs.go.tmpl`
- `internal/cbor/helpers_uuid_google.go`, `helpers_uuid_gofrs.go`
- `internal/lib/filter.uuid_google.go.tmpl`, `filter.uuid_gofrs.go.tmpl`

Adding stdlib `uuid` as a third supported library follows the same mechanical pattern and lets
users drop both third-party dependencies. Detection lives in `core/parser/field.go` and
`core/codegen/def/def.go`.

### 4.2 `testing/synctest.Sleep` — limited value

`internal/tests` contains 6 `time.Sleep` calls, but those are integration tests against a real
SurrealDB container; a synctest bubble does not cover real network I/O. Only useful for
pure-Go timing logic (expiry purge scheduling in `runExpiryPurge`, cache TTL in
`internal/cache.go`) if that logic is ever tested in isolation.

### 4.3 Not applicable

- `strings.CutLast` / `bytes.CutLast` — no `LastIndex` usage anywhere in the codebase.
- `encoding/json` v2 options — SOM serialises via CBOR, not JSON.
- `net/url.Clone`, `hash/maphash.Hasher`, `math/rand/v2`, `math/big` — no current use case.

## 5. Rollout constraint

The generated code inherits the user's Go version, not SOM's. Emitting generic methods in the
embed templates raises the floor for every consumer of SOM.

During this phase of the project only the latest Go version is supported (see `SECURITY.md`), so
requiring 1.27 in generated code is acceptable policy-wise. If a softer path is ever wanted,
`core/util/gomod` already reads the target module's `go` directive and could gate emission.

## 6. Suggested order

1. `RawResult.Scan[T]` / `ScanOne[T]` (2.1)
2. `*dbConn` methods (2.3)
3. Standard library `uuid` support (4.1)
4. Typed projections / group-by (2.2)
