package s3

import (
	"context"
	"fmt"
	"hash"

	"github.com/minio/minio-go/v7"
)

func (r *repository) Upload(ctx context.Context, path, objectName string, h hash.Hash) (*minio.ObjectInfo, error) {
	opts := minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"checksum-sha256": fmt.Sprintf("%x", h.Sum(nil)),
		},

		ContentType: r.conf.ContentType,
	}

	_, err := r.client.FPutObject(ctx, r.conf.BucketName, objectName, path, opts)
	if err != nil {
		return nil, fmt.Errorf("FPutObject:%w", err)
	}

	stat, err := r.client.StatObject(ctx, r.conf.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("StatObject:%w", err)
	}

	return &stat, nil
}
