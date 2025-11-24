//go:build embed

package repo

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type ID = models.RecordID

type Database interface {
	// Create accepts either a table name (string) for server-generated IDs or a RecordID for specific IDs.
	// The omit parameter specifies fields to exclude from the returned result.
	Create(ctx context.Context, what any, data any, omit []string) ([]byte, error)
	// Select retrieves a record by ID. The omit parameter specifies fields to exclude from the result.
	Select(ctx context.Context, id *ID, omit []string) ([]byte, error)
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
	Live(ctx context.Context, statement string, vars map[string]any) (<-chan []byte, error)
	// Update modifies a record. The omit parameter specifies fields to exclude from the returned result.
	Update(ctx context.Context, id *ID, data any, omit []string) ([]byte, error)
	Delete(ctx context.Context, id *ID) ([]byte, error)

	Marshal(val any) ([]byte, error)
	Unmarshal(buf []byte, val any) error
	Close() error
}

// Config holds the configuration for connecting to the SurrealDB instance.
type Config struct {
	Address   string
	Namespace string
	Database  string
	Username  string
	Password  string
}

type ClientImpl struct {
	db Database
}

// surrealDBWrapper wraps the official surrealdb.go client to implement the Database interface.
type surrealDBWrapper struct {
	db *surrealdb.DB
}

func (w *surrealDBWrapper) Create(ctx context.Context, what any, data any, omit []string) ([]byte, error) {
	// For CREATE, we use RETURN AFTER and then SELECT with OMIT if needed,
	// since SurrealDB doesn't support OMIT with CREATE/UPDATE RETURN clauses.
	var statement string
	var selectID string

	switch v := what.(type) {
	case *newRecordID:
		statement = fmt.Sprintf("CREATE %s CONTENT $data RETURN AFTER", v.String())
		selectID = v.String()
	case string:
		statement = fmt.Sprintf("CREATE %s CONTENT $data RETURN AFTER", v)
	case models.RecordID:
		statement = fmt.Sprintf("CREATE %s CONTENT $data RETURN AFTER", v.String())
		selectID = v.String()
	case models.Table:
		statement = fmt.Sprintf("CREATE %s CONTENT $data RETURN AFTER", string(v))
	default:
		return nil, fmt.Errorf("invalid type for 'what' parameter: %T (expected string, RecordID, or Table)", what)
	}

	queryResult, err := surrealdb.Query[[]any](ctx, w.db, statement, map[string]any{"data": data})
	if err != nil {
		return nil, fmt.Errorf("failed to execute: %w", err)
	}
	if queryResult == nil || len(*queryResult) == 0 {
		return nil, fmt.Errorf("empty response from create")
	}
	if (*queryResult)[0].Error != nil {
		return nil, fmt.Errorf("create failed: %w", (*queryResult)[0].Error)
	}
	resultArray := (*queryResult)[0].Result
	if len(resultArray) == 0 {
		return nil, fmt.Errorf("empty result array from create")
	}

	// If no OMIT needed, return directly
	if len(omit) == 0 {
		return cbor.Marshal(resultArray[0])
	}

	// Extract ID from result for the follow-up SELECT
	resultMap, ok := resultArray[0].(map[string]any)
	if !ok {
		return cbor.Marshal(resultArray[0])
	}
	if id, exists := resultMap["id"]; exists {
		selectID = fmt.Sprintf("%v", id)
	}
	if selectID == "" {
		return cbor.Marshal(resultArray[0])
	}

	// Do a follow-up SELECT with OMIT
	selectStatement := fmt.Sprintf("SELECT * OMIT %s FROM %s", strings.Join(omit, ", "), selectID)
	selectResult, err := surrealdb.Query[[]any](ctx, w.db, selectStatement, nil)
	if err != nil {
		return cbor.Marshal(resultArray[0])
	}
	if selectResult == nil || len(*selectResult) == 0 {
		return cbor.Marshal(resultArray[0])
	}
	if (*selectResult)[0].Error != nil {
		return cbor.Marshal(resultArray[0])
	}
	selectArray := (*selectResult)[0].Result
	if len(selectArray) == 0 {
		return cbor.Marshal(resultArray[0])
	}
	return cbor.Marshal(selectArray[0])
}

func (w *surrealDBWrapper) Select(ctx context.Context, id *ID, omit []string) ([]byte, error) {
	if id == nil {
		return nil, fmt.Errorf("id cannot be nil")
	}

	// Build SELECT query with OMIT if needed
	statement := fmt.Sprintf("SELECT * FROM %s", id.String())
	if len(omit) > 0 {
		statement = fmt.Sprintf("SELECT * OMIT %s FROM %s", strings.Join(omit, ", "), id.String())
	}

	queryResult, err := surrealdb.Query[[]any](ctx, w.db, statement, nil)
	if err != nil {
		return nil, err
	}
	if queryResult == nil || len(*queryResult) == 0 {
		return nil, fmt.Errorf("empty response from select")
	}
	if (*queryResult)[0].Error != nil {
		return nil, fmt.Errorf("select failed: %w", (*queryResult)[0].Error)
	}
	resultArray := (*queryResult)[0].Result
	if len(resultArray) == 0 {
		return nil, fmt.Errorf("record not found")
	}
	return cbor.Marshal(resultArray[0])
}

func (w *surrealDBWrapper) Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error) {
	result, err := surrealdb.Query[any](ctx, w.db, statement, vars)
	if err != nil {
		return nil, err
	}

	// Check for errors in individual query results.
	// The surrealdb.Query function returns *[]QueryResult[T], where each result can have its own error.
	if result != nil {
		for i, qr := range *result {
			if qr.Error != nil {
				return nil, fmt.Errorf("query statement %d failed: %w", i, qr.Error)
			}
		}
	}

	return cbor.Marshal(result)
}

func (w *surrealDBWrapper) Live(ctx context.Context, statement string, vars map[string]any) (<-chan []byte, error) {
	// NOTE: SurrealDB does not yet support proper variable handling for live queries.
	// To circumvent this limitation, params are registered in the database before issuing
	// the actual live query. Those params are given the values of the variables passed to
	// this method. This way, the live query can be filtered by said params.
	//
	// References:
	// Bug: Using variables in filters does not emit live messages (https://github.com/surrealdb/surrealdb/issues/2623)
	// Bug: LQ params should be evaluated before registering (https://github.com/surrealdb/surrealdb/issues/2641)
	// Bug: parameters do not work with live queries (https://github.com/surrealdb/surrealdb/issues/3602)
	// Feature: Live Query WHERE clause should process Params (https://github.com/surrealdb/surrealdb/issues/4026)

	// Generate a random prefix to prevent param name collisions.
	varPrefix, err := randString(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random string: %w", err)
	}

	// Create DEFINE PARAM statements for each variable.
	params := make(map[string]string, len(vars))
	for key := range vars {
		newKey := varPrefix + "_" + key
		params[newKey] = "DEFINE PARAM $" + newKey + " VALUE $" + key
		statement = strings.ReplaceAll(statement, "$"+key, "$"+newKey)
	}

	// Prepend DEFINE PARAM statements to the query.
	if len(params) > 0 {
		var paramDefs strings.Builder
		for _, value := range params {
			paramDefs.WriteString(value + "; ")
		}
		statement = paramDefs.String() + statement
	}

	// Execute the LIVE statement via Query call.
	result, err := surrealdb.Query[models.UUID](ctx, w.db, statement, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to execute live query: %w", err)
	}

	// Extract the live query ID from the result.
	// The last result contains the live query UUID.
	queryIndex := len(params)
	if result == nil || len(*result) <= queryIndex {
		return nil, fmt.Errorf("empty response from live query")
	}

	lastResult := (*result)[queryIndex]
	if lastResult.Error != nil {
		return nil, fmt.Errorf("live query error: %w", lastResult.Error)
	}

	liveID := lastResult.Result.String()

	notifications, err := w.db.LiveNotifications(liveID)
	if err != nil {
		return nil, fmt.Errorf("failed to get live notifications: %w", err)
	}

	out := make(chan []byte)
	go func() {
		defer close(out)
		defer func() {
			// Clean up the defined params when the live query ends.
			cleanupCtx := context.Background()
			for newKey := range params {
				_, _ = surrealdb.Query[any](cleanupCtx, w.db, fmt.Sprintf("REMOVE PARAM $%s;", newKey), nil)
			}
		}()

		for notif := range notifications {
			data, err := cbor.Marshal(notif)
			if err != nil {
				// TODO: add logger to overall client config and use it here
				slog.ErrorContext(ctx, "failed to marshal live notification", "error", err)
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- data:
			}
		}
	}()

	return out, nil
}

func (w *surrealDBWrapper) Update(ctx context.Context, id *ID, data any, omit []string) ([]byte, error) {
	if id == nil {
		return nil, fmt.Errorf("id cannot be nil")
	}

	// For UPDATE, we use MERGE (not CONTENT) so missing fields are preserved.
	// This is important for password fields that may be empty after OMIT reads.
	// CONTENT replaces the entire document, MERGE only updates provided fields.
	// We also use RETURN AFTER and then SELECT with OMIT if needed,
	// since SurrealDB doesn't support OMIT with CREATE/UPDATE RETURN clauses.
	statement := fmt.Sprintf("UPDATE %s MERGE $data RETURN AFTER", id.String())
	queryResult, err := surrealdb.Query[[]any](ctx, w.db, statement, map[string]any{"data": data})
	if err != nil {
		return nil, err
	}
	if queryResult == nil || len(*queryResult) == 0 {
		return nil, fmt.Errorf("empty response from update")
	}
	if (*queryResult)[0].Error != nil {
		return nil, fmt.Errorf("update failed: %w", (*queryResult)[0].Error)
	}
	resultArray := (*queryResult)[0].Result
	if len(resultArray) == 0 {
		return nil, fmt.Errorf("empty result array from update")
	}

	// If no OMIT needed, return directly
	if len(omit) == 0 {
		return cbor.Marshal(resultArray[0])
	}

	// Do a follow-up SELECT with OMIT
	selectStatement := fmt.Sprintf("SELECT * OMIT %s FROM %s", strings.Join(omit, ", "), id.String())
	selectResult, err := surrealdb.Query[[]any](ctx, w.db, selectStatement, nil)
	if err != nil {
		return cbor.Marshal(resultArray[0])
	}
	if selectResult == nil || len(*selectResult) == 0 {
		return cbor.Marshal(resultArray[0])
	}
	if (*selectResult)[0].Error != nil {
		return cbor.Marshal(resultArray[0])
	}
	selectArray := (*selectResult)[0].Result
	if len(selectArray) == 0 {
		return cbor.Marshal(resultArray[0])
	}
	return cbor.Marshal(selectArray[0])
}

func (w *surrealDBWrapper) Delete(ctx context.Context, id *ID) ([]byte, error) {
	if id == nil {
		return nil, fmt.Errorf("id cannot be nil")
	}

	result, err := surrealdb.Delete[any](ctx, w.db, *id)
	if err != nil {
		return nil, err
	}
	return cbor.Marshal(result)
}

func (w *surrealDBWrapper) Marshal(val any) ([]byte, error) {
	return cbor.Marshal(val)
}

func (w *surrealDBWrapper) Unmarshal(buf []byte, val any) error {
	return cbor.Unmarshal(buf, val)
}

func (w *surrealDBWrapper) Close() error {
	return w.db.Close(context.Background())
}

func NewClient(ctx context.Context, conf Config) (*ClientImpl, error) {
	db, err := surrealdb.FromEndpointURLString(ctx, conf.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SurrealDB: %w", err)
	}

	if conf.Username != "" && conf.Password != "" {
		_, err := db.SignIn(ctx, surrealdb.Auth{
			Username: conf.Username,
			Password: conf.Password,
		})
		if err != nil {
			_ = db.Close(ctx)
			return nil, fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	if conf.Namespace != "" && conf.Database != "" {
		_, err := surrealdb.Query[any](ctx, db,
			fmt.Sprintf("DEFINE NAMESPACE IF NOT EXISTS %s;", conf.Namespace),
			nil)
		if err != nil {
			_ = db.Close(ctx)
			return nil, fmt.Errorf("failed to create namespace: %w", err)
		}

		if err := db.Use(ctx, conf.Namespace, ""); err != nil {
			_ = db.Close(ctx)
			return nil, fmt.Errorf("failed to set namespace: %w", err)
		}

		_, err = surrealdb.Query[any](ctx, db,
			fmt.Sprintf("DEFINE DATABASE IF NOT EXISTS %s;", conf.Database),
			nil)
		if err != nil {
			_ = db.Close(ctx)
			return nil, fmt.Errorf("failed to create database: %w", err)
		}

		if err := db.Use(ctx, conf.Namespace, conf.Database); err != nil {
			_ = db.Close(ctx)
			return nil, fmt.Errorf("failed to set namespace/database: %w", err)
		}
	}

	wrapper := &surrealDBWrapper{db: db}

	return &ClientImpl{
		db: wrapper,
	}, nil
}

func (c *ClientImpl) Close() {
	_ = c.db.Close()
}

// randString generates a random alphanumeric string of length n.
func randString(n int) (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	byteSlice := make([]byte, n)

	for index := range byteSlice {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		byteSlice[index] = letterBytes[randInt.Int64()]
	}

	return string(byteSlice), nil
}
