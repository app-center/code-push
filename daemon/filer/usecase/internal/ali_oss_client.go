package internal

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/funnyecho/code-push/daemon/filer/domain"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func NewAliOssClient(schemeService domain.ISchemeService) *AliOssClient {
	return &AliOssClient{
		schemeService: schemeService,
	}
}

type AliOssClient struct {
	schemeService domain.ISchemeService
	scheme        *domain.AliOssScheme
	client        *oss.Client
}

func (c *AliOssClient) UpdateScheme(config *domain.AliOssScheme) error {
	err := c.schemeService.UpdateAliOssScheme(config)
	if err == nil {
		c.scheme = config
	}

	return err
}

func (c *AliOssClient) SignFetchURL(key []byte) ([]byte, error) {
	bucket, err := c.GetPackageBucket()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get package bucket: %s", ossBucketPackage)
	}

	url, err := bucket.SignURL(string(key), oss.HTTPGet, 5*60)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to sign url in get method with key: %s", key)
	}

	return []byte(url), nil
}

func (c *AliOssClient) GetClient() (*oss.Client, error) {
	if c.scheme == nil {
		s, err := c.schemeService.AliOssScheme()
		if err != nil {
			return nil, err
		}

		c.scheme = s
	}

	if c.client == nil {
		client, err := oss.New(string(c.scheme.Endpoint), string(c.scheme.AccessKeyId), string(c.scheme.AccessKeySecret))
		if err != nil {
			return nil, err
		}

		c.client = client
	}

	return c.client, nil
}

func (c *AliOssClient) GetPackageBucket() (*oss.Bucket, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get oss client")
	}

	return client.Bucket(ossBucketPackage)
}

func (c *AliOssClient) GeneratePackageObjectKey() string {
	return util.EncodeBase64(uuid.NewV4().String())
}

const (
	ossBucketPackage = "package"
)
