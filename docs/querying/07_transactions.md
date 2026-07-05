# Transactions

SOM supports client-side transactions that group multiple operations into an atomic unit. Either all operations succeed, or none are applied.

## Starting a Transaction

```go
txCtx, cancel := som.TxStart(ctx)
defer cancel()
```

`TxStart` returns a new context carrying the transaction state and a cancel function. Always defer the cancel function to ensure cleanup on errors.

## Committing

```go
err := som.TxCommit(txCtx)
if err != nil {
    return err
}
```

## Cancelling / Rolling Back

Call the cancel function returned by `TxStart`, or use `TxCancel` explicitly:

```go
err := som.TxCancel(txCtx)
```

The deferred `cancel()` also rolls back if the transaction hasn't been committed.

## Complete Example

```go
txCtx, cancel := som.TxStart(ctx)
defer cancel()

user := &model.User{Name: "Alice", Email: "alice@example.com"}
if err := client.UserRepo().Create(txCtx, user); err != nil {
    return err // cancel() rolls back on return
}

profile := &model.Profile{UserID: user.ID(), Bio: "Hello!"}
if err := client.ProfileRepo().Create(txCtx, profile); err != nil {
    return err // cancel() rolls back on return
}

// Atomically commit both operations
return som.TxCommit(txCtx)
```

## Cross-Repository Transactions

Transactions work across multiple repositories. All operations using the same transaction context are grouped:

```go
txCtx, cancel := som.TxStart(ctx)
defer cancel()

err := client.UserRepo().Create(txCtx, user)
err = client.PostRepo().Create(txCtx, post)
err = client.FollowsRepo().Relate().From(user1).To(user2).Create(txCtx, follows)

return som.TxCommit(txCtx)
```

## Behavior Notes

- **No nesting**: Calling `TxStart` on a context that already has an active transaction will panic with `ErrTransactionAlreadyActive`.
- **Idempotent close**: Calling commit or cancel multiple times is safe; subsequent calls return `ErrTransactionClosed`.
- **Cache bypass**: Operations within a transaction bypass the context cache and always return fresh data from the database.
- **Live queries**: Live queries are not supported within transactions (`ErrLiveNotSupportedInTx`).

## Error Types

| Error | When |
|-------|------|
| `ErrTransactionClosed` | Operation attempted after commit or cancel |
| `ErrTransactionAlreadyActive` | Nested `TxStart` attempted (panics) |
| `ErrLiveNotSupportedInTx` | Live query attempted within a transaction |
