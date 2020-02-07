package model

import (
	"time"
)

type Branch struct {
	id         string    `json:"branch_id"`
	name       string    `json:"branch_name"`
	authHost   string    `json:"branch_auth_host"`
	encToken   string    `json:"branch_enc_token"`
	createTime time.Time `json:"create_time"`
}

func (b *Branch) BranchId() string {
	return b.id
}

func (b *Branch) BranchName() string {
	return b.name
}

func (b *Branch) BranchAuthHost() string {
	return b.authHost
}

func (b *Branch) BranchEncToken() string {
	return b.encToken
}

func (b *Branch) BranchCreateTime() time.Time {
	return b.createTime
}

type BranchConfig struct {
	Id         string
	Name       string
	AuthHost   string
	EncToken   string
	CreateTime time.Time
}

func NewBranch(config BranchConfig) *Branch {
	return &Branch{
		id:         config.Id,
		name:       config.Name,
		authHost:   config.AuthHost,
		encToken:   config.EncToken,
		createTime: config.CreateTime,
	}
}

type BranchList = []*Branch
