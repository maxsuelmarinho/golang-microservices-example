package model

import (
	"fmt"
)

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (a *Account) ToString() string {
	return fmt.Sprintf("%s %s", a.ID, a.Name)
}
