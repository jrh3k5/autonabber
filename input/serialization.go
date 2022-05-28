package input

type serializedBudgetChanges struct {
	Changes []*serializedBudgetChange `yaml:"changes"`
}

type serializedBudgetChange struct {
	Name           string                           `yaml:"name"`
	CategoryGroups []*serializedBudgetCategoryGroup `yaml:"category_groups"`
}

type serializedBudgetCategoryGroup struct {
	Name       string                            `yaml:"name"`
	Categories []*serializedBudgetCategoryChange `yaml:"categories"`
}

type serializedBudgetCategoryChange struct {
	Name   string `yaml:"name"`
	Change string `yaml:"change"`
}
