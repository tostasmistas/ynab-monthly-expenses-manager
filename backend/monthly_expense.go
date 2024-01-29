package backend

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// MonthlyExpense represents a monthly expense with its YNAB category id, payee name, amount, and memo
type MonthlyExpense struct {
	CategoryId *string         `json:"category_id" mapstructure:"category_id" fake:"{uuid}"`
	PayeeName  *string         `json:"payee_name" mapstructure:"payee_name" fake:"{company}"`
	Amount     decimal.Decimal `json:"amount" mapstructure:"amount" fake:"skip"`
	Memo       *string         `json:"memo" mapstructure:"memo" fake:"{sentence}"`
}

// MonthlyExpenses represents a collection of monthly expenses per category for a specific YNAB budget and account
type MonthlyExpenses struct {
	BudgetId  string                     `json:"budget_id" mapstructure:"budget_id" fake:"{uuid}"`
	AccountId string                     `json:"account_id" mapstructure:"account_id" fake:"{uuid}"`
	Expenses  map[string]*MonthlyExpense `json:"expenses" mapstructure:"expenses" fake:"skip"`
}

// CombinedMonthlyExpenses represents a collection of monthly expenses, combining both the shared and individual monthly expenses
type CombinedMonthlyExpenses struct {
	SharedMonthlyExpenses     *MonthlyExpenses `json:"shared_monthly_expenses"`
	IndividualMonthlyExpenses *MonthlyExpenses `json:"individual_monthly_expenses"`
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
		currentMonth := time.Now()
		nextMonth := time.Date(currentMonth.Year(), currentMonth.Month()+1, 1, 0, 0, 0, 0, currentMonth.Location())
		return nextMonth.Format("January 2006")
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
	pastMonth := time.Date(currentMonth.Year(), currentMonth.Month()-1, 1, 0, 0, 0, 0, currentMonth.Location())

	return fmt.Sprintf("%s - %d %s to %d %s",
		pastMonth.Format("January 2006"),
		billingCycleStart, pastMonth.Format("January"),
		billingCycleEnd, currentMonth.Format("January"),
	)
}

// GetIndividualMonthlyExpensePayeeName returns the predefined payee name for an individual monthly expense
func GetIndividualMonthlyExpensePayeeName(payeeName string) string {
	return fmt.Sprintf("Transfer: %s", payeeName)
}

// GetIndividualMonthlyExpenseMemo returns the predefined memo for an individual monthly expense
func GetIndividualMonthlyExpenseMemo() string {
	return fmt.Sprintf("%s - Household Expenses", time.Now().Format("January 2006"))
}

// SplitSharedMonthlyExpenses calculates the individual share for each monthly expense category
func (combinedMonthlyExpenses *CombinedMonthlyExpenses) SplitSharedMonthlyExpenses() {
	sharedMonthlyExpenses := combinedMonthlyExpenses.SharedMonthlyExpenses
	individualMonthlyExpenses := combinedMonthlyExpenses.IndividualMonthlyExpenses

	roundUp := rand.Float64() <= 0.4

	categoryNames := maps.Keys(sharedMonthlyExpenses.Expenses)
	slices.Sort(categoryNames)

	for _, categoryName := range categoryNames {
		sharedMonthlyExpense := sharedMonthlyExpenses.Expenses[categoryName]
		individualMonthlyExpense := individualMonthlyExpenses.Expenses[categoryName]

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

		individualMonthlyExpense.Amount = splitExpenseAmount
	}
}

// CreateSharedMonthlyExpensesTransactions creates the YNAB transactions for the shared monthly expenses
func (combinedMonthlyExpenses *CombinedMonthlyExpenses) CreateSharedMonthlyExpensesTransactions(client APIClient) bool {
	sharedMonthlyExpenses := combinedMonthlyExpenses.SharedMonthlyExpenses
	individualMonthlyExpenses := combinedMonthlyExpenses.IndividualMonthlyExpenses

	var transactions []SaveTransaction

	var subTransactionsForMyIndividualShare []SaveSubTransaction
	var totalMyIndividualShareAmount decimal.Decimal

	var subTransactionsForOtherIndividualShare []SaveSubTransaction
	var totalOtherIndividualShareAmount decimal.Decimal

	for categoryName, monthlyExpense := range sharedMonthlyExpenses.Expenses {
		transactionAmount := monthlyExpense.Amount

		transactions = append(transactions,
			createTransaction(
				sharedMonthlyExpenses.AccountId,
				transactionAmount.Neg(),
				monthlyExpense.PayeeName,
				monthlyExpense.CategoryId,
				monthlyExpense.Memo,
				nil,
			),
		)

		myIndividualShareAmount := individualMonthlyExpenses.Expenses[categoryName].Amount
		totalMyIndividualShareAmount = totalMyIndividualShareAmount.Add(myIndividualShareAmount)

		otherIndividualShareAmount := transactionAmount.Sub(myIndividualShareAmount)
		totalOtherIndividualShareAmount = totalOtherIndividualShareAmount.Add(otherIndividualShareAmount)

		subTransactionsForMyIndividualShare = append(subTransactionsForMyIndividualShare,
			createSubTransaction(
				myIndividualShareAmount,
				monthlyExpense.CategoryId,
			),
		)

		subTransactionsForOtherIndividualShare = append(subTransactionsForOtherIndividualShare,
			createSubTransaction(
				otherIndividualShareAmount,
				monthlyExpense.CategoryId,
			),
		)
	}

	transactions = append(transactions,
		createTransaction(
			sharedMonthlyExpenses.AccountId,
			totalMyIndividualShareAmount,
			to.StringPtr(GetIndividualMonthlyExpensePayeeName("Magui")),
			nil,
			to.StringPtr(GetIndividualMonthlyExpenseMemo()),
			subTransactionsForMyIndividualShare,
		),
		createTransaction(
			sharedMonthlyExpenses.AccountId,
			totalOtherIndividualShareAmount,
			to.StringPtr(GetIndividualMonthlyExpensePayeeName("Jão")),
			nil,
			to.StringPtr(GetIndividualMonthlyExpenseMemo()),
			subTransactionsForOtherIndividualShare,
		),
	)

	_, err := client.CreateTransactions(sharedMonthlyExpenses.BudgetId, transactions)

	return err == nil
}

// CreateIndividualMonthlyExpensesTransactions creates the YNAB transactions for the individual monthly expenses
func (combinedMonthlyExpenses *CombinedMonthlyExpenses) CreateIndividualMonthlyExpensesTransactions(client APIClient) bool {
	individualMonthlyExpenses := combinedMonthlyExpenses.IndividualMonthlyExpenses

	var sampleExpense MonthlyExpense

	var subTransactions []SaveSubTransaction
	var totalTransactionAmount decimal.Decimal

	for _, monthlyExpense := range individualMonthlyExpenses.Expenses {
		sampleExpense = *monthlyExpense

		subTransactionAmount := monthlyExpense.Amount
		totalTransactionAmount = totalTransactionAmount.Add(subTransactionAmount)

		subTransactions = append(subTransactions,
			createSubTransaction(
				subTransactionAmount.Neg(),
				monthlyExpense.CategoryId,
			),
		)
	}

	transaction := createTransaction(
		individualMonthlyExpenses.AccountId,
		totalTransactionAmount.Neg(),
		sampleExpense.PayeeName,
		nil,
		sampleExpense.Memo,
		subTransactions,
	)

	_, err := client.CreateTransaction(individualMonthlyExpenses.BudgetId, transaction)

	return err == nil
}

// createTransaction creates a new SaveTransaction instance
func createTransaction(accountId string, amount decimal.Decimal, payeeName *string, categoryId *string, memo *string, subTransactions []SaveSubTransaction) SaveTransaction {
	return SaveTransaction{
		AccountId:       to.StringPtr(accountId),
		Date:            time.Now().Format("2006-01-02"),
		Amount:          amount.Mul(decimal.NewFromInt(1000)).IntPart(),
		PayeeName:       payeeName,
		CategoryId:      categoryId,
		Memo:            memo,
		Cleared:         "uncleared",
		Approved:        false,
		SubTransactions: subTransactions,
	}
}

// createSubTransaction creates a new SaveSubTransaction instance
func createSubTransaction(amount decimal.Decimal, categoryId *string) SaveSubTransaction {
	return SaveSubTransaction{
		Amount:     amount.Mul(decimal.NewFromInt(1000)).IntPart(),
		CategoryId: categoryId,
	}
}
