package model

import "time"

type Env struct {
	branchId   string    `json:"branch_id"`
	id         string    `json:"env_id"`
	name       string    `json:"env_name"`
	encToken   string    `json:"env_enc_token"`
	createTime time.Time `json:"create_time"`
}

func (e *Env) BranchId() string {
	return e.branchId
}

func (e *Env) Id() string {
	return e.id
}

func (e *Env) Name() string {
	return e.name
}

func (e *Env) EncToken() string {
	return e.encToken
}

func (e *Env) CreateTime() time.Time {
	return e.createTime
}

type EnvConfig struct {
	BranchId   string
	Id         string
	Name       string
	EncToken   string
	CreateTime time.Time
}

func NewEnv(config EnvConfig) Env {
	return Env{
		branchId:   config.BranchId,
		id:         config.Id,
		name:       config.Name,
		encToken:   config.EncToken,
		createTime: config.CreateTime,
	}
}

type EnvMap = map[string]Env
type EnvList = []Env
