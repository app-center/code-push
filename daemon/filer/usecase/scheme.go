package usecase

import (
	"github.com/funnyecho/code-push/daemon/filer/domain"
	"github.com/funnyecho/code-push/daemon/filer/usecase/internal"
)

type IScheme interface {
	UpdateAliOssScheme(config *AliOssSchemeConfig) error
}

func NewSchemeUseCase(config SchemeUseCaseConfig) IScheme {
	return &schemeUseCase{
		aliOssClient: internal.NewAliOssClient(config.SchemeService),
	}
}

type schemeUseCase struct {
	aliOssClient *internal.AliOssClient
}

func (s *schemeUseCase) UpdateAliOssScheme(config *AliOssSchemeConfig) error {
	params := &domain.AliOssScheme{}

	if config.Endpoint != nil {
		params.Endpoint = config.Endpoint
	}

	if config.AccessKeyId != nil {
		params.AccessKeyId = config.AccessKeyId
	}

	if config.AccessKeySecret != nil {
		params.AccessKeySecret = config.AccessKeySecret
	}

	return s.aliOssClient.UpdateScheme(params)
}

type SchemeUseCaseConfig struct {
	SchemeService domain.ISchemeService
}

type AliOssSchemeConfig struct {
	Endpoint        []byte
	AccessKeyId     []byte
	AccessKeySecret []byte
}
