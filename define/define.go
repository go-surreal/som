package define

func test() {
	Def()
}

func Define() *Define {
	return &Define{}
}

type Def struct{}

func (d *Def) Table(name string) *DefTable {
	return &DefTable{
		statement: "DEFINE TABLE " + name + " SCHEMAFULL",
	}
}

type DefTable struct {
	statement string
}

func (t *DefTable) WithFields() string {}

func (d *Def) Field(name string) *DefField {}

type DefField struct {
}

// DEFINE [
//	NAMESPACE @name
//	| DATABASE @name
//	| LOGIN @name ON [ NAMESPACE | DATABASE ] [ PASSWORD @pass | PASSHASH @hash ]
//	| TOKEN @name ON [ NAMESPACE | DATABASE ] TYPE @type VALUE @value
//	| SCOPE @name
//	| TABLE @name
//		[ DROP ]
//		[ SCHEMAFULL | SCHEMALESS ]
//		[ AS SELECT @projections
//			FROM @tables
//			[ WHERE @condition ]
//			[ GROUP [ BY ] @groups ]
//		]
//		[ PERMISSIONS [ NONE | FULL
//			| FOR select @expression
//			| FOR create @expression
//			| FOR update @expression
//			| FOR delete @expression
//		] ]
//	| EVENT @name ON [ TABLE ] @table WHEN @expression THEN @expression
//	| FIELD @name ON [ TABLE ] @table
//		[ TYPE @type ]
//		[ VALUE @expression ]
//		[ ASSERT @expression ]
//		[ PERMISSIONS [ NONE | FULL
//			| FOR select @expression
//			| FOR create @expression
//			| FOR update @expression
//			| FOR delete @expression
//		] ]
//	| INDEX @name ON [ TABLE ] @table [ FIELDS | COLUMNS ] @fields [ UNIQUE ]
// ]
