# Validation

A model, or any type used within it, can validate itself by providing a `Validate() error` method.
som detects those methods while generating the code and emits the calls directly into the write
operations of the repository. Types without such a method get no call at all — there is no runtime
interface check and no reflection involved.

## Defining a Rule

```go
package model

import (
    "errors"
    "strings"

    "yourproject/gen/som"
)

type User struct {
    som.Node[som.ULID]

    Name  string
    Email Email
}

func (u *User) Validate() error {
    if u.Name == "" {
        return errors.New("name must not be empty")
    }

    return nil
}

// Email is a named scalar that validates itself, wherever it is used.
type Email string

func (e Email) Validate() error {
    if !strings.Contains(string(e), "@") {
        return errors.New("not an email address")
    }

    return nil
}
```

Both a value and a pointer receiver work, and a method promoted from an embedded type counts as
well. The rule applies to every type som encounters within a model: the model itself, nested
structs, named scalars and the elements of a slice — including types from other packages.

## Where It Runs

Validation runs on every write of a full model:

| Operation | Validated |
|---|---|
| `Create`, `CreateWithID`, `Insert`, `Update` | yes |
| Sink `Create`, `Insert` | yes |
| `Relate().<Edge>().Create` | yes |
| `Read`, `Query`, `Refresh`, `Resolve` | no |
| `Delete`, `Erase`, `Restore` | no |

**Validation never runs on read.** A record stored before a rule existed, or written by another
application, would otherwise be impossible to load — and therefore impossible to fix.

Within a write, the checks run **after** the before-hooks (`OnBeforeCreate`, `OnBeforeUpdate` and
the hook interfaces of the model) and before the statement is sent to the database. A hook is thus
free to fill in or correct values, and what the hooks leave behind is what gets validated. Nothing
is written if a check fails.

Values are checked from the leaves up: the deepest value is checked before the struct holding it,
and the model itself is checked last. The first failing check aborts the write.

## Handling the Error

A failed check is wrapped into a `som.ValidationError`, which names the offending value:

```go
err := client.UserRepo().Create(ctx, user)

var validationErr *som.ValidationError
if errors.As(err, &validationErr) {
    fmt.Println(validationErr.Path) // e.g. "Contacts[1].Email"
    fmt.Println(validationErr.Err)  // the error the Validate method returned
}
```

`Path` is the path of the value within the model, with slice indexes in brackets. It is empty when
the model type itself rejected the write. The error returned by `Validate` stays reachable through
`errors.Is` and `errors.As`, since `ValidationError` unwraps to it.

A caller that only needs to tell a rejected write from a failed one can match a sentinel instead:
`som.ErrValidation` for a failed `Validate` method, or the broader `som.ErrInvalid`, which covers
every write rejected because of a value it carried — including a
[field constraint](12_constraints.md) the database enforced:

```go
if errors.Is(err, som.ErrInvalid) {
    http.Error(w, err.Error(), http.StatusUnprocessableEntity)
}
```

## Checking Ahead of a Write

The walkers live in the generated `validate` package, with one exported entry point per model.
Calling it runs the exact same checks the repository runs, which is useful to reject input before
anything else is done with it:

```go
if err := validate.User(user); err != nil {
    return err
}
```

## Record Links Are Not Followed

A field pointing to another node or edge is never validated as part of the model holding it. Such a
record is written through its own repository, and a link that was not fetched holds nothing but its
id — validating it would fail on data that is simply not loaded.

## Together With Field Constraints

A [field constraint](12_constraints.md) declared via a som tag is enforced by the database, this
feature is enforced by the generated Go code before the statement is sent. When a rule is expressed
in both places, the `Validate` method rejects the write first and the database constraint never
fires.

Which to reach for:

- **Tag constraints** for invariants the database has to guarantee no matter who writes — other
  applications, a migration, a manual query.
- **`Validate` methods** for business rules that need Go logic, several fields at once, or context
  the schema cannot express.

## Opting Out

Some types validate themselves in ways that are expensive or simply not wanted, which is common for
third party types. A field can opt out with the `novalidate` tag, which excludes the field and
everything below it:

```go
type Area struct {
    som.Node[som.ULID]

    // The geometry is checked by the writer, not on every write.
    Shape geom.Polygon `som:"novalidate"`
}
```
