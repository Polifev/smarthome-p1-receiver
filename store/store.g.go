package store

import "github.com/polifev/smarthome-p1-receiver/model"

type PowerStore interface {
	PutData(data model.PowerData) error
}
