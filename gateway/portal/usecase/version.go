package usecase

import "github.com/funnyecho/code-push/gateway/portal"

func (u *useCase) ReleaseVersion(params *portal.VersionReleaseParams) error {
	return u.codePush.ReleaseVersion(params)
}

func (u *useCase) GetVersion(envId, appVersion []byte) (*portal.Version, error) {
	return u.codePush.GetVersion(envId, appVersion)
}

func (u *useCase) VersionStrictCompatQuery(envId, appVersion []byte) (*portal.VersionCompatQueryResult, error) {
	return u.codePush.VersionStrictCompatQuery(envId, appVersion)
}
