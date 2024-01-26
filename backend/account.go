package backend

// Account represents a YNAB account
// This struct corresponds to the data structure defined in the YNAB API documentation
type Account struct {
	Id                  string           `json:"id"`
	Name                string           `json:"name"`
	Type                string           `json:"type"`
	OnBudget            bool             `json:"on_budget"`
	Closed              bool             `json:"closed"`
	Note                string           `json:"note"`
	Balance             int64            `json:"balance"`
	ClearedBalance      int64            `json:"cleared_balance"`
	UnclearedBalance    int64            `json:"uncleared_balance"`
	TransferPayeeId     string           `json:"transfer_payee_id"`
	DirectImportLinked  bool             `json:"direct_import_linked"`
	DirectImportInError bool             `json:"direct_import_in_error"`
	LastReconciledAt    string           `json:"last_reconciled_at"`
	DebtOriginalBalance int64            `json:"debt_original_balance"`
	DebtInterestRates   map[string]int64 `json:"debt_interest_rates"`
	DebtMinimumPayments map[string]int64 `json:"debt_minimum_payments"`
	DebtEscrowAmounts   map[string]int64 `json:"debt_escrow_amounts"`
	Deleted             bool             `json:"deleted"`
}

// Accounts represents a collection of YNAB accounts
type Accounts []Account

// SharedMonthlyExpensesAccountName is the predefined name of the YNAB account designated for shared monthly expenses
const SharedMonthlyExpensesAccountName string = "Millennium bcp"

// IndividualMonthlyExpensesAccountName is the predefined name of the YNAB account designated for individual monthly expenses
const IndividualMonthlyExpensesAccountName string = "CGD"

// GetMonthlyExpensesAccount fetches the YNAB account designated for monthly expenses based on its name
func (accounts *Accounts) GetMonthlyExpensesAccount(accountName string) Account {
	for _, account := range *accounts {
		if !account.Closed && !account.Deleted && account.Name == accountName {
			return account
		}
	}

	return Account{}
}
