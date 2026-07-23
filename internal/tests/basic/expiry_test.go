package basic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"som.test/gen/som/repo"
	"som.test/model"
	"gotest.tools/v3/assert"
)

func TestExpiry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// expires_at is set by the database on creation and lies in the future.
	eph := model.Ephemeral{Label: "temp"}
	err := client.EphemeralRepo().Create(ctx, &eph)
	assert.NilError(t, err)
	assert.Assert(t, eph.ID() != "")
	assert.Assert(t, !eph.Expiry.ExpiresAt().IsZero(), "expires_at should be populated on create")
	assert.Assert(t, eph.Expiry.ExpiresAt().After(time.Now()), "expires_at should be in the future")

	all, err := client.EphemeralRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(all), "record should be visible before expiry")

	time.Sleep(1500 * time.Millisecond)

	// Filter-on-read: expired records are excluded from normal queries.
	afterExpiry, err := client.EphemeralRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(afterExpiry), "expired record should be filtered out")

	// The record still physically exists until purged; WithExpired reveals it.
	withExpired, err := client.EphemeralRepo().Query().WithExpired().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(withExpired), "expired record should still exist until purged")
}

// TestExpiryPurge verifies the background purge goroutine deletes expired records.
func TestExpiryPurge(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	id := dbCounter.Add(1)
	namespace := fmt.Sprintf("ns_expiry_%d", id)
	database := fmt.Sprintf("db_expiry_%d", id)

	client, err := repo.NewClient(ctx, repo.Config{
		Address:             "ws://" + sharedEndpoint,
		Username:            sharedUsername,
		Password:            sharedPassword,
		Namespace:           namespace,
		Database:            database,
		ExpiryPurgeInterval: 250 * time.Millisecond,
	})
	assert.NilError(t, err)
	defer func() {
		_, _ = client.Raw(context.Background(), fmt.Sprintf("REMOVE DATABASE %s", database), nil)
		_, _ = client.Raw(context.Background(), fmt.Sprintf("REMOVE NAMESPACE %s", namespace), nil)
		client.Close()
	}()

	err = client.ApplySchema(ctx)
	assert.NilError(t, err)

	eph := model.Ephemeral{Label: "purge-me"}
	err = client.EphemeralRepo().Create(ctx, &eph)
	assert.NilError(t, err)

	// Wait for expiry (1s) plus at least one purge tick (250ms).
	time.Sleep(2 * time.Second)

	// Even including expired records, the purge goroutine should have removed it.
	remaining, err := client.EphemeralRepo().Query().WithExpired().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(remaining), "expired record should be purged from the database")
}
