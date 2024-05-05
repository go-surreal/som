package som

type DefineRoot struct{}

func Define() *DefineRoot {
	return &DefineRoot{}
}

func (d *DefineRoot) Analyzer() *DefineAnalyzer {
	return &DefineAnalyzer{}
}

func (d *DefineRoot) Index() *DefineIndex {
	return &DefineIndex{}
}

func (d *DefineRoot) Constraint() *DefineConstraint {
	return &DefineConstraint{}
}

type DefineAnalyzer struct{}

type DefineIndex struct{}

type DefineConstraint struct{}
