package model

import (
	"fmt"
)

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ServedBy string `json:"servedBy"`
	Quote    Quote  `json:"quote"`
}

type Quote struct {
	Text     string `json:"quote"`
	ServedBy string `json:"ipAddress"`
	Language string `json:"language"`
}

func (a *Account) ToString() string {
	return fmt.Sprintf("%s %s %s", a.ID, a.Name, a.ServedBy)
}
