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

DEFINE TABLE user SCHEMAFULL;
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
DEFINE FIELD login ON TABLE user TYPE object;
DEFINE FIELD login.username ON TABLE user TYPE string;
DEFINE FIELD login.password ON TABLE user TYPE string;
DEFINE FIELD role ON TABLE user TYPE string ASSERT $value INSIDE ["", "admin", "user"];
DEFINE FIELD groups ON TABLE user TYPE option<array>;
DEFINE FIELD groups.* ON TABLE user TYPE option<record(group)>;
DEFINE FIELD main_group ON TABLE user TYPE option<record(group)>;
DEFINE FIELD main_group_ptr ON TABLE user TYPE option<record(group)>;
DEFINE FIELD other ON TABLE user TYPE option<array>;
DEFINE FIELD other.* ON TABLE user TYPE string;
DEFINE FIELD more ON TABLE user TYPE option<array>;
DEFINE FIELD more.* ON TABLE user TYPE float;
DEFINE FIELD roles ON TABLE user TYPE option<array>;
DEFINE FIELD roles.* ON TABLE user TYPE string ASSERT $value INSIDE ["", "admin", "user"];
DEFINE FIELD string_ptr ON TABLE user TYPE option<string>;
DEFINE FIELD int_ptr ON TABLE user TYPE option<int>;
DEFINE FIELD struct_ptr ON TABLE user TYPE option<object>;
DEFINE FIELD struct_ptr.string_ptr ON TABLE user TYPE option<string>;
DEFINE FIELD struct_ptr.int_ptr ON TABLE user TYPE option<int>;
DEFINE FIELD struct_ptr.time_ptr ON TABLE user TYPE option<datetime>;
DEFINE FIELD struct_ptr.uuid_ptr ON TABLE user TYPE option<string> ASSERT $value == NONE OR $value == NULL OR is::uuid($value);
DEFINE FIELD string_ptr_slice ON TABLE user TYPE option<array>;
DEFINE FIELD string_ptr_slice.* ON TABLE user TYPE option<string>;
DEFINE FIELD string_slice_ptr ON TABLE user TYPE option<array>;
DEFINE FIELD string_slice_ptr.* ON TABLE user TYPE string;
DEFINE FIELD struct_ptr_slice ON TABLE user TYPE option<array>;
DEFINE FIELD struct_ptr_slice.* ON TABLE user TYPE option<object>;
DEFINE FIELD struct_ptr_slice.*.string_ptr ON TABLE user TYPE option<string>;
DEFINE FIELD struct_ptr_slice.*.int_ptr ON TABLE user TYPE option<int>;
DEFINE FIELD struct_ptr_slice.*.time_ptr ON TABLE user TYPE option<datetime>;
DEFINE FIELD struct_ptr_slice.*.uuid_ptr ON TABLE user TYPE option<string> ASSERT $value == NONE OR $value == NULL OR is::uuid($value);
DEFINE FIELD struct_ptr_slice_ptr ON TABLE user TYPE option<array>;
DEFINE FIELD struct_ptr_slice_ptr.* ON TABLE user TYPE option<object>;
DEFINE FIELD struct_ptr_slice_ptr.*.string_ptr ON TABLE user TYPE option<string>;
DEFINE FIELD struct_ptr_slice_ptr.*.int_ptr ON TABLE user TYPE option<int>;
DEFINE FIELD struct_ptr_slice_ptr.*.time_ptr ON TABLE user TYPE option<datetime>;
DEFINE FIELD struct_ptr_slice_ptr.*.uuid_ptr ON TABLE user TYPE option<string> ASSERT $value == NONE OR $value == NULL OR is::uuid($value);
DEFINE FIELD enum_ptr_slice ON TABLE user TYPE option<array>;
DEFINE FIELD enum_ptr_slice.* ON TABLE user TYPE option<string> ASSERT $value == NULL OR $value INSIDE ["", "admin", "user"];
DEFINE FIELD node_ptr_slice ON TABLE user TYPE option<array>;
DEFINE FIELD node_ptr_slice.* ON TABLE user TYPE option<record(group)>;
DEFINE FIELD node_ptr_slice_ptr ON TABLE user TYPE option<array>;
DEFINE FIELD node_ptr_slice_ptr.* ON TABLE user TYPE option<record(group)>;
DEFINE FIELD slice_slice ON TABLE user TYPE option<array>;
DEFINE FIELD slice_slice.* ON TABLE user TYPE option<array>;
DEFINE FIELD time ON TABLE user TYPE datetime;
DEFINE FIELD time_ptr ON TABLE user TYPE option<datetime>;
DEFINE FIELD duration ON TABLE user TYPE duration VALUE type::duration($value);
DEFINE FIELD duration_ptr ON TABLE user TYPE option<duration> VALUE type::duration($value);
DEFINE FIELD uuid ON TABLE user TYPE string ASSERT is::uuid($value);
DEFINE FIELD uuid_ptr ON TABLE user TYPE option<string> ASSERT $value == NONE OR $value == NULL OR is::uuid($value);

DEFINE TABLE group SCHEMAFULL;
DEFINE FIELD created_at ON TABLE group TYPE datetime;
DEFINE FIELD updated_at ON TABLE group TYPE datetime;
DEFINE FIELD name ON TABLE group TYPE string;

DEFINE TABLE group_member SCHEMAFULL;
DEFINE FIELD created_at ON TABLE group_member TYPE datetime;
DEFINE FIELD updated_at ON TABLE group_member TYPE datetime;
DEFINE FIELD meta ON TABLE group_member TYPE object;
DEFINE FIELD meta.is_admin ON TABLE group_member TYPE bool;
DEFINE FIELD meta.is_active ON TABLE group_member TYPE bool;

COMMIT TRANSACTION;
`
