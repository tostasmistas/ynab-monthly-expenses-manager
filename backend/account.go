package backend

// Account represents an account
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

type Accounts []Account

const SharedMonthlyExpensesAccountName string = "Millennium bcp"
const IndividualMonthlyExpensesAccountName string = "CGD"

func (accounts *Accounts) GetMonthlyExpensesAccount() Account {
	for _, account := range *accounts {
		if !account.Closed && !account.Deleted &&
			(account.Name == SharedMonthlyExpensesAccountName ||
				account.Name == IndividualMonthlyExpensesAccountName) {
			return account
		}
	}

	return Account{}
}
