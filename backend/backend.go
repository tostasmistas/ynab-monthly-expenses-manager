package backend

import (
	"context"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/forPelevin/gomoji"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Backend struct {
	Context                   context.Context
	APIClient                 *APIClient
	SharedMonthlyExpenses     *MonthlyExpenses
	IndividualMonthlyExpenses *MonthlyExpenses
}

func SetupBackend() *Backend {
	var apiClient APIClient
	apiClient.Client = resty.New()
	apiClient.Configure()

	budgets, _ := apiClient.GetBudgets()
	var sharedBudget, individualBudget BudgetSummary

	for _, budget := range budgets {
		if strings.Contains(budget.Name, SharedBudgetName) {
			sharedBudget = budget
		} else if strings.Contains(budget.Name, IndividualBudgetName) {
			individualBudget = budget
		}
	}

	sharedAccount := sharedBudget.Accounts.GetMonthlyExpensesAccount()
	sharedCategories, _ := apiClient.GetCategories(sharedBudget.Id)

	individualAccount := individualBudget.Accounts.GetMonthlyExpensesAccount()
	individualCategories, _ := apiClient.GetCategories(individualBudget.Id)

	sharedMonthlyExpenses := MonthlyExpenses{
		BudgetId:  sharedBudget.Id,
		AccountId: sharedAccount.Id,
		Expenses:  make(map[string]*MonthlyExpense),
	}

	individualMonthlyExpenses := MonthlyExpenses{
		BudgetId:  individualBudget.Id,
		AccountId: individualAccount.Id,
		Expenses:  make(map[string]*MonthlyExpense),
	}

	for _, category := range sharedCategories.GetMonthlyExpensesCategories() {
		categoryName := gomoji.RemoveEmojis(category.Name)

		sharedMonthlyExpenses.Expenses[categoryName] = &MonthlyExpense{
			CategoryId: to.StringPtr(category.Id),
			PayeeName:  to.StringPtr(GetSharedMonthlyExpensePayeeName(categoryName)),
			Memo:       to.StringPtr(GetSharedMonthlyExpenseMemo(categoryName)),
		}
	}

	for _, category := range individualCategories.GetMonthlyExpensesCategories() {
		categoryName := gomoji.RemoveEmojis(category.Name)

		individualMonthlyExpenses.Expenses[categoryName] = &MonthlyExpense{
			CategoryId: to.StringPtr(category.Id),
			PayeeName:  to.StringPtr(GetIndividualMonthlyExpensePayeeName()),
			Memo:       to.StringPtr(GetIndividualMonthlyExpenseMemo()),
		}
	}

	return &Backend{
		APIClient:                 &apiClient,
		SharedMonthlyExpenses:     &sharedMonthlyExpenses,
		IndividualMonthlyExpenses: &individualMonthlyExpenses,
	}
}

func (backend *Backend) Startup(context context.Context) {
	backend.Context = context

	runtime.EventsOn(context, "sharedMonthlyExpensesInput", func(args ...interface{}) {
		var sharedMonthlyExpenses MonthlyExpenses

		decoderConfig := &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				StringToDecimalHookFunc(),
			),
			Result: &sharedMonthlyExpenses,
		}
		decoder, _ := mapstructure.NewDecoder(decoderConfig)
		decoder.Decode(args[0])

		SplitSharedMonthlyExpenses(sharedMonthlyExpenses, backend.IndividualMonthlyExpenses)

		runtime.EventsEmit(context, "sharedMonthlyExpensesSplit", backend.IndividualMonthlyExpenses)
	})
}

func (backend *Backend) DomReady(context context.Context) {
	runtime.EventsEmit(context, "backendSetupComplete",
		backend.SharedMonthlyExpenses.IsValid() && backend.IndividualMonthlyExpenses.IsValid(),
	)
}

func (backend *Backend) GetSharedMonthlyExpenses() *MonthlyExpenses {
	return backend.SharedMonthlyExpenses
}

func (backend *Backend) CreateMonthlyExpensesTransactions(sharedMonthlyExpenses *MonthlyExpenses, individualMonthlyExpenses *MonthlyExpenses) bool {
	return sharedMonthlyExpenses.CreateSharedMonthlyExpensesTransactions(*backend.APIClient) &&
		individualMonthlyExpenses.CreateIndividualMonthlyExpensesTransactions(*backend.APIClient)
}
