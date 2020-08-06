package alioss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/funnyecho/code-push/daemon/filer"
	"github.com/funnyecho/code-push/daemon/filer/usecase"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"io"
)

func NewAliOssAdapter(endpoint, accessKeyId, accessKeySecret string) (usecase.AliOssAdapter, error) {
	if endpoint == "" || accessKeyId == "" || accessKeySecret == "" {
		return nil, errors.Wrap(filer.ErrParamsInvalid, "endpoints, accessKeyId, accessKeySecret is required")
	}

	return &aliOss{
		endpoint:        endpoint,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}, nil
}

type aliOss struct {
	endpoint        string
	accessKeyId     string
	accessKeySecret string

	client *oss.Client
}

func (o *aliOss) SignFetchURL(key []byte) ([]byte, error) {
	bucket, err := o.getPackageBucket()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get package bucket: %s", ossBucketPackage)
	}

	url, err := bucket.SignURL(string(key), oss.HTTPGet, 5*60)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to sign url in get method with key: %s", key)
	}

	return []byte(url), nil
}

func (o *aliOss) Upload(stream io.Reader) ([]byte, error) {
	bucket, bucketErr := o.getPackageBucket()
	if bucketErr != nil {
		return nil, bucketErr
	}

	objectKey := o.generatePackageObjectKey()

	uploadErr := bucket.PutObject(
		objectKey,
		stream,
	)

	if uploadErr != nil {
		return nil, errors.Wrapf(uploadErr, "failed to upload to bucket: %s", bucket.BucketName)
	}

	return []byte(objectKey), nil
}

func (o *aliOss) getPackageBucket() (*oss.Bucket, error) {
	client, err := o.getClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get oss client")
	}

	return client.Bucket(ossBucketPackage)
}

func (o *aliOss) generatePackageObjectKey() string {
	return util.EncodeBase64(uuid.NewV4().String())
}

func (o *aliOss) getClient() (*oss.Client, error) {
	if o.client == nil {
		client, err := oss.New(o.endpoint, o.accessKeyId, o.accessKeySecret)
		if err != nil {
			return nil, err
		}

		o.client = client
	}

	return o.client, nil
}

const (
	ossBucketPackage = "package"
)
