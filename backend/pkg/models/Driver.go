package models

import "time"

type Driver struct {
	ID             string
	Name           string
	Surname        string
	Patronymic     string
	BirthDate      time.Time
	PassportSeries string
	Snils          string
	LicenseSeries  string
}
