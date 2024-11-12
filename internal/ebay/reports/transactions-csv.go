package reports

import (
	"bufio"
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionsReport struct {
	Transactions []Transaction
}

type Transaction struct {
	TransactionCreationDate       time.Time
	Type                          string
	OrderNumber                   string
	LegacyOrderID                 string
	BuyerUsername                 string
	BuyerName                     string
	PostToCity                    string
	PostToProvinceRegionState     string
	PostToPostcode                string
	PostToCountry                 string
	NetAmount                     decimal.Decimal
	PayoutCurrency                string
	PayoutDate                    string
	PayoutID                      string
	PayoutMethod                  string
	PayoutStatus                  string
	ReasonForHold                 string
	ItemID                        string
	TransactionID                 string
	ItemTitle                     string
	CustomLabel                   string
	Quantity                      int
	ItemSubtotal                  decimal.NullDecimal
	PostageAndPackaging           string
	SellerCollectedTax            string
	EBayCollectedTax              string
	SellerSpecifiedVATRate        decimal.Decimal
	FinalValueFeeFixed            decimal.NullDecimal
	FinalValueFeeVariable         decimal.NullDecimal
	RegulatoryOperatingFee        decimal.NullDecimal
	VeryHighItemNotAsDescribedFee decimal.NullDecimal
	BelowStandardPerformanceFee   decimal.NullDecimal
	InternationalFee              decimal.NullDecimal
	GrossTransactionAmount        decimal.Decimal
	TransactionCurrency           string
	ExchangeRate                  decimal.NullDecimal
	ReferenceID                   string
	Description                   string
}

func LoadTransactionsFile(ctx context.Context, path string) (TransactionsReport, error) {
	var report TransactionsReport
	file, err := os.Open(path)
	if err != nil {
		return report, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for i := 0; i < 11; i++ {
		_, err := reader.ReadString('\n')
		if err != nil {
			return report, err
		}
	}

	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return report, err
	}

	for _, record := range records[1:] { // Skipping header
		creationDate, err := parseDate(record[0])
		if err != nil {
			return report, err
		}
		netAmount, err := parseDecimal(record[10])
		if err != nil {
			return report, err
		}
		quantity, err := parseInt(record[21])
		if err != nil {
			return report, err
		}
		itemSubtotal, err := parseNullDecimal(record[22])
		if err != nil {
			return report, err
		}
		sellerSpecifiedVATRate, err := parsePercentage(record[26])
		if err != nil {
			return report, err
		}
		finalValueFeeFixed, err := parseNullDecimal(record[27])
		if err != nil {
			return report, err
		}
		finalValueFeeVariable, err := parseNullDecimal(record[28])
		if err != nil {
			return report, err
		}
		regulatoryOperatingFee, err := parseNullDecimal(record[29])
		if err != nil {
			return report, err
		}
		veryHighItemNotAsDescribedFee, err := parseNullDecimal(record[30])
		if err != nil {
			return report, err
		}
		belowStandardPerformanceFee, err := parseNullDecimal(record[31])
		if err != nil {
			return report, err
		}
		internationalFee, err := parseNullDecimal(record[32])
		if err != nil {
			return report, err
		}
		grossTransactionAmount, err := parseDecimal(record[33])
		if err != nil {
			return report, err
		}
		exchangeRate, err := parseNullDecimal(record[35])
		if err != nil {
			return report, err
		}

		report.Transactions = append(report.Transactions, Transaction{
			TransactionCreationDate:       creationDate,
			Type:                          record[1],
			OrderNumber:                   record[2],
			LegacyOrderID:                 record[3],
			BuyerUsername:                 record[4],
			BuyerName:                     record[5],
			PostToCity:                    record[6],
			PostToProvinceRegionState:     record[7],
			PostToPostcode:                record[8],
			PostToCountry:                 record[9],
			NetAmount:                     netAmount,
			PayoutCurrency:                record[11],
			PayoutDate:                    record[12],
			PayoutID:                      record[13],
			PayoutMethod:                  record[14],
			PayoutStatus:                  record[15],
			ReasonForHold:                 record[16],
			ItemID:                        record[17],
			TransactionID:                 record[18],
			ItemTitle:                     record[19],
			CustomLabel:                   record[20],
			Quantity:                      quantity,
			ItemSubtotal:                  itemSubtotal,
			PostageAndPackaging:           record[23],
			SellerCollectedTax:            record[24],
			EBayCollectedTax:              record[25],
			SellerSpecifiedVATRate:        sellerSpecifiedVATRate,
			FinalValueFeeFixed:            finalValueFeeFixed,
			FinalValueFeeVariable:         finalValueFeeVariable,
			RegulatoryOperatingFee:        regulatoryOperatingFee,
			VeryHighItemNotAsDescribedFee: veryHighItemNotAsDescribedFee,
			BelowStandardPerformanceFee:   belowStandardPerformanceFee,
			InternationalFee:              internationalFee,
			GrossTransactionAmount:        grossTransactionAmount,
			TransactionCurrency:           record[34],
			ExchangeRate:                  exchangeRate,
			ReferenceID:                   record[36],
			Description:                   record[37],
		})
	}

	return report, err
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("2 Jan 2006", value)
}

func parsePercentage(value string) (decimal.Decimal, error) {
	return parseDecimal(value[:len(value)-1])
}

func parseDecimal(value string) (decimal.Decimal, error) {
	return decimal.NewFromString(value)
}

func parseNullDecimal(value string) (decimal.NullDecimal, error) {
	if value == "--" {
		return decimal.NullDecimal{}, nil
	}

	d, err := decimal.NewFromString(value)
	if err != nil {
		return decimal.NullDecimal{}, nil
	}
	return decimal.NewNullDecimal(d), nil
}

func parseInt(value string) (int, error) {
	if value == "--" {
		return 0, nil
	}

	return strconv.Atoi(value)
}
