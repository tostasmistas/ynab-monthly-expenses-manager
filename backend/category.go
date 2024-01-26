package backend

import (
	"fmt"
	"strings"
)

// CategoryGroup represents a YNAB category group
// This struct corresponds to the data structure defined in the YNAB API documentation
type CategoryGroup struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Hidden  bool   `json:"hidden"`
	Deleted bool   `json:"deleted"`
}

// Category represents a YNAB category
// This struct corresponds to the data structure defined in the YNAB API documentation
type Category struct {
	Id                     string `json:"id"`
	CategoryGroupId        string `json:"category_group_id"`
	CategoryGroupName      string `json:"category_group_name"`
	Name                   string `json:"name"`
	Hidden                 bool   `json:"hidden"`
	Note                   string `json:"note"`
	Budgeted               int64  `json:"budgeted"`
	Activity               int64  `json:"activity"`
	Balance                int64  `json:"balance"`
	GoalType               string `json:"goal_type"`
	GoalDay                int32  `json:"goal_day"`
	GoalCadence            int32  `json:"goal_cadence"`
	GoalCadenceFrequency   int32  `json:"goal_cadence_frequency"`
	GoalCreationMonth      string `json:"goal_creation_month"`
	GoalTarget             int64  `json:"goal_target"`
	GoalTargetMonth        string `json:"goal_target_month"`
	GoalPercentageComplete int32  `json:"goal_percentage_complete"`
	GoalMonthsToBudget     int32  `json:"goal_months_to_budget"`
	GoalUnderFunded        int64  `json:"goal_under_funded"`
	GoalOverallFunded      int64  `json:"goal_overall_funded"`
	GoalOverallLeft        int64  `json:"goal_overall_left"`
	Deleted                bool   `json:"deleted"`
}

// CategoryGroupWithCategories represents a YNAB category group with its categories
// This struct corresponds to the data structure defined in the YNAB API documentation
type CategoryGroupWithCategories struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Hidden     bool       `json:"hidden"`
	Deleted    bool       `json:"deleted"`
	Categories []Category `json:"categories"`
}

// CategoryGroupsWithCategories represents a collection of YNAB category groups with their categories.
type CategoryGroupsWithCategories []CategoryGroupWithCategories

// GetCategories fetches the YNAB categories of a YNAB budget
// GET https://api.ynab.com/v1/budgets/{budget_id}/categories
func (client *APIClient) GetCategories(budgetId string) (CategoryGroupsWithCategories, error) {
	categoriesResponse := struct {
		Data struct {
			CategoryGroups  CategoryGroupsWithCategories `json:"category_groups"`
			ServerKnowledge int64                        `json:"server_knowledge"`
		} `json:"data"`
	}{}

	response, err := client.Client.R().
		SetResult(&categoriesResponse).
		Get(fmt.Sprintf("budgets/%s/categories", budgetId))

	if err = client.ValidateResponse(response, err); err != nil {
		return nil, err
	}

	return categoriesResponse.Data.CategoryGroups, nil
}

// GetMonthlyExpensesCategories fetches the YNAB categories related to monthly expenses
func (categoryGroups *CategoryGroupsWithCategories) GetMonthlyExpensesCategories() []Category {
	var monthlyExpensesCategories []Category

	for _, categoryGroup := range *categoryGroups {
		if strings.Contains(categoryGroup.Name, "Obligatory Monthly Expenses") {
			for _, category := range categoryGroup.Categories {
				if !category.Hidden && !strings.Contains(category.Name, "Bank Fees") {
					monthlyExpensesCategories = append(monthlyExpensesCategories, category)
				}
			}
		}
	}

	return monthlyExpensesCategories
}
