// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package som

import(
	"fmt"
)
	
func (c *ClientImpl) ApplySchema() error {
	_, err := c.db.Query(tmpl, nil)
	if err != nil {
		return fmt.Errorf("could not apply schema: %v", err)
	}

	return nil
}

var tmpl = `

BEGIN TRANSACTION;

DEFINE TABLE fields_like_db_response SCHEMAFULL;
DEFINE FIELD id ON TABLE fields_like_db_response TYPE record ASSERT $value != NONE AND $value != NULL AND $value != "";
DEFINE FIELD time ON TABLE fields_like_db_response TYPE string;
DEFINE FIELD status ON TABLE fields_like_db_response TYPE string;
DEFINE FIELD detail ON TABLE fields_like_db_response TYPE string;
DEFINE FIELD result ON TABLE fields_like_db_response TYPE option<array>;
DEFINE FIELD result.* ON TABLE fields_like_db_response TYPE string;

COMMIT TRANSACTION;
`
