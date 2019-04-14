package model

import "fmt"

type Quote struct {
	HardwareArchitecture string `json:"hardwareArchitecture"`
	OperatingSystem      string `json:"operatingSystem"`
	IPAddress            string `json:"ipAddress"`
	Quote                string `json:"quote"`
	Language             string `json:"language"`
}

func (q *Quote) ToString() string {
	return fmt.Sprintf("%s %s %s %s %s", q.HardwareArchitecture, q.OperatingSystem, q.IPAddress, q.Language, q.Quote)
}
