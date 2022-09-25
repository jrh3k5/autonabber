package undo

import "time"

type BudgetChangesUndo struct {
	ApplicationDateTime time.Time
	AppliedBudgetID     string
	ChangeName          string
	Groups              []*BudgetGroupChangeUndo
}

type BudgetGroupChangeUndo struct {
	GroupName  string
	Categories []*BudgetCategoryChangeUndo
}

type BudgetCategoryChangeUndo struct {
	CategoryName       string
	DollarChangeAmount int64
	CentChangeAmount   int16
}
