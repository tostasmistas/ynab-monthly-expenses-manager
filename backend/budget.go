package backend

// BudgetSummary represents the summary of a budget
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

type DateFormat struct {
	Format string `json:"format"`
}

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

type Budgets struct {
	SharedBudget     BudgetSummary `json:"shared_budget"`
	IndividualBudget BudgetSummary `json:"individual_budget"`
}

const SharedBudgetName string = "Casa Reis-Pereira"
const IndividualBudgetName string = "Magui"

// GetBudgets fetches the list of budgets
// GET https://api.ynab.com/v1/budgets
func (client *APIClient) GetBudgets() ([]BudgetSummary, error) {
	budgetsResponse := struct {
		Data struct {
			Budgets       []BudgetSummary `json:"budgets"`
			DefaultBudget BudgetSummary   `json:"default_budget"`
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
