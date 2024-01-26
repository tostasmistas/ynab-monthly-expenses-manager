package backend

import (
	"github.com/forPelevin/gomoji"
)

// BudgetSummary represents the summary of a YNAB budget
// This struct corresponds to the data structure defined in the YNAB API documentation
type BudgetSummary struct {
	Id             string         `json:"id"`
	Name           string         `json:"name"`
	LastModifiedOn string         `json:"last_modified_on"`
	FirstMonth     string         `json:"first_month"`
	LastMonth      string         `json:"last_month"`
	DateFormat     DateFormat     `json:"date_format"`
	CurrencyFormat CurrencyFormat `json:"currency_format"`
	Accounts       Accounts       `json:"accounts"`
}

// DateFormat represents the date format setting for a YNAB budget
// This struct corresponds to the data structure defined in the YNAB API documentation
type DateFormat struct {
	Format string `json:"format"`
}

// CurrencyFormat represents the currency format setting for a YNAB budget
// This struct corresponds to the data structure defined in the YNAB API documentation
type CurrencyFormat struct {
	IsoCode          string `json:"iso_code"`
	ExampleFormat    string `json:"example_format"`
	DecimalDigits    int32  `json:"decimal_digits"`
	DecimalSeparator string `json:"decimal_separator"`
	SymbolFirst      bool   `json:"symbol_first"`
	GroupSeparator   string `json:"group_separator"`
	CurrencySymbol   string `json:"currency_symbol"`
	DisplaySymbol    bool   `json:"display_symbol"`
}

// Budgets represents a collection of YNAB budgets
type Budgets []BudgetSummary

// SharedBudgetName is the predefined name of the shared YNAB budget
const SharedBudgetName string = "Casa Reis-Pereira"

// IndividualBudgetName is the predefined name of the individual YNAB budget
const IndividualBudgetName string = "Magui"

// GetBudgets fetches the list of YNAB budgets
// GET https://api.ynab.com/v1/budgets
func (client *APIClient) GetBudgets() (Budgets, error) {
	budgetsResponse := struct {
		Data struct {
			Budgets       Budgets       `json:"budgets"`
			DefaultBudget BudgetSummary `json:"default_budget"`
		} `json:"data"`
	}{}

	response, err := client.Client.R().
		SetQueryParams(map[string]string{
			"include_accounts": "true",
		}).
		SetResult(&budgetsResponse).
		Get("budgets")

	if err = client.ValidateResponse(response, err); err != nil {
		return nil, err
	}

	return budgetsResponse.Data.Budgets, nil
}

// GetBudget fetches a YNAB budget based on its name
func (budgets *Budgets) GetBudget(budgetName string) BudgetSummary {
	for _, budget := range *budgets {
		if gomoji.RemoveEmojis(budget.Name) == budgetName {
			return budget
		}
	}

	return BudgetSummary{}
}
