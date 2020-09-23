package alioss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/funnyecho/code-push/daemon"
	"github.com/funnyecho/code-push/daemon/usecase"
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/funnyecho/code-push/pkg/util"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"io"
	"path"
)

func NewAliOssAdapter(endpoint, bucket, accessKeyId, accessKeySecret string, logger log.Logger) (usecase.AliOssAdapter, error) {
	if endpoint == "" || bucket == "" || accessKeyId == "" || accessKeySecret == "" {
		return nil, errors.Wrap(daemon.ErrParamsInvalid, "endpoints, bucket, accessKeyId, accessKeySecret is required")
	}

	return &aliOss{
		endpoint:        endpoint,
		bucket:          bucket,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		Logger:          logger,
	}, nil
}

type aliOss struct {
	endpoint        string
	bucket          string
	accessKeyId     string
	accessKeySecret string

	log.Logger

	client *oss.Client
}

func (o *aliOss) SignFetchURL(key []byte) ([]byte, error) {
	bucket, err := o.getPackageBucket()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get package bucket: %s", o.bucket)
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

	return client.Bucket(o.bucket)
}

func (o *aliOss) generatePackageObjectKey() string {
	return path.Join(ossDir, util.EncodeBase64(uuid.NewV4().String()))
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
	ossDir = "package"
)
