package backend

import (
	"context"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/forPelevin/gomoji"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Backend encapsulates the YNAB API client and the shared and individual monthly expenses
type Backend struct {
	Context                 context.Context
	APIClient               *APIClient
	CombinedMonthlyExpenses *CombinedMonthlyExpenses
}

// SetupBackend creates a new Backend instance
func SetupBackend() *Backend {
	var apiClient APIClient
	apiClient.Client = resty.New()
	apiClient.Configure()

	budgets, _ := apiClient.GetBudgets()

	sharedBudget := budgets.GetBudget(SharedBudgetName)
	sharedMonthlyExpensesAccount := sharedBudget.Accounts.GetMonthlyExpensesAccount(SharedMonthlyExpensesAccountName)
	sharedCategories, _ := apiClient.GetCategories(sharedBudget.Id)

	individualBudget := budgets.GetBudget(IndividualBudgetName)
	individualMonthlyExpensesAccount := individualBudget.Accounts.GetMonthlyExpensesAccount(IndividualMonthlyExpensesAccountName)
	individualCategories, _ := apiClient.GetCategories(individualBudget.Id)

	sharedMonthlyExpenses := MonthlyExpenses{
		BudgetId:  sharedBudget.Id,
		AccountId: sharedMonthlyExpensesAccount.Id,
		Expenses:  make(map[string]*MonthlyExpense),
	}

	individualMonthlyExpenses := MonthlyExpenses{
		BudgetId:  individualBudget.Id,
		AccountId: individualMonthlyExpensesAccount.Id,
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
			PayeeName:  to.StringPtr(GetIndividualMonthlyExpensePayeeName(SharedBudgetName)),
			Memo:       to.StringPtr(GetIndividualMonthlyExpenseMemo()),
		}
	}

	return &Backend{
		APIClient: &apiClient,
		CombinedMonthlyExpenses: &CombinedMonthlyExpenses{
			SharedMonthlyExpenses:     &sharedMonthlyExpenses,
			IndividualMonthlyExpenses: &individualMonthlyExpenses,
		},
	}
}

// Startup sets the backend context and registers an event handler to listen for the "sharedMonthlyExpensesInput" event
// When this event occurs the individual share for each monthly expense category is calculated and then the "sharedMonthlyExpensesSplit" event is emitted
func (backend *Backend) Startup(context context.Context) {
	backend.Context = context

	runtime.EventsOn(context, "sharedMonthlyExpensesInput", func(args ...interface{}) {
		decoderConfig := &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				StringToDecimalHookFunc(),
			),
			Result: &backend.CombinedMonthlyExpenses.SharedMonthlyExpenses,
		}
		decoder, _ := mapstructure.NewDecoder(decoderConfig)
		decoder.Decode(args[0])

		backend.CombinedMonthlyExpenses.SplitSharedMonthlyExpenses()

		runtime.EventsEmit(context, "sharedMonthlyExpensesSplit", backend.CombinedMonthlyExpenses.IndividualMonthlyExpenses)
	})
}

// DomReady emits the "backendSetupComplete" event indicating if both the shared and individual monthly expenses are valid as that is a requirement for the application
func (backend *Backend) DomReady(context context.Context) {
	runtime.EventsEmit(context, "backendSetupComplete",
		backend.CombinedMonthlyExpenses.SharedMonthlyExpenses.IsValid() &&
			backend.CombinedMonthlyExpenses.IndividualMonthlyExpenses.IsValid(),
	)
}

// GetSharedMonthlyExpenses returns the shared monthly expenses
func (backend *Backend) GetSharedMonthlyExpenses() *MonthlyExpenses {
	return backend.CombinedMonthlyExpenses.SharedMonthlyExpenses
}

// CreateMonthlyExpensesTransactions creates YNAB transactions for the shared and individual monthly expenses
func (backend *Backend) CreateMonthlyExpensesTransactions(combinedMonthlyExpenses *CombinedMonthlyExpenses) bool {
	return combinedMonthlyExpenses.CreateSharedMonthlyExpensesTransactions(*backend.APIClient) &&
		combinedMonthlyExpenses.CreateIndividualMonthlyExpensesTransactions(*backend.APIClient)
}
