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
err = client.UserRepo().Relate().Follows().Create(txCtx, follows)

return som.TxCommit(txCtx)
```

## Locked Reads

A plain `Read` inside a transaction sees a consistent snapshot, but the transaction still commits if the record was changed by someone else in the meantime. When a record is only read to make a decision, that silently builds on stale data. `ReadForUpdate` closes this gap using SurrealDB's `SELECT ... FOR UPDATE` (requires SurrealDB 3.3+):

```go
func placeOrder(ctx context.Context, client som.Client, accountID string) error {
    txCtx, cancel := som.TxStart(ctx)
    defer cancel()

    account, exists, err := client.AccountRepo().ReadForUpdate(txCtx, accountID)
    if err != nil {
        return err
    }
    if !exists || account.Status != "active" {
        return errors.New("account is not active")
    }

    if err := client.OrderRepo().Create(txCtx, &model.Order{Account: account}); err != nil {
        return err
    }

    return som.TxCommit(txCtx)
}
```

If another transaction modifies the account before `TxCommit`, the commit fails and `som.IsTransactionConflict(err)` reports `true`. The whole transaction can then be retried:

```go
for {
    err := placeOrder(ctx, client, accountID)
    if !som.IsTransactionConflict(err) {
        return err
    }
}
```

- **No lock is taken**: Nobody waits. The conflict check runs at commit time, so plan for retries.
- **Missing records count too**: Reading a non-existent ID registers it as well. If a concurrent transaction creates the record, the commit fails. This makes get-or-create logic safe.
- **Only needed for records you don't write**: A record the transaction also writes is already checked for conflicts on commit.
- **Single records only**: The read is bound to a record ID. Queries over a table are not protected, so concurrently created records matching a filter are not detected.
- **Transaction required**: Outside of a transaction, `ReadForUpdate` returns `ErrTransactionRequired`. It never uses the cache.

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
| `ErrTransactionRequired` | `ReadForUpdate` called without an active transaction |
