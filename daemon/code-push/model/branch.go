package model

import (
	"time"
)

type Branch struct {
	BranchId	   string `json:"branch_id"`
	BranchName     string    `json:"branch_name"`
	BranchAuthHost string    `json:"branch_auth_host"`
	BranchEncToken string    `json:"branch_enc_token"`
	CreateTime     time.Time `json:"create_time"`
}

func NewBranch(id, name, authHost, encToken string, createTime time.Time) *Branch {
	return &Branch{
		BranchId:       id,
		BranchName:     name,
		BranchAuthHost: authHost,
		BranchEncToken: encToken,
		CreateTime:     createTime,
	}
}

type BranchList = []*Branch
