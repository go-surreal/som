package basic

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestIterate(t *testing.T) {
	ctx := context.Background()
	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	t.Run("empty result set", func(t *testing.T) {
		count := 0
		for _, err := range client.SpecialTypesRepo().Query().Iterate(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			count++
		}
		assert.Equal(t, 0, count)
	})

	t.Run("single batch", func(t *testing.T) {
		// Create 3 records, batch size 10
		for i := 0; i < 3; i++ {
			g := &model.SpecialTypes{Name: fmt.Sprintf("single-batch-%d", i)}
			err := client.SpecialTypesRepo().Create(ctx, g)
			assert.NilError(t, err)
		}

		count := 0
		for _, err := range client.SpecialTypesRepo().Query().
			Where(filter.SpecialTypes.Name.Contains("single-batch-").True()).
			Iterate(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			count++
		}
		assert.Equal(t, 3, count)
	})

	t.Run("multiple batches", func(t *testing.T) {
		// Create 25 records, batch size 10 -> 3 batches
		for i := 0; i < 25; i++ {
			g := &model.SpecialTypes{Name: fmt.Sprintf("multi-batch-%d", i)}
			err := client.SpecialTypesRepo().Create(ctx, g)
			assert.NilError(t, err)
		}

		count := 0
		for _, err := range client.SpecialTypesRepo().Query().
			Where(filter.SpecialTypes.Name.Contains("multi-batch-").True()).
			Iterate(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			count++
		}
		assert.Equal(t, 25, count)
	})

	t.Run("exact batch boundary", func(t *testing.T) {
		// Create exactly 20 records, batch size 10 -> 2 full batches
		for i := 0; i < 20; i++ {
			g := &model.SpecialTypes{Name: fmt.Sprintf("exact-batch-%d", i)}
			err := client.SpecialTypesRepo().Create(ctx, g)
			assert.NilError(t, err)
		}

		count := 0
		for _, err := range client.SpecialTypesRepo().Query().
			Where(filter.SpecialTypes.Name.Contains("exact-batch-").True()).
			Iterate(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			count++
		}
		assert.Equal(t, 20, count)
	})

	t.Run("early termination", func(t *testing.T) {
		// Create 15 records
		for i := 0; i < 15; i++ {
			g := &model.SpecialTypes{Name: fmt.Sprintf("early-term-%d", i)}
			err := client.SpecialTypesRepo().Create(ctx, g)
			assert.NilError(t, err)
		}

		count := 0
		for _, err := range client.SpecialTypesRepo().Query().
			Where(filter.SpecialTypes.Name.Contains("early-term-").True()).
			Iterate(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			count++
			if count >= 5 {
				break
			}
		}
		assert.Equal(t, 5, count)
	})
}

func TestIterateID(t *testing.T) {
	ctx := context.Background()
	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	t.Run("empty result set", func(t *testing.T) {
		count := 0
		for _, err := range client.SpecialTypesRepo().Query().IterateID(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			count++
		}
		assert.Equal(t, 0, count)
	})

	t.Run("single batch", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			g := &model.SpecialTypes{Name: fmt.Sprintf("id-single-batch-%d", i)}
			err := client.SpecialTypesRepo().Create(ctx, g)
			assert.NilError(t, err)
		}

		count := 0
		for id, err := range client.SpecialTypesRepo().Query().
			Where(filter.SpecialTypes.Name.Contains("id-single-batch-").True()).
			IterateID(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			assert.Assert(t, id != "", "ID should not be empty")
			count++
		}
		assert.Equal(t, 3, count)
	})

	t.Run("multiple batches", func(t *testing.T) {
		for i := 0; i < 25; i++ {
			g := &model.SpecialTypes{Name: fmt.Sprintf("id-multi-batch-%d", i)}
			err := client.SpecialTypesRepo().Create(ctx, g)
			assert.NilError(t, err)
		}

		count := 0
		for id, err := range client.SpecialTypesRepo().Query().
			Where(filter.SpecialTypes.Name.Contains("id-multi-batch-").True()).
			IterateID(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			assert.Assert(t, id != "", "ID should not be empty")
			count++
		}
		assert.Equal(t, 25, count)
	})

	t.Run("ID format check", func(t *testing.T) {
		g := &model.SpecialTypes{Name: "id-format-check"}
		err := client.SpecialTypesRepo().Create(ctx, g)
		assert.NilError(t, err)

		for id, err := range client.SpecialTypesRepo().Query().
			Where(filter.SpecialTypes.Name.Contains("id-format-check").True()).
			IterateID(ctx, 10) {
			if err != nil {
				t.Fatal(err)
			}
			assert.Assert(t, len(id) > 0, "ID should not be empty")
			assert.Assert(t, id[:14] == "special_types:", "ID should start with 'special_types:', got: %s", id)
		}
	})
}
