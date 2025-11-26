package field

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestDatabaseObject(t *testing.T) {
	t.Parallel()

	dbObj := DatabaseObject{
		Name:   "SomeNameCAPS",
		Fields: nil,
	}

	assert.Equal(t, "SomeNameCAPS", dbObj.NameGo())
	assert.Equal(t, "someNameCaps", dbObj.NameGoLower())
	assert.Equal(t, "some_name_caps", dbObj.NameDatabase())
	assert.Equal(t, "object.some_name_caps.go", dbObj.FileName())
	assert.Check(t, dbObj.GetFields() == nil)
}
