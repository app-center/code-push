package sys

import "time"

type Branch struct {
	ID         string
	Name       string
	EncToken   string
	CreateTime time.Time
}
