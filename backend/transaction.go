package backend

import (
	"fmt"
)

// TransactionSummary represents the summary of a YNAB transaction
// This struct corresponds to the data structure defined in the YNAB API documentation
type TransactionSummary struct {
	Id                      string `json:"id"`
	Date                    string `json:"date"`
	Amount                  int64  `json:"amount"`
	Memo                    string `json:"memo"`
	Cleared                 string `json:"cleared"`
	Approved                bool   `json:"approved"`
	FlagColor               string `json:"flag_color"`
	AccountId               string `json:"account_id"`
	PayeeId                 string `json:"payee_id"`
	CategoryId              string `json:"category_id"`
	TransferAccountId       string `json:"transfer_account_id"`
	TransferTransactionId   string `json:"transfer_transaction_id"`
	MatchedTransactionId    string `json:"matched_transaction_id"`
	ImportId                string `json:"import_id"`
	ImportPayeeName         string `json:"import_payee_name"`
	ImportPayeeNameOriginal string `json:"import_payee_name_original"`
	DebtTransactionType     string `json:"debt_transaction_type"`
	Deleted                 bool   `json:"deleted"`
}

// TransactionDetail represents the details of a YNAB transaction
// This struct corresponds to the data structure defined in the YNAB API documentation
type TransactionDetail struct {
	TransactionSummary
	AccountName     string           `json:"account_name"`
	PayeeName       string           `json:"payee_name"`
	CategoryName    string           `json:"category_name"`
	SubTransactions []SubTransaction `json:"subtransactions"`
}

// SubTransaction represents a sub-transaction of a YNAB transaction
// This struct corresponds to the data structure defined in the YNAB API documentation
type SubTransaction struct {
	Id                    string `json:"id"`
	TransactionId         string `json:"transaction_id"`
	Amount                int64  `json:"amount"`
	Memo                  string `json:"memo"`
	PayeeId               string `json:"payee_id"`
	PayeeName             string `json:"payee_name"`
	CategoryId            string `json:"category_id"`
	CategoryName          string `json:"category_name"`
	TransferAccountId     string `json:"transfer_account_id"`
	TransferTransactionId string `json:"transfer_transaction_id"`
	Deleted               bool   `json:"deleted"`
}

// SaveTransaction represents the schema for creating a new YNAB transaction
// This struct corresponds to the data structure defined in the YNAB API documentation
type SaveTransaction struct {
	AccountId       *string              `json:"account_id"`
	Date            string               `json:"date"`
	Amount          int64                `json:"amount"`
	PayeeId         *string              `json:"payee_id"`
	PayeeName       *string              `json:"payee_name"`
	CategoryId      *string              `json:"category_id"`
	Memo            *string              `json:"memo"`
	Cleared         string               `json:"cleared"`
	Approved        bool                 `json:"approved"`
	FlagColor       *string              `json:"flag_color"`
	ImportId        *string              `json:"import_id"`
	SubTransactions []SaveSubTransaction `json:"subtransactions"`
}

// SaveSubTransaction represents the schema for creating a new sub-transaction of a YNAB transaction
// This struct corresponds to the data structure defined in the YNAB API documentation
type SaveSubTransaction struct {
	Amount     int64   `json:"amount"`
	PayeeId    *string `json:"payee_id"`
	PayeeName  *string `json:"payee_name"`
	CategoryId *string `json:"category_id"`
	Memo       *string `json:"memo"`
}

// CreateTransaction creates a new YNAB transaction for a YNAB budget
// POST https://api.ynab.com/v1/budgets/{budget_id}/transactions
func (client *APIClient) CreateTransaction(budgetId string, transaction SaveTransaction) (TransactionDetail, error) {
	transactionBody := struct {
		Transaction SaveTransaction `json:"transaction"`
	}{
		Transaction: transaction,
	}

	transactionResponse := struct {
		Data struct {
			TransactionIds []string          `json:"transaction_ids"`
			Transaction    TransactionDetail `json:"transaction"`
		} `json:"data"`
	}{}

	response, err := client.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(transactionBody).
		SetResult(&transactionResponse).
		Post(fmt.Sprintf("budgets/%s/transactions", budgetId))

	if err = client.ValidateResponse(response, err); err != nil {
		return TransactionDetail{}, err
	}

	return transactionResponse.Data.Transaction, nil
}

// CreateTransactions creates new YNAB transactions for a YNAB budget
// POST https://api.ynab.com/v1/budgets/{budget_id}/transactions
func (client *APIClient) CreateTransactions(budgetId string, transactions []SaveTransaction) ([]TransactionDetail, error) {
	transactionsBody := struct {
		Transactions []SaveTransaction `json:"transactions"`
	}{
		Transactions: transactions,
	}

	transactionsResponse := struct {
		Data struct {
			TransactionIds []string            `json:"transaction_ids"`
			Transactions   []TransactionDetail `json:"transactions"`
		} `json:"data"`
	}{}

	response, err := client.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(transactionsBody).
		SetResult(&transactionsResponse).
		Post(fmt.Sprintf("budgets/%s/transactions", budgetId))

	if err = client.ValidateResponse(response, err); err != nil {
		return nil, err
	}

	return transactionsResponse.Data.Transactions, nil
}
