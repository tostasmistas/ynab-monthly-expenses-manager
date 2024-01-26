package backend

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/shopspring/decimal"
)

// MonthlyExpense represents a monthly expense with its YNAB category id, payee name, amount, and memo
type MonthlyExpense struct {
	CategoryId *string         `json:"category_id" mapstructure:"category_id"`
	PayeeName  *string         `json:"payee_name" mapstructure:"payee_name"`
	Amount     decimal.Decimal `json:"amount" mapstructure:"amount"`
	Memo       *string         `json:"memo" mapstructure:"memo"`
}

// MonthlyExpenses represents a collection of monthly expenses per category for a specific YNAB budget and account
type MonthlyExpenses struct {
	BudgetId  string                     `json:"budget_id" mapstructure:"budget_id"`
	AccountId string                     `json:"account_id" mapstructure:"account_id"`
	Expenses  map[string]*MonthlyExpense `json:"expenses" mapstructure:"expenses"`
}

// IsValid checks if a collection of monthly expenses is valid by ensuring a non-empty YNAB budget id and account id, and having exactly 4 monthly expenses
func (monthlyExpenses *MonthlyExpenses) IsValid() bool {
	return monthlyExpenses.AccountId != "" &&
		monthlyExpenses.BudgetId != "" &&
		len(monthlyExpenses.Expenses) == 4
}

// GetSharedMonthlyExpensePayeeName returns the predefined payee name for a given shared monthly expense category
func GetSharedMonthlyExpensePayeeName(categoryName string) string {
	var payeeName string

	switch categoryName {
	case "Condominium":
		payeeName = "Loja do Condomínio"
	case "Electricity":
		payeeName = "EDP"
	case "Water":
		payeeName = "EPAL"
	case "TV / Internet / Phone":
		payeeName = "Vodafone"
	}

	return payeeName
}

// GetSharedMonthlyExpenseMemo returns the predefined memo for a given shared monthly expense category
func GetSharedMonthlyExpenseMemo(categoryName string) string {
	switch categoryName {
	case "Condominium":
		return time.Now().AddDate(0, 1, 0).Format("January 2006")
	case "Electricity":
		return getSharedMonthlyExpenseBillingCycleMemo(11, 10)
	case "Water":
		return getSharedMonthlyExpenseBillingCycleMemo(4, 3)
	case "TV / Internet / Phone":
		return fmt.Sprintf("%s & %s",
			getSharedMonthlyExpenseBillingCycleMemo(9, 8),
			strings.Split(getSharedMonthlyExpenseBillingCycleMemo(16, 15), "- ")[1],
		)
	default:
		return ""
	}
}

// getSharedMonthlyExpenseBillingCycleMemo returns the predefined memo, based on the billing cycle, of a shared expense
func getSharedMonthlyExpenseBillingCycleMemo(billingCycleStart int, billingCycleEnd int) string {
	currentMonth := time.Now()
	pastMonth := time.Now().AddDate(0, -1, 0)

	return fmt.Sprintf("%s - %d %s to %d %s",
		pastMonth.Format("January 2006"),
		billingCycleStart, pastMonth.Format("January"),
		billingCycleEnd, currentMonth.Format("January"),
	)
}

// GetIndividualMonthlyExpensePayeeName returns the predefined payee name for an individual monthly expense
func GetIndividualMonthlyExpensePayeeName() string {
	return fmt.Sprintf("Transfer: %s", SharedBudgetName)
}

// GetIndividualMonthlyExpenseMemo returns the predefined memo for an individual monthly expense
func GetIndividualMonthlyExpenseMemo() string {
	return time.Now().Format("January 2006")
}

// SplitSharedMonthlyExpenses calculates the individual share for each monthly expense category
func SplitSharedMonthlyExpenses(sharedMonthlyExpenses MonthlyExpenses, individualMonthlyExpenses *MonthlyExpenses) {
	roundUp := rand.Float64() <= 0.4

	for categoryName, sharedMonthlyExpense := range sharedMonthlyExpenses.Expenses {
		splitExpenseAmount := sharedMonthlyExpense.Amount.Div(decimal.NewFromInt(2))
		roundedSplitExpenseAmount := sharedMonthlyExpense.Amount.DivRound(decimal.NewFromInt(2), 2)

		if !splitExpenseAmount.Equal(roundedSplitExpenseAmount) {
			if roundUp {
				splitExpenseAmount = splitExpenseAmount.RoundUp(2)
			} else {
				splitExpenseAmount = splitExpenseAmount.RoundDown(2)
			}

			roundUp = !roundUp
		}

		individualMonthlyExpense := individualMonthlyExpenses.Expenses[categoryName]
		individualMonthlyExpense.Amount = splitExpenseAmount
	}
}

// CreateIndividualMonthlyExpensesTransactions creates the YNAB transactions for the individual monthly expenses
func (monthlyExpenses *MonthlyExpenses) CreateIndividualMonthlyExpensesTransactions(client APIClient) bool {
	var subTransactions []SaveSubTransaction
	var sampleExpense MonthlyExpense
	var totalTransactionAmount int64

	for _, monthlyExpense := range monthlyExpenses.Expenses {
		sampleExpense = *monthlyExpense
		subTransactionAmount := monthlyExpense.Amount.Mul(decimal.NewFromInt(1000)).IntPart()
		totalTransactionAmount += subTransactionAmount

		subTransactions = append(subTransactions,
			SaveSubTransaction{
				Amount:     -subTransactionAmount,
				CategoryId: monthlyExpense.CategoryId,
			},
		)
	}

	transaction := SaveTransaction{
		AccountId:       to.StringPtr(monthlyExpenses.AccountId),
		Date:            time.Now().Format("2006-01-02"),
		Amount:          -totalTransactionAmount,
		PayeeName:       sampleExpense.PayeeName,
		CategoryId:      nil,
		Memo:            sampleExpense.Memo,
		Cleared:         "uncleared",
		Approved:        false,
		SubTransactions: subTransactions,
	}

	_, err := client.CreateTransaction(monthlyExpenses.BudgetId, transaction)

	return err == nil
}

// CreateSharedMonthlyExpensesTransactions creates the YNAB transactions for the shared monthly expenses
func (monthlyExpenses *MonthlyExpenses) CreateSharedMonthlyExpensesTransactions(client APIClient) bool {
	var transactions []SaveTransaction

	for _, monthlyExpense := range monthlyExpenses.Expenses {
		transactions = append(transactions,
			SaveTransaction{
				AccountId:  to.StringPtr(monthlyExpenses.AccountId),
				Date:       time.Now().Format("2006-01-02"),
				Amount:     -monthlyExpense.Amount.Mul(decimal.NewFromInt(1000)).IntPart(),
				PayeeName:  monthlyExpense.PayeeName,
				CategoryId: monthlyExpense.CategoryId,
				Memo:       monthlyExpense.Memo,
				Cleared:    "uncleared",
				Approved:   false,
			},
		)
	}

	_, err := client.CreateTransactions(monthlyExpenses.BudgetId, transactions)

	return err == nil
}
