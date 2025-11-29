//go:build embed

package repo

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed schema/schema.surql
var schema string

func (c *ClientImpl) ApplySchema(ctx context.Context) error {
	_, err := c.db.Query(ctx, schema, nil)
	if err != nil {
		return fmt.Errorf("could not apply schema: %v", err)
	}

	return nil
}
