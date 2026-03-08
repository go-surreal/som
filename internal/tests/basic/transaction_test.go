package basic

import (
	"context"
	"sync"
	"testing"

	som "github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestTransactionCommit(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	node := &model.SpecialTypes{Name: "tx-commit"}
	err := client.SpecialTypesRepo().Create(txCtx, node)
	assert.NilError(t, err)
	assert.Assert(t, node.ID() != "")

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)

	read, exists, err := client.SpecialTypesRepo().Read(ctx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "tx-commit", read.Name)
}

func TestTransactionCancel(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)

	node := &model.SpecialTypes{Name: "tx-cancel"}
	err := client.SpecialTypesRepo().Create(txCtx, node)
	assert.NilError(t, err)
	assert.Assert(t, node.ID() != "")

	cancel()

	_, exists, err := client.SpecialTypesRepo().Read(ctx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, !exists, "record should not exist after cancel")
}

func TestTransactionMultipleOperations(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	node1 := &model.SpecialTypes{Name: "tx-multi-1"}
	err := client.SpecialTypesRepo().Create(txCtx, node1)
	assert.NilError(t, err)

	node2 := &model.SpecialTypes{Name: "tx-multi-2"}
	err = client.SpecialTypesRepo().Create(txCtx, node2)
	assert.NilError(t, err)

	node3 := &model.SpecialTypes{Name: "tx-multi-3"}
	err = client.SpecialTypesRepo().Create(txCtx, node3)
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)

	read1, exists1, err := client.SpecialTypesRepo().Read(ctx, string(node1.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists1)
	assert.Equal(t, "tx-multi-1", read1.Name)

	read2, exists2, err := client.SpecialTypesRepo().Read(ctx, string(node2.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists2)
	assert.Equal(t, "tx-multi-2", read2.Name)

	read3, exists3, err := client.SpecialTypesRepo().Read(ctx, string(node3.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists3)
	assert.Equal(t, "tx-multi-3", read3.Name)
}

func TestTransactionMultipleRepos(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	special := &model.SpecialTypes{Name: "tx-cross-special"}
	err := client.SpecialTypesRepo().Create(txCtx, special)
	assert.NilError(t, err)

	relation := &model.SpecialRelation{Title: "tx-cross-relation"}
	err = client.SpecialRelationRepo().Create(txCtx, relation)
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)

	readSpecial, exists, err := client.SpecialTypesRepo().Read(ctx, string(special.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "tx-cross-special", readSpecial.Name)

	readRelation, exists, err := client.SpecialRelationRepo().Read(ctx, string(relation.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "tx-cross-relation", readRelation.Title)
}

func TestTransactionCancelMultipleRepos(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)

	special := &model.SpecialTypes{Name: "tx-cancel-special"}
	err := client.SpecialTypesRepo().Create(txCtx, special)
	assert.NilError(t, err)

	relation := &model.SpecialRelation{Title: "tx-cancel-relation"}
	err = client.SpecialRelationRepo().Create(txCtx, relation)
	assert.NilError(t, err)

	cancel()

	_, exists, err := client.SpecialTypesRepo().Read(ctx, string(special.ID()))
	assert.NilError(t, err)
	assert.Assert(t, !exists)

	_, exists, err = client.SpecialRelationRepo().Read(ctx, string(relation.ID()))
	assert.NilError(t, err)
	assert.Assert(t, !exists)
}

func TestTransactionCommitWithoutOperations(t *testing.T) {
	ctx := context.Background()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	err := som.TxCommit(txCtx)
	assert.NilError(t, err, "commit on unused transaction should be a no-op")
}

func TestTransactionCancelWithoutOperations(t *testing.T) {
	ctx := context.Background()

	_, cancel := som.TxStart(ctx)

	cancel()
}

func TestTransactionCommitThenCommitFails(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	node := &model.SpecialTypes{Name: "double-commit"}
	err := client.SpecialTypesRepo().Create(txCtx, node)
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.ErrorIs(t, err, som.ErrTransactionClosed)
}

func TestTransactionCancelAfterCommit(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	node := &model.SpecialTypes{Name: "commit-then-cancel"}
	err := client.SpecialTypesRepo().Create(txCtx, node)
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)

	err = som.TxCancel(txCtx)
	assert.NilError(t, err, "cancel after commit should be idempotent no-op")
}

func TestTransactionCancelIdempotent(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)

	node := &model.SpecialTypes{Name: "double-cancel"}
	err := client.SpecialTypesRepo().Create(txCtx, node)
	assert.NilError(t, err)

	cancel()

	err = som.TxCancel(txCtx)
	assert.NilError(t, err, "second cancel should be idempotent no-op")
}

func TestTransactionOperationAfterCommitFails(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	node := &model.SpecialTypes{Name: "op-after-commit"}
	err := client.SpecialTypesRepo().Create(txCtx, node)
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)

	node2 := &model.SpecialTypes{Name: "should-fail"}
	err = client.SpecialTypesRepo().Create(txCtx, node2)
	assert.Assert(t, err != nil, "operation after commit should fail")
}

func TestTransactionOperationAfterCancelFails(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)

	node := &model.SpecialTypes{Name: "op-after-cancel"}
	err := client.SpecialTypesRepo().Create(txCtx, node)
	assert.NilError(t, err)

	cancel()

	node2 := &model.SpecialTypes{Name: "should-fail"}
	err = client.SpecialTypesRepo().Create(txCtx, node2)
	assert.Assert(t, err != nil, "operation after cancel should fail")
}

func TestTxStartPanicsIfAlreadyActive(t *testing.T) {
	ctx := context.Background()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		som.TxStart(txCtx)
		return false
	}(), "expected panic when nesting transactions")
}

func TestTransactionReadBypassesCache(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := &model.SpecialTypes{Name: "cache-bypass"}
	err := client.SpecialTypesRepo().Create(ctx, node)
	assert.NilError(t, err)

	cachedCtx, cacheCleanup := som.WithCache[model.SpecialTypes](ctx)
	defer cacheCleanup()

	read1, exists, err := client.SpecialTypesRepo().Read(cachedCtx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "cache-bypass", read1.Name)

	txCachedCtx, txCancel := som.TxStart(cachedCtx)
	defer txCancel()

	read2, exists, err := client.SpecialTypesRepo().Read(txCachedCtx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "cache-bypass", read2.Name)

	assert.Assert(t, read1 != read2, "transaction read should bypass cache and return a different pointer")
}

func TestCommitWithoutTransaction(t *testing.T) {
	ctx := context.Background()

	err := som.TxCommit(ctx)
	assert.NilError(t, err, "commit on context without transaction should be no-op")
}

func TestCancelWithoutTransaction(t *testing.T) {
	ctx := context.Background()

	err := som.TxCancel(ctx)
	assert.NilError(t, err, "cancel on context without transaction should be no-op")
}

func TestTransactionConcurrentEnsureTx(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errs := make([]error, 10)

	for i := range 10 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			node := &model.SpecialTypes{Name: "concurrent"}
			errs[idx] = client.SpecialTypesRepo().Create(txCtx, node)
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		assert.NilError(t, err, "concurrent create %d failed", i)
	}

	err := som.TxCommit(txCtx)
	assert.NilError(t, err)
}

func TestTransactionSoftDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := &model.SpecialTypes{Name: "to-soft-delete"}
	err := client.SpecialTypesRepo().Create(ctx, node)
	assert.NilError(t, err)
	assert.Assert(t, !node.SoftDelete.IsDeleted())

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	err = client.SpecialTypesRepo().Delete(txCtx, node)
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)

	read, exists, err := client.SpecialTypesRepo().Read(ctx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Assert(t, read.SoftDelete.IsDeleted(), "record should be soft-deleted after committed transaction")
}

func TestTransactionSoftDeleteCancelled(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := &model.SpecialTypes{Name: "delete-cancel"}
	err := client.SpecialTypesRepo().Create(ctx, node)
	assert.NilError(t, err)

	txCtx, cancel := som.TxStart(ctx)

	err = client.SpecialTypesRepo().Delete(txCtx, node)
	assert.NilError(t, err)

	cancel()

	read, exists, err := client.SpecialTypesRepo().Read(ctx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Assert(t, !read.SoftDelete.IsDeleted(), "record should NOT be soft-deleted after cancelled transaction")
}
