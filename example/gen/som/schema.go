// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package som

import(
	"fmt"
)
	
func (c *Client) ApplySchema() error {
	_, err := c.db.Query(tmpl, nil)
	if err != nil {
		return fmt.Errorf("could not apply schema: %v", err)
	}

	return nil
}

var tmpl = `

DEFINE TABLE user SCHEMAFULL;
DEFINE FIELD created_at ON TABLE user TYPE datetime ASSERT $value != NULL;
DEFINE FIELD updated_at ON TABLE user TYPE datetime ASSERT $value != NULL;
DEFINE FIELD string ON TABLE user TYPE string ASSERT $value != NULL;
DEFINE FIELD int ON TABLE user TYPE int ASSERT $value != NULL;
DEFINE FIELD int_32 ON TABLE user TYPE int ASSERT $value != NULL;
DEFINE FIELD int_64 ON TABLE user TYPE int ASSERT $value != NULL;
DEFINE FIELD float_32 ON TABLE user TYPE float ASSERT $value != NULL;
DEFINE FIELD float_64 ON TABLE user TYPE float ASSERT $value != NULL;
DEFINE FIELD bool ON TABLE user TYPE bool ASSERT $value != NULL;
DEFINE FIELD bool_2 ON TABLE user TYPE bool ASSERT $value != NULL;
DEFINE FIELD uuid ON TABLE user TYPE string ASSERT $value != NULL;
DEFINE FIELD login ON TABLE user TYPE object ASSERT $value != NULL;
DEFINE FIELD login.username ON TABLE user TYPE string ASSERT $value != NULL;
DEFINE FIELD login.password ON TABLE user TYPE string ASSERT $value != NULL;
DEFINE FIELD role ON TABLE user TYPE string ASSERT $value != NULL;
DEFINE FIELD groups ON TABLE user TYPE array;
DEFINE FIELD groups.* ON TABLE user TYPE record(group);
DEFINE FIELD main_group ON TABLE user TYPE record(group);
DEFINE FIELD other ON TABLE user TYPE array;
DEFINE FIELD other.* ON TABLE user TYPE string ASSERT $value != NULL;
DEFINE FIELD more ON TABLE user TYPE array;
DEFINE FIELD more.* ON TABLE user TYPE float ASSERT $value != NULL;
DEFINE FIELD roles ON TABLE user TYPE array;
DEFINE FIELD roles.* ON TABLE user TYPE string ASSERT $value != NULL;
DEFINE FIELD string_ptr ON TABLE user TYPE string;
DEFINE FIELD int_ptr ON TABLE user TYPE int;
DEFINE FIELD time_ptr ON TABLE user TYPE datetime;
DEFINE FIELD uuid_ptr ON TABLE user TYPE string;
DEFINE FIELD struct_ptr ON TABLE user TYPE object;
DEFINE FIELD struct_ptr.string_ptr ON TABLE user TYPE string;
DEFINE FIELD struct_ptr.int_ptr ON TABLE user TYPE int;
DEFINE FIELD struct_ptr.time_ptr ON TABLE user TYPE datetime;
DEFINE FIELD struct_ptr.uuid_ptr ON TABLE user TYPE string;
DEFINE FIELD string_ptr_slice ON TABLE user TYPE array;
DEFINE FIELD string_ptr_slice.* ON TABLE user TYPE string;
DEFINE FIELD string_slice_ptr ON TABLE user TYPE array;
DEFINE FIELD string_slice_ptr.* ON TABLE user TYPE string ASSERT $value != NULL;
DEFINE FIELD struct_ptr_slice ON TABLE user TYPE array;
DEFINE FIELD struct_ptr_slice.* ON TABLE user TYPE object;
DEFINE FIELD struct_ptr_slice.*.string_ptr ON TABLE user TYPE string;
DEFINE FIELD struct_ptr_slice.*.int_ptr ON TABLE user TYPE int;
DEFINE FIELD struct_ptr_slice.*.time_ptr ON TABLE user TYPE datetime;
DEFINE FIELD struct_ptr_slice.*.uuid_ptr ON TABLE user TYPE string;
DEFINE FIELD struct_ptr_slice_ptr ON TABLE user TYPE array;
DEFINE FIELD struct_ptr_slice_ptr.* ON TABLE user TYPE object;
DEFINE FIELD struct_ptr_slice_ptr.*.string_ptr ON TABLE user TYPE string;
DEFINE FIELD struct_ptr_slice_ptr.*.int_ptr ON TABLE user TYPE int;
DEFINE FIELD struct_ptr_slice_ptr.*.time_ptr ON TABLE user TYPE datetime;
DEFINE FIELD struct_ptr_slice_ptr.*.uuid_ptr ON TABLE user TYPE string;
DEFINE FIELD enum_ptr_slice ON TABLE user TYPE array;
DEFINE FIELD enum_ptr_slice.* ON TABLE user TYPE string;
DEFINE FIELD node_ptr_slice ON TABLE user TYPE array;
DEFINE FIELD node_ptr_slice.* ON TABLE user TYPE record(group);
DEFINE FIELD node_ptr_slice_ptr ON TABLE user TYPE array;
DEFINE FIELD node_ptr_slice_ptr.* ON TABLE user TYPE record(group);
DEFINE FIELD slice_slice ON TABLE user TYPE array;
DEFINE FIELD slice_slice.* ON TABLE user TYPE array;

DEFINE TABLE group SCHEMAFULL;
DEFINE FIELD created_at ON TABLE group TYPE datetime ASSERT $value != NULL;
DEFINE FIELD updated_at ON TABLE group TYPE datetime ASSERT $value != NULL;
DEFINE FIELD name ON TABLE group TYPE string ASSERT $value != NULL;

DEFINE TABLE member_of SCHEMAFULL;
DEFINE FIELD created_at ON TABLE member_of TYPE datetime ASSERT $value != NULL;
DEFINE FIELD updated_at ON TABLE member_of TYPE datetime ASSERT $value != NULL;
DEFINE FIELD meta ON TABLE member_of TYPE object ASSERT $value != NULL;
DEFINE FIELD meta.is_admin ON TABLE member_of TYPE bool ASSERT $value != NULL;
DEFINE FIELD meta.is_active ON TABLE member_of TYPE bool ASSERT $value != NULL;
`
