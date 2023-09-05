package field

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestNodeTable(t *testing.T) {
	nodeTable := NodeTable{
		Name: "User",
	}

	assert.Equal(t, "node.user.go", nodeTable.FileName())
	assert.Equal(t, "User", nodeTable.NameGo())
	assert.Equal(t, "user", nodeTable.NameGoLower())
	assert.Equal(t, "user", nodeTable.NameDatabase())

	nodeTable = NodeTable{
		Name:       "Account",
		Timestamps: true,
	}

	assert.Equal(t, "node.account.go", nodeTable.FileName())
	assert.Equal(t, "Account", nodeTable.NameGo())
	assert.Equal(t, "account", nodeTable.NameGoLower())
	assert.Equal(t, "account", nodeTable.NameDatabase())
	assert.Equal(t, true, nodeTable.HasTimestamps())

	nodeTable = NodeTable{
		Name: "FieldsLikeDBResponse",
	}

	assert.Equal(t, "node.fields_like_db_response.go", nodeTable.FileName())
	assert.Equal(t, "FieldsLikeDBResponse", nodeTable.NameGo())
	assert.Equal(t, "fieldsLikeDBResponse", nodeTable.NameGoLower())
	assert.Equal(t, "fields_like_db_response", nodeTable.NameDatabase())
}
