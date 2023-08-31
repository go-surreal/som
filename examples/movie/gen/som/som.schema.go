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

DEFINE TABLE person SCHEMAFULL;
DEFINE FIELD id ON TABLE person TYPE record<person> ASSERT $value != NONE AND $value != NULL AND $value != "";
DEFINE FIELD name ON TABLE person TYPE string;

DEFINE TABLE movie SCHEMAFULL;
DEFINE FIELD id ON TABLE movie TYPE record<movie> ASSERT $value != NONE AND $value != NULL AND $value != "";
DEFINE FIELD title ON TABLE movie TYPE string;

DEFINE TABLE directed SCHEMAFULL;

DEFINE TABLE acted_in SCHEMAFULL;
`
