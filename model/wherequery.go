package model

type WhereQuery struct {
	Path  string
	Op    string
	Value interface{}
}
