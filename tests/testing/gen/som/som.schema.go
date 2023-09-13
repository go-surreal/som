// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package som

import(
	"context"
	"fmt"
)
	
func (c *ClientImpl) ApplySchema(ctx context.Context) error {
	_, err := c.db.Query(ctx, tmpl, nil)
	if err != nil {
		return fmt.Errorf("could not apply schema: %v", err)
	}

	return nil
}

var tmpl = `

DEFINE TABLE fields_like_db_response SCHEMAFULL;
DEFINE FIELD id ON TABLE fields_like_db_response TYPE record<fields_like_db_response> ASSERT $value != NONE AND $value != NULL AND $value != "";
DEFINE FIELD time ON TABLE fields_like_db_response TYPE string;
DEFINE FIELD status ON TABLE fields_like_db_response TYPE string;
DEFINE FIELD detail ON TABLE fields_like_db_response TYPE string;
DEFINE FIELD result ON TABLE fields_like_db_response TYPE option<array | null>;
DEFINE FIELD result.* ON TABLE fields_like_db_response TYPE string;
`