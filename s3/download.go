package s3

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
)

func (r *repository) Download(ctx context.Context, path, objectName string) (*minio.ObjectInfo, error) {
	opts := minio.GetObjectOptions{
		Checksum: true,
	}

	if err := r.client.FGetObject(ctx, r.conf.BucketName, objectName, path, opts); err != nil {
		return nil, fmt.Errorf("FGetObject:%w", err)
	}

	stat, err := r.client.StatObject(ctx, r.conf.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("StatObject:%w", err)
	}

	return &stat, nil
}

func (r *repository) URLDownload(ctx context.Context, objectName string) (string, error) {
	fileName := objectName[strings.LastIndex(objectName, "/")+1 : int(len(objectName))]
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+fileName+"\"")

	// Generates a presigned url which expires
	presignedURL, err := r.client.PresignedGetObject(ctx, r.conf.BucketName, objectName, r.conf.Expire, reqParams)
	if err != nil {
		return "", fmt.Errorf("PresignedGetObject:%w", err)
	}
	url := presignedURL.String()

	return url, nil
}
