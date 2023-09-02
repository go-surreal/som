// Code generated by github.com/marcbinz/som, DO NOT EDIT.

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

DEFINE TABLE user SCHEMAFULL;
DEFINE FIELD id ON TABLE user TYPE record<user> ASSERT $value != NONE AND $value != NULL AND $value != "";
DEFINE FIELD created_at ON TABLE user TYPE datetime;
DEFINE FIELD updated_at ON TABLE user TYPE datetime;
DEFINE FIELD string ON TABLE user TYPE string;
DEFINE FIELD int ON TABLE user TYPE int;
DEFINE FIELD int_32 ON TABLE user TYPE int;
DEFINE FIELD int_64 ON TABLE user TYPE int;
DEFINE FIELD float_32 ON TABLE user TYPE float;
DEFINE FIELD float_64 ON TABLE user TYPE float;
DEFINE FIELD bool ON TABLE user TYPE bool;
DEFINE FIELD bool_2 ON TABLE user TYPE bool;
DEFINE FIELD uuid ON TABLE user TYPE string ASSERT is::uuid($value);
DEFINE FIELD login ON TABLE user TYPE object;
DEFINE FIELD login.username ON TABLE user TYPE string;
DEFINE FIELD login.password ON TABLE user TYPE string;
DEFINE FIELD role ON TABLE user TYPE string ASSERT $value INSIDE ["", "admin", "user"];
DEFINE FIELD groups ON TABLE user TYPE option<array | null>;
DEFINE FIELD groups.* ON TABLE user TYPE option<record<group> | null>;
DEFINE FIELD main_group ON TABLE user TYPE option<record<group> | null>;
DEFINE FIELD main_group_ptr ON TABLE user TYPE option<record<group> | null>;
DEFINE FIELD other ON TABLE user TYPE option<array | null>;
DEFINE FIELD other.* ON TABLE user TYPE string;
DEFINE FIELD more ON TABLE user TYPE option<array | null>;
DEFINE FIELD more.* ON TABLE user TYPE float;
DEFINE FIELD roles ON TABLE user TYPE option<array | null>;
DEFINE FIELD roles.* ON TABLE user TYPE string ASSERT $value INSIDE ["", "admin", "user"];
DEFINE FIELD string_ptr ON TABLE user TYPE option<string | null>;
DEFINE FIELD int_ptr ON TABLE user TYPE option<int | null>;
DEFINE FIELD time_ptr ON TABLE user TYPE option<datetime | null>;
DEFINE FIELD uuid_ptr ON TABLE user TYPE option<string | null> ASSERT $value == NONE OR $value == NULL OR is::uuid($value);
DEFINE FIELD struct_ptr ON TABLE user TYPE option<object | null>;
DEFINE FIELD struct_ptr.string_ptr ON TABLE user TYPE option<string | null>;
DEFINE FIELD struct_ptr.int_ptr ON TABLE user TYPE option<int | null>;
DEFINE FIELD struct_ptr.time_ptr ON TABLE user TYPE option<datetime | null>;
DEFINE FIELD struct_ptr.uuid_ptr ON TABLE user TYPE option<string | null> ASSERT $value == NONE OR $value == NULL OR is::uuid($value);
DEFINE FIELD string_ptr_slice ON TABLE user TYPE option<array | null>;
DEFINE FIELD string_ptr_slice.* ON TABLE user TYPE option<string | null>;
DEFINE FIELD string_slice_ptr ON TABLE user TYPE option<array | null>;
DEFINE FIELD string_slice_ptr.* ON TABLE user TYPE string;
DEFINE FIELD struct_ptr_slice ON TABLE user TYPE option<array | null>;
DEFINE FIELD struct_ptr_slice.* ON TABLE user TYPE option<object | null>;
DEFINE FIELD struct_ptr_slice.*.string_ptr ON TABLE user TYPE option<string | null>;
DEFINE FIELD struct_ptr_slice.*.int_ptr ON TABLE user TYPE option<int | null>;
DEFINE FIELD struct_ptr_slice.*.time_ptr ON TABLE user TYPE option<datetime | null>;
DEFINE FIELD struct_ptr_slice.*.uuid_ptr ON TABLE user TYPE option<string | null> ASSERT $value == NONE OR $value == NULL OR is::uuid($value);
DEFINE FIELD struct_ptr_slice_ptr ON TABLE user TYPE option<array | null>;
DEFINE FIELD struct_ptr_slice_ptr.* ON TABLE user TYPE option<object | null>;
DEFINE FIELD struct_ptr_slice_ptr.*.string_ptr ON TABLE user TYPE option<string | null>;
DEFINE FIELD struct_ptr_slice_ptr.*.int_ptr ON TABLE user TYPE option<int | null>;
DEFINE FIELD struct_ptr_slice_ptr.*.time_ptr ON TABLE user TYPE option<datetime | null>;
DEFINE FIELD struct_ptr_slice_ptr.*.uuid_ptr ON TABLE user TYPE option<string | null> ASSERT $value == NONE OR $value == NULL OR is::uuid($value);
DEFINE FIELD enum_ptr_slice ON TABLE user TYPE option<array | null>;
DEFINE FIELD enum_ptr_slice.* ON TABLE user TYPE option<string> ASSERT $value == NULL OR $value INSIDE ["", "admin", "user"];
DEFINE FIELD node_ptr_slice ON TABLE user TYPE option<array | null>;
DEFINE FIELD node_ptr_slice.* ON TABLE user TYPE option<record<group> | null>;
DEFINE FIELD node_ptr_slice_ptr ON TABLE user TYPE option<array | null>;
DEFINE FIELD node_ptr_slice_ptr.* ON TABLE user TYPE option<record<group> | null>;
DEFINE FIELD slice_slice ON TABLE user TYPE option<array | null>;
DEFINE FIELD slice_slice.* ON TABLE user TYPE option<array | null>;

DEFINE TABLE group SCHEMAFULL;
DEFINE FIELD id ON TABLE group TYPE record<group> ASSERT $value != NONE AND $value != NULL AND $value != "";
DEFINE FIELD created_at ON TABLE group TYPE datetime;
DEFINE FIELD updated_at ON TABLE group TYPE datetime;
DEFINE FIELD name ON TABLE group TYPE string;

DEFINE TABLE group_member SCHEMAFULL;
DEFINE FIELD created_at ON TABLE group_member TYPE datetime;
DEFINE FIELD updated_at ON TABLE group_member TYPE datetime;
DEFINE FIELD meta ON TABLE group_member TYPE object;
DEFINE FIELD meta.is_admin ON TABLE group_member TYPE bool;
DEFINE FIELD meta.is_active ON TABLE group_member TYPE bool;
`
