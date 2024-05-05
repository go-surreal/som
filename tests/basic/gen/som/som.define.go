package som

type DefineRoot struct{}

func Define() *DefineRoot {
	return &DefineRoot{}
}

func (d *DefineRoot) Analyzer(name string) *DefineAnalyzer {
	return &DefineAnalyzer{}
}

func (d *DefineRoot) Constraint() *DefineConstraint {
	return &DefineConstraint{}
}

type DefineAnalyzer struct{}

func (d *DefineAnalyzer) Tokenizers(a ...any) *DefineAnalyzer {
	return &DefineAnalyzer{}
}

func (d *DefineAnalyzer) Filters(a ...any) *DefineAnalyzer {
	return &DefineAnalyzer{}
}

type DefineConstraint struct{}

func (d *DefineRoot) Model() *DefineModel {
	return &DefineModel{}
}

type DefineModel struct{}

func (d *DefineModel) User() *DefineTable[any] {
	return &DefineTable[any]{}
}

type DefineTable[T any] struct{}

func (d *DefineTable[T]) Index(name string) *DefineIndex[T] {
	return &DefineIndex[T]{}
}

type DefineIndex[T any] struct{}

func (d *DefineIndex[T]) On(a ...any) *DefineIndex[T] {
	return &DefineIndex[T]{}
}

func (d *DefineIndex[T]) Unique() *DefineIndex[T] {
	return &DefineIndex[T]{}
}

func (d *DefineIndex[T]) Search(a ...any) *DefineIndex[T] {
	return &DefineIndex[T]{}
}
