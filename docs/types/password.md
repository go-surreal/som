# Password Type

The password type provides secure password storage with automatic hashing and protection against accidental exposure.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `som.Password[Algorithm]` / `*som.Password[Algorithm]` |
| Database Schema | `string` / `option<string>` |
| CBOR Encoding | Direct (as string) |
| Sortable | No |

## Supported Algorithms

SOM supports four password hashing algorithms:

| Algorithm | Go Type | SurrealDB Function |
|-----------|---------|-------------------|
| Bcrypt | `som.Password[som.Bcrypt]` | `crypto::bcrypt::generate` |
| Argon2 | `som.Password[som.Argon2]` | `crypto::argon2::generate` |
| Pbkdf2 | `som.Password[som.Pbkdf2]` | `crypto::pbkdf2::generate` |
| Scrypt | `som.Password[som.Scrypt]` | `crypto::scrypt::generate` |

## Definition

```go
type User struct {
    som.Node

    Username string
    Password som.Password[som.Bcrypt]  // Required, using Bcrypt
}

type Admin struct {
    som.Node

    Email    string
    Password som.Password[som.Argon2]   // Required, using Argon2
    Recovery *som.Password[som.Bcrypt]  // Optional backup password
}
```

## How It Works

### Automatic Hashing

Passwords are automatically hashed by SurrealDB when stored. The hashing only occurs when:

1. A new record is created with a password
2. The password value changes (different from `$before`)

This means updating other fields won't re-hash an unchanged password.

### Security: SELECT Permissions

Password fields have `PERMISSIONS FOR SELECT NONE`, which means:

- Passwords are **never returned** in query results
- Even with direct database access, hashes are not exposed
- The field will be empty when reading records

```surql
DEFINE FIELD password ON TABLE user TYPE string
    VALUE IF $value != NONE AND $value != NULL AND $value != "" AND $value != $before
          THEN crypto::bcrypt::generate($value)
          ELSE $value END
    PERMISSIONS FOR SELECT NONE;
```

## Creating Passwords

```go
user := &model.User{
    Username: "alice",
    Password: som.Password[som.Bcrypt]{Value: "plaintext_password"},
}

err := client.UserRepo().Create(ctx, user)
// Password is automatically hashed in the database
```

## Updating Passwords

```go
// Read user (password field will be empty)
user, _, _ := client.UserRepo().Read(ctx, userID)

// Update password
user.Password = som.Password[som.Bcrypt]{Value: "new_password"}
err := client.UserRepo().Update(ctx, user)
// Only the new password is hashed; other fields unchanged
```

## Password Verification

SOM provides a type-safe `Verify()` filter method for password verification:

```go
// Type-safe password verification
user, ok, err := client.UserRepo().Query().
    Filter(
        where.User.Username.Equal("alice"),
        where.User.Password.Verify("plaintext_password"),
    ).
    First(ctx)

if err != nil {
    // Handle error
}
if !ok {
    // Invalid credentials
}
// User authenticated successfully
```

The `Verify()` method automatically uses the correct crypto comparison function based on the password's algorithm type (Bcrypt, Argon2, etc.).

### Raw Query Alternative

You can also use raw SurrealQL queries if needed:

```go
// Verify login credentials using a raw query
query := `SELECT * FROM user WHERE username = $username AND crypto::bcrypt::compare(password, $password)`
```

For Argon2:
```go
query := `SELECT * FROM user WHERE username = $username AND crypto::argon2::compare(password, $password)`
```

## Filter Operations

Password fields have limited filter operations due to their secure nature:

### Password Verification

```go
// Verify password against stored hash
where.User.Password.Verify("plaintext_password")
```

This uses the appropriate crypto comparison function for the password's algorithm.

### Nil Operations (Pointer Types Only)

```go
// Check if optional password is set
where.Admin.Recovery.IsNil()

// Check if optional password exists
where.Admin.Recovery.IsNotNil()
```

### Zero Value Check

```go
// Check if password is empty (not set)
where.User.Password.Zero(true)

// Check if password has a value
where.User.Password.Zero(false)
```

## Algorithm Comparison

| Algorithm | Speed | Memory | Security | Use Case |
|-----------|-------|--------|----------|----------|
| **Bcrypt** | Slow | Low | High | General purpose, widely supported |
| **Argon2** | Configurable | High | Very High | Modern applications, recommended |
| **Pbkdf2** | Fast | Low | Good | Legacy compatibility |
| **Scrypt** | Slow | High | High | Memory-hard requirements |

### Recommendations

- **New applications**: Use `som.Argon2` for best security
- **Compatibility needs**: Use `som.Bcrypt` for wide support
- **Memory-constrained**: Use `som.Pbkdf2`

## Complete Example

```go
package main

import (
    "context"
    "yourproject/gen/som"
    "yourproject/gen/som/where"
    "yourproject/model"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create user with password
    user := &model.User{
        Username: "alice",
        Email:    "alice@example.com",
        Password: som.Password[som.Bcrypt]{Value: "secure_password_123"},
    }
    client.UserRepo().Create(ctx, user)

    // Reading user - password field is empty (protected)
    found, _, _ := client.UserRepo().Read(ctx, user.ID())
    // found.Password.Value is "" (empty, never returned)

    // Change password
    found.Password = som.Password[som.Bcrypt]{Value: "new_secure_password"}
    client.UserRepo().Update(ctx, found)

    // Find users with recovery password set
    adminsWithRecovery, _ := client.AdminRepo().Query().
        Filter(where.Admin.Recovery.IsNotNil()).
        All(ctx)
}
```

## Security Best Practices

1. **Never log passwords**: Even though they're hashed, avoid logging password fields
2. **Use strong passwords**: The hashing algorithm doesn't protect weak passwords
3. **Choose appropriate algorithm**: Argon2 for new apps, Bcrypt for compatibility
4. **Rate limit authentication**: Protect against brute-force attacks at the application level
5. **Use HTTPS**: Always transmit passwords over encrypted connections

## Limitations

- Passwords cannot be sorted (not meaningful for hashed values)
- Passwords cannot be compared with equality filters (hashes differ each time)
- Reading a password field always returns empty value
