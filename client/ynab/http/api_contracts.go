package http

type budgetsResponse struct {
	Data *budgetData `json:"data"`
}

type budgetData struct {
	Budgets []*budgetDetails `json:"budgets"`
}

type budgetDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type budgetCategoriesResponse struct {
	Data *budgetCategoriesData `json:"data"`
}

type budgetCategoriesData struct {
	CategoryGroups []*budgetCategoryGroup `json:"category_groups"`
}

type budgetCategoryGroup struct {
	Name       string            `json:"name"`
	Categories []*budgetCategory `json:"categories"`
	Hidden     bool              `json:"hidden"`
}

type budgetCategory struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Budgeted int64  `json:"budgeted"`
	Balance  int64  `json:"balance"`
	Hidden   bool   `json:"hidden"`
}

type categoryPatchRequest struct {
	Category *patchedCategory `json:"category"`
}

type patchedCategory struct {
	Budgeted int64 `json:"budgeted"`
}

type transactionsContainer struct {
	Data *transactions `json:"data"`
}

type transactions struct {
	Transactions []*transaction `json:"transactions"`
}

type transaction struct {
	Date   string `json:"date"`
	Amount int64  `json:"amount"`
}
