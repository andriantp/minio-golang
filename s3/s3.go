package s3

import (
	"context"
	"hash"
	"time"

	"github.com/minio/minio-go/v7"
)

type Cloud struct {
	Region          string
	Endpoint        string
	Secure          bool
	Token           string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	ContentType     string
	Expire          time.Duration
}

type repository struct {
	client *minio.Client
	conf   Cloud
}

func NewCloud(conf Cloud) (RepositoryI, error) {
	return &repository{
		conf:   conf,
		client: nil,
	}, nil
}

type RepositoryI interface {
	New(ctx context.Context) error

	Ensure(ctx context.Context) error
	Clean(ctx context.Context) error

	FileList(ctx context.Context, directory string) ([]string, error)
	Upload(ctx context.Context, path, objectName string, h hash.Hash) (*minio.ObjectInfo, error)
	Download(ctx context.Context, path, objectName string) (*minio.ObjectInfo, error)
	URLDownload(ctx context.Context, objectName string) (string, error)
}
