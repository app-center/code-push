package model

import "time"

type Env struct {
	BranchId    string    `json:"branch_id"`
	EnvId       string    `json:"env_id"`
	EnvName     string    `json:"env_name"`
	EnvEncToken string    `json:"env_enc_token"`
	CreateTime  time.Time `json:"create_time"`
}

func NewEnv(
	branchId, envId, name, encToken string,
	createTime time.Time,
) *Env {
	return &Env{
		BranchId:    branchId,
		EnvId:       envId,
		EnvName:     name,
		EnvEncToken: encToken,
		CreateTime:  createTime,
	}
}

type EnvMap = map[string]*Env
type EnvList = []*Env
