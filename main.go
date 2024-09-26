package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/parz3val/krisp-account-manager/twilio_utils"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func monospacedCenteredLabel(text binding.String) *widget.Label {
	return widget.NewLabelWithData(text)
}

func main() {
	a := app.New()
	w := a.NewWindow("KRISP ACCOUNT MANAGER")
	config := twilio_utils.Settings{
		AccountSid: "AccountSid",
		AuthToken:  "AuthToken",
	}
	var cachedAccounts []api.ApiV2010Account

	fmt.Println(config)
	msg_binding_string := binding.NewString()
	msg_binding_string.Set("NOTIFICATIONS")
	// msg_string, _ := msg_binding_string.Get()
	msg_label := monospacedCenteredLabel(msg_binding_string)

	suspended_button := widget.NewButton("Suspended Subaccounts with Number", func() {
		// check if the config is empty or not
		if config.AccountSid == "" || config.AuthToken == "" {
			msg_binding_string.Set("Please enter Account Sid and Auth Token")
			return
		}
		accounts, err := twilio_utils.FetchAccounts(config, "suspended")

		msg_binding_string.Set("Fetching Accounts...")

		if err == nil {
			cachedAccounts = accounts
			filename := twilio_utils.GenerateFilename("suspended with number")
			err := twilio_utils.WriteAccountsToCSV(cachedAccounts, filename)
			if err == nil {
				msg := fmt.Sprintf("Accounts saved to %s", filename)
				msg_binding_string.Set(msg)
			}
		}
	})
	active_button := widget.NewButton("Active subaccounts having numbers", func() {
		// check if the config is empty or not
		if config.AccountSid == "" || config.AuthToken == "" {
			msg_binding_string.Set("Please enter Account Sid and Auth Token")
			return
		}

		msg_binding_string.Set("Fetching Accounts...")

		accounts, err := twilio_utils.FetchAccounts(config, "active")
		if err != nil {
			msg_binding_string.Set("Error fetching accounts")
			return
		}
		msg_binding_string.Set("Accounts Fetched!")
		msg_binding_string.Set("Fetching Numbers for sub accounts!")
		filename := twilio_utils.GenerateFilename("active with number")
		twilio_utils.AccountsProcessor(accounts, filename, false)
		msg_binding_string.Set("Accounts saved to " + filename)

	})
	active_without_button := widget.NewButton("Active subaccounts with no numbers", func() {
		// check if the config is empty or not
		if config.AccountSid == "" || config.AuthToken == "" {
			msg_binding_string.Set("Please enter Account Sid and Auth Token")
			return
		}
		accounts, err := twilio_utils.FetchAccounts(config, "active")
		if err == nil {
			cachedAccounts = accounts
			msg_binding_string.Set("Accounts Fetched!")
			msg_binding_string.Set("Checking the accounts to see if they have numbers")
			filename := twilio_utils.GenerateFilename("active without number")
			twilio_utils.AccountsProcessor(cachedAccounts, filename, true)
		}
	})

	account_sid_input := widget.NewPasswordEntry()
	account_sid_input.SetPlaceHolder("Enter Account Sid")
	auth_token_input := widget.NewPasswordEntry()
	auth_token_input.SetPlaceHolder("Enter Account Auth Token")

	save_button := widget.NewButton("Save", func() {
		// check if the input text is empty or not
		if account_sid_input.Text == "" || auth_token_input.Text == "" {
			msg_binding_string.Set("Please enter Account Sid and Auth Token")
			return
		}
		config.AccountSid = account_sid_input.Text
		config.AuthToken = auth_token_input.Text
		msg_binding_string.Set("Account SID and AUTH Token Saved")
	})

	msg_grid := container.New(layout.NewCenterLayout(), msg_label)
	input_grid := container.New(layout.NewGridLayout(1), account_sid_input, auth_token_input, save_button)
	button_grid := container.New(layout.NewGridLayout(1), suspended_button, active_button, active_without_button)

	switchable_tabs := container.NewAppTabs(
		container.NewTabItem("Account Settings", input_grid),
		container.NewTabItem("Dashboard", button_grid),
	)
	page_grid := container.New(layout.NewGridLayout(1), switchable_tabs, msg_grid)
	pageContent := container.New(layout.NewStackLayout(), page_grid)

	w.SetContent(pageContent)
	w.ShowAndRun()
}
