package util

import (
	"fmt"
	"strings"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	CAD Currency = "CAD"
)

// SupportedCurrencies returns a slice of all supported currencies
func SupportedCurrencies() []string {
	return []string{string(USD), string(EUR), string(CAD)}
}

func IsCurrencySupported(currency any) bool {
	switch c := currency.(type) {
	case Currency:
		return isSupportedCurrency(c)
	case string:
		return isSupportedCurrency(Currency(c))
	default:
		return false
	}
}

// isSupportedCurrency is a helper function to check if a Currency is supported
func isSupportedCurrency(c Currency) bool {
	switch c {
	case USD, EUR, CAD:
		return true
	default:
		return false
	}
}

// ParseCurrency converts a string to a Currency if it's supported
func ParseCurrency(s string) (Currency, error) {
	c := Currency(s)
	if !isSupportedCurrency(c) {
		return "", fmt.Errorf("unsupported currency: %s", s)
	}
	return c, nil
}

// String implements the Stringer interface for Currency
func (c Currency) String() string {
	return string(c)
}

// Define a custom error type
type UnsupportedCurrencyError struct {
	Currency string
}

func (e *UnsupportedCurrencyError) Error() string {
	return fmt.Sprintf("unsupported currency: %s. Supported currencies are: %s",
		e.Currency, strings.Join(SupportedCurrencies(), ", "))
}
