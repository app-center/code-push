package daemon

import "time"

type Branch struct {
	ID         string
	Name       string
	EncToken   string
	CreateTime time.Time
}

type Env struct {
	BranchId   string
	ID         string
	Name       string
	EncToken   string
	CreateTime time.Time
}

type Version struct {
	EnvId            string
	AppVersion       string
	CompatAppVersion string
	MustUpdate       bool
	Changelog        string
	PackageFileKey   string
	CreateTime       time.Time
}

type BranchList = []*Branch
type EnvList = []*Env
type VersionList = []*Version

type File struct {
	Key        string
	Value      string
	Desc       string
	CreateTime time.Time
	FileMD5    string
	FileSize   int64
}

type FileKey []byte
type FileValue []byte
type FileDesc []byte

type FileMeta struct {
	FileMD5  string
	FileSize int64
}

type AccessTokenIssuer int32

const (
	AccessTokenIssuerSYS    AccessTokenIssuer = 0
	AccessTokenIssuerPORTAL AccessTokenIssuer = 1
	AccessTokenIssuerCLIENT AccessTokenIssuer = 2
)

type AccessTokenClaims struct {
	Issuer   AccessTokenIssuer
	Subject  string
	Audience []byte
}
