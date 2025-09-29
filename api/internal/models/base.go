package models

type Model interface {
	GetID() int64
	SetID(id int64)
}
