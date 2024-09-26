package twilio_utils

import (
	// "github.com/kevinburke/twilio-go"
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func FetchAccounts(settings Settings, status string) ([]api.ApiV2010Account, error) {
	fmt.Printf("AccountSid: %v\n", settings.AccountSid)
	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: settings.AccountSid,
			Password: settings.AuthToken,
		},
	)
	params := &api.ListAccountParams{}
	params.SetStatus(status)
	// params.SetPageSize(1000)
	resp, err := client.Api.ListAccount(params)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func WriteAccountsToCSV(accounts []api.ApiV2010Account, filename string) error {
	// open a file for writing
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	header := []string{"AuthToken", "DateCreated", "DateUpdated", "FriendlyName", "Sid", "Status"}
	writer := csv.NewWriter(f)
	defer writer.Flush()
	err = writer.Write(header)

	if err != nil {
		return err
	}
	for record := range accounts {
		recordLine := []string{
			*accounts[record].AuthToken,
			*accounts[record].DateCreated,
			*accounts[record].DateUpdated,
			*accounts[record].FriendlyName,
			*accounts[record].Sid,
			*accounts[record].Status,
		}
		err = writer.Write(recordLine)
		if err != nil {
			fmt.Printf("Error writing record to file: %v\n", err)
		}

	}
	return nil

}

type Accountant struct {
	AccountSid   string
	AuthToken    string
	FriendlyName string
	DateCreated  string
	DateUpdated  string
	numbers      []string
	numeros      []Numeros
}

type Numeros struct {
	PhoneNumber string
	DateCreated string
	DateUpdated string
	AccountName string
	AccountSid  string
	Status      string
	VoiceUrl    string
	SmsUrl      string
}

func AccountsProcessor(accounts []api.ApiV2010Account, filename string, check_nums bool) {
	const batchSize = 50
	// Create wait-group for the current batch
	var wg sync.WaitGroup
	var chWait sync.WaitGroup

	// create a channel of Accountant
	accountant := make(chan Accountant)

	for i := 0; i < len(accounts); i += batchSize {
		// increase the wait-group counter
		wg.Add(1)
		fmt.Printf("Adding to the wait group, i: %v\n, i + batchSize : %v\n", i, i+batchSize)
		go fetchNumbersBatch(accounts[i:min(i+batchSize, len(accounts))], accountant, &wg, &chWait, check_nums)
		wg.Wait()
	}

	// close channels when all the channels are doing writing
	go func() {
		chWait.Wait()
		close(accountant)
	}()

	// create a slice of Accountant
	accountants := make([]Accountant, 0)
	// create slice of Numeros
	// numeros := make([]Numeros, 0)

	// read from accountants channel into the accountants slice
	for accountant_ := range accountant {
		accountants = append(accountants, accountant_)
		// if accountant_.numeros != nil {
		// 	numeros = append(numeros, accountant_.numeros...)
		// }
	}

	// save the accountants to csv
	err := saveAccountantsToCSV(accountants, filename)
	if err != nil {
		fmt.Printf("Error saving accountants to csv: %v\n", err)
	}
	// save numeros to csv
	// err = saveNumerosToCSV(numeros, "numeros.csv")

}

func fetchNumbersBatch(accounts []api.ApiV2010Account, accountant chan Accountant, outerWg *sync.WaitGroup, chWait *sync.WaitGroup, check_nums bool) {
	// create inner wait-group for the current batch
	defer outerWg.Done()
	var innerWg sync.WaitGroup
	for i := range accounts {
		innerWg.Add(1)
		go fetchAllNumbers(accounts[i], accountant, &innerWg, chWait, check_nums)

	}
	innerWg.Wait()
	fmt.Printf("Inner wait-group done\n")

}
func fetchAllNumbers(account api.ApiV2010Account, accountant chan Accountant, wg *sync.WaitGroup, chWait *sync.WaitGroup, check_nums bool) {
	defer wg.Done()
	// var chG sync.WaitGroup
	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: *account.Sid,
			Password: *account.AuthToken,
		},
	)
	// make a list of numbers string called numero
	numero := []string{}
	detailNumero := []Numeros{}

	resp, err := client.Api.ListIncomingPhoneNumber(nil)
	if !check_nums {
		for record := range resp {
			numero = append(numero, *resp[record].PhoneNumber)
			numeros := Numeros{
				PhoneNumber: func() string {
					if resp[record].PhoneNumber == nil {
						return ""
					}
					return *resp[record].PhoneNumber
				}(),
				DateCreated: func() string {
					if resp[record].DateCreated == nil {
						return ""
					}
					return *resp[record].DateCreated
				}(),
				DateUpdated: func() string {
					if resp[record].DateUpdated == nil {
						return ""
					}
					return *resp[record].DateUpdated

				}(),
				AccountName: *account.FriendlyName,
				AccountSid:  *account.Sid,
				Status: func() string {
					if resp[record].Status == nil {
						return ""
					}
					return *resp[record].Status
				}(),
				VoiceUrl: func() string {
					if resp[record].VoiceUrl == nil {
						return ""
					}
					return *resp[record].VoiceUrl

				}(),
				SmsUrl: func() string {
					if resp[record].SmsUrl == nil {
						return ""
					}
					return *resp[record].SmsUrl
				}(),
			}
			detailNumero = append(detailNumero, numeros)

		}
		// Create an accountant object
		__accountant := Accountant{
			AccountSid:   *account.Sid,
			AuthToken:    *account.AuthToken,
			FriendlyName: *account.FriendlyName,
			DateCreated:  *account.DateCreated,
			DateUpdated:  *account.DateUpdated,
			numbers:      numero,
			numeros:      detailNumero,
		}

		chWait.Add(1)
		go func() {
			defer chWait.Done()
			if len(__accountant.numbers) > 0 {
				accountant <- __accountant
			}
		}()
		// wg.Done()
		// chG.Wait()

		if err != nil {
			fmt.Printf("Error fetching numbers: %v\n", err)
		}

	} else {
		if err != nil || len(resp) < 1 {
			__accountant := Accountant{
				AccountSid:   *account.Sid,
				AuthToken:    *account.AuthToken,
				FriendlyName: *account.FriendlyName,
				DateCreated:  *account.DateCreated,
				DateUpdated:  *account.DateUpdated,
			}

			chWait.Add(1)
			go func() {
				defer chWait.Done()
				accountant <- __accountant
			}()
		}

	}

}

func saveAccountantsToCSV(accountants []Accountant, filename string) error {
	// open a file for writing
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	header := []string{"AccountSid", "AuthToken", "DateCreated", "DateUpdated", "FriendlyName", "Numbers"}
	writer := csv.NewWriter(f)
	defer writer.Flush()
	err = writer.Write(header)

	if err != nil {
		return err
	}
	for record := range accountants {
		recordLine := []string{
			accountants[record].AccountSid,
			accountants[record].AuthToken,
			accountants[record].DateCreated,
			accountants[record].DateUpdated,
			accountants[record].FriendlyName,
			fmt.Sprintf("%v", accountants[record].numbers),
		}
		err = writer.Write(recordLine)
		if err != nil {
			fmt.Printf("Error writing record to file: %v\n", err)
		}

	}
	return nil

}

func saveNumerosToCSV(numeros []Numeros, filename string) error {
	// open a file for writing
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	header := []string{"PhoneNumber", "DateCreated", "DateUpdated", "AccountName", "AccountSid", "Status", "VoiceUrl", "SmsUrl"}
	writer := csv.NewWriter(f)
	defer writer.Flush()
	err = writer.Write(header)

	if err != nil {
		return err
	}
	for record := range numeros {
		recordLine := []string{
			numeros[record].PhoneNumber,
			numeros[record].DateCreated,
			numeros[record].DateUpdated,
			numeros[record].AccountName,
			numeros[record].AccountSid,
			numeros[record].Status,
			numeros[record].VoiceUrl,
			numeros[record].SmsUrl,
		}
		err = writer.Write(recordLine)
		if err != nil {
			fmt.Printf("Error writing record to file: %v\n", err)
		}

	}
	return nil

}
