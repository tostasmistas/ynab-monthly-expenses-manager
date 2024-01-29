package backend

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var categoryNames = [...]string{"Condominium", "Electricity", "TV / Internet / Phone", "Water"}

func TestIsValid(t *testing.T) {
	testCases := map[string]struct {
		monthlyExpenses  *MonthlyExpenses
		expectedValidity bool
	}{
		"invalid monthly expenses - no budget id": {
			monthlyExpenses: &MonthlyExpenses{
				AccountId: gofakeit.UUID(),
				Expenses:  createFakeExpenses(nil),
			},
			expectedValidity: false,
		},
		"invalid monthly expenses - no account id": {
			monthlyExpenses: &MonthlyExpenses{
				BudgetId: gofakeit.UUID(),
				Expenses: createFakeExpenses(nil),
			},
			expectedValidity: false,
		},
		"invalid monthly expenses - no expenses": {
			monthlyExpenses: &MonthlyExpenses{
				BudgetId:  gofakeit.UUID(),
				AccountId: gofakeit.UUID(),
			},
			expectedValidity: false,
		},
		"valid monthly expenses": {
			monthlyExpenses:  createFakeMonthlyExpenses(nil),
			expectedValidity: true,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, testCase.expectedValidity, testCase.monthlyExpenses.IsValid(),
				fmt.Sprintf("Expected monthly expenses validity to be %t", testCase.expectedValidity))
		})
	}
}

func TestSplitSharedMonthlyExpenses(t *testing.T) {
	sharedExpenseAmounts := map[string]float64{
		"Condominium":           245.75,
		"Electricity":           130.52,
		"TV / Internet / Phone": 85.90,
		"Water":                 60.25,
	}
	individualExpenseAmounts := map[string]float64{
		"Condominium":           0.00,
		"Electricity":           0.00,
		"TV / Internet / Phone": 0.00,
		"Water":                 0.00,
	}

	testCases := map[string]struct {
		roundUp                          bool
		expectedIndividualExpenseAmounts map[string]float64
	}{
		"round-up first": {
			roundUp: true,
			expectedIndividualExpenseAmounts: map[string]float64{
				"Condominium":           122.88,
				"Electricity":           65.26,
				"TV / Internet / Phone": 42.95,
				"Water":                 30.12,
			},
		},
		"round-down first": {
			roundUp: false,
			expectedIndividualExpenseAmounts: map[string]float64{
				"Condominium":           122.87,
				"Electricity":           65.26,
				"TV / Internet / Phone": 42.95,
				"Water":                 30.13,
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(rand.Float64, func() float64 {
				if testCase.roundUp {
					return 0.0
				} else {
					return 0.9
				}
			})
			defer patches.Reset()

			combinedMonthlyExpenses := CombinedMonthlyExpenses{
				SharedMonthlyExpenses:     createFakeMonthlyExpenses(sharedExpenseAmounts),
				IndividualMonthlyExpenses: createFakeMonthlyExpenses(individualExpenseAmounts),
			}
			combinedMonthlyExpenses.SplitSharedMonthlyExpenses()

			for _, categoryName := range categoryNames {
				sharedMonthlyExpense := combinedMonthlyExpenses.SharedMonthlyExpenses.Expenses[categoryName]
				individualMonthlyExpense := combinedMonthlyExpenses.IndividualMonthlyExpenses.Expenses[categoryName]

				expectedIndividualExpenseAmount := decimal.NewFromFloat(testCase.expectedIndividualExpenseAmounts[categoryName])
				actualIndividualExpenseAmount := individualMonthlyExpense.Amount

				assert.True(t, expectedIndividualExpenseAmount.Equal(actualIndividualExpenseAmount),
					fmt.Sprintf("Expected individual share for category '%s' to be %s (half of shared expense of %s), but got %s",
						categoryName, expectedIndividualExpenseAmount.String(), sharedMonthlyExpense.Amount.String(), actualIndividualExpenseAmount.String()))
			}
		})
	}
}

func createFakeMonthlyExpense(expenseAmount float64) *MonthlyExpense {
	monthlyExpense := &MonthlyExpense{}
	gofakeit.Struct(monthlyExpense)

	monthlyExpense.Amount = decimal.NewFromFloat(expenseAmount)

	return monthlyExpense
}

func createFakeExpenses(expenseAmounts map[string]float64) map[string]*MonthlyExpense {
	expenses := make(map[string]*MonthlyExpense, len(categoryNames))

	for _, categoryName := range categoryNames {
		monthlyExpense := createFakeMonthlyExpense(expenseAmounts[categoryName])
		expenses[categoryName] = monthlyExpense
	}

	return expenses
}

func createFakeMonthlyExpenses(expenseAmounts map[string]float64) *MonthlyExpenses {
	monthlyExpenses := MonthlyExpenses{}
	gofakeit.Struct(&monthlyExpenses)

	monthlyExpenses.Expenses = createFakeExpenses(expenseAmounts)

	return &monthlyExpenses
}
