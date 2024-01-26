<p align="center">
  <img src="build/appicon.png" width="80px" height="80px">
</p>

<h2 align="center">YNAB Monthly Expenses Manager</h2>

<br />

**YNAB Monthly Expenses Manager** is a personalized application crafted to simplify the process of managing and inputting monthly household expenses into [YNAB](https://www.ynab.com/), a zero-based budgeting tool.

<p align="center">
  <img width="700" alt="Screenshot 2024-01-24 at 22 24 10" src="https://github.com/tostasmistas/ynab-monthly-expenses-manager/assets/11311824/1ea183bb-cbbb-4d61-bdb2-0e96dc44f738">
</p>

## 💰 About

As a long-time YNAB user managing both shared and individual budgets, I grappled each month with the repetitive task of manually entering and splitting household expenses.
This challenge motivated me to create a personalized application that simplifies and automates these processes.

### 🚀 Features 

- **Expense input**: Effortless input of the total monthly household expenses by category.

- **Expense split**: Automatically calculates the individual share for each expense category.

- **YNAB integration**: Seamlessly integrates with YNAB, automatically inputting the shared and individual expense transactions.

<br />

> [!NOTE]  
> This application is customized based on my YNAB budget data and is primarily intended for personal use.
While it may not be suitable for general use, feel free to explore and use it as a reference for building your own personalized solution.
If you have questions or suggestions for improvement, let's connect and make managing YNAB budgets even more efficient! 🤝

## 📖 How to use

1. **Expense input**

Input the total monthly expense for each category - `Condominium`, `Electricity`, `Water`, and `TV / Internet / Phone` - under the card named `Total Monthly Expenses`.

2. **Expense split**

After entering the total monthly expenses for all categories, clicking the `Split` button triggers the application to calculate and display, under the card named `Individual share`, a breakdown of the individual share for each expense category.

The application utilizes a rounding algorithm to fairly distribute shared monthly expenses among individuals:

- The algorithm processes each expense category within the shared monthly expenses, calculating the individual share by dividing each expense amount by 2 in consideration of a two-person household.
- In cases where the resulting individual share is not an exact division by 2, requiring rounding to adhere to the 2-decimal place constraint inherent in monetary values, the algorithm employs a balanced rounding approach:
  - Initially, variability is introduced by assigning a 50% chance of rounding up the expense. Then, subsequent rounding conditions alternate based on the previous rounding result, i.e., if the current expense is rounded up, the subsequent one will be rounded down if needed, and vice versa.
  - This variability in rounding is designed for fairness. Even if only one expense requires rounding each month, the algorithm introduces a 50% chance of rounding up initially, ensuring an equitable distribution of rounding over time, providing both individuals with an equal chance of experiencing rounded-up or rounded-down amounts.

<details>
<summary>Expense split example</summary>

<br />

Suppose the total shared monthly expenses are:

- Condominium: `245.75€`
- Electricity: `130.52€`
- Water: `60.25€`
- TV / Internet / Phone: `85.90€`

Applying the algorithm step by step:

- Condominium:
  - The individual share is `245.75€ / 2 = 122.875€`, which is not an exact division by 2, requiring rounding. The initial 50% chance of rounding up is considered. 
  - Assuming that chance dictated a rounding down, the individual share is rounded down to `122.87€`.

- Electricity:
  - The individual share is `130.52€ / 2 = 65.26€`, which is an exact division by 2, so no rounding is required.

- Water:
  - The individual share is `60.25€ / 2 = 30.125€`, which is not an exact division by 2, requiring rounding.
  - The rounding is determined by the toggle from the Condominium category, so the individual share is rounded up to `30.13€`.

- TV / Internet / Phone:
  - The individual share is `85.90€ / 2 = 42.95€`, which is an exact division by 2, so no rounding is required.

</details>

3. **YNAB integration**

After inputting and splitting the monthly household expenses, the final step is seamless integration with YNAB, initiated by clicking the `Import` button.

For the shared expenses, under the shared budget in YNAB and for the shared monthly expenses account, distinct transactions are created for each category. These transactions detail the expense amount, the payee, and the billing cycle in the memo field.

For the individual share, under the individual budget in YNAB and for the individual monthly expenses account, a main transaction is created encompassing the total individual share amount. Sub-transactions are nested within, capturing each individual's share for every expense category.

> [!WARNING]  
> Without a valid YNAB Personal Access Token, the application won't load properly. This token should be configured in the constant `AccessToken` located in the `backend/api_client.go` file.

## 🧑‍💻 Development mode

This application is built using [Wails](https://wails.io/) and uses Go on the backend and React on the frontend.

#### Requirements

- [Wails v2](https://github.com/wailsapp/wails)
- [Go v1.21+](https://go.dev/doc/install)
- [Node.js v15+](https://nodejs.org/en/download/)

To develop the application locally clone the repository and in the root directory run the command `wails dev` and in the frontend directory run the command `npm run dev`.
