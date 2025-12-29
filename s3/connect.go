package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (r *repository) New(ctx context.Context) error {
	creds := credentials.NewStaticV4(r.conf.AccessKeyID, r.conf.SecretAccessKey, "")
	opts := minio.Options{
		Creds:  creds,
		Secure: r.conf.Secure,
		Region: r.conf.Region,
	}
	client, err := minio.New(r.conf.Endpoint, &opts)
	if err != nil {
		return err
	}
	fmt.Println("[info] [New] succesed connected to MinioConnect")
	r.client = client

	return nil
}
