package models

type EntitySession struct {
	Token   string
	Entity  Entity
	Expired bool
}
