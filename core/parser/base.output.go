package parser

// Output is the complete result of parsing a model package.
type Output struct {
	PkgPath    string
	Nodes      []*Node
	Unions     []*Union
	Edges      []*Edge
	Views      []*View
	Sinks      []*Sink
	Fragments  []*Fragment
	Structs    []*Struct
	Enums      []*Enum
	EnumValues []*EnumValue

	Define *DefineOutput

	UsedFeatures *UsedFeatures
}

// UsedFeatures records which optional third-party packages the model uses, so
// that the generated code only pulls in what is actually needed.
type UsedFeatures struct {
	UsesGoogleUUID        bool
	UsesGofrsUUID         bool
	UsesStdUUID           bool
	UsesOrbGeo            bool
	UsesSimplefeaturesGeo bool
	UsesGoGeomGeo         bool
}

// featureContributor is implemented by the field types that depend on an
// optional third-party package.
type featureContributor interface {
	contributeFeatures(features *UsedFeatures)
}

func collectUsedFeatures(output *Output) *UsedFeatures {
	features := &UsedFeatures{}

	collect := func(fields []Field) {
		for _, f := range fields {
			if c, ok := f.(featureContributor); ok {
				c.contributeFeatures(features)
			}
		}
	}

	for _, node := range output.Nodes {
		collect(node.Fields)
	}

	for _, edge := range output.Edges {
		collect(edge.Fields)
	}

	for _, str := range output.Structs {
		collect(str.Fields)
	}

	return features
}
