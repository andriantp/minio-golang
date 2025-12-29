package s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
)

func (r *repository) FileList(ctx context.Context, directory string) ([]string, error) {
	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Prefix:    directory,
		Recursive: true,
	}

	var nameFile []string
	objectCh := r.client.ListObjects(ctx, r.conf.BucketName, opts)
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("ListObjects:%w", object.Err)
		}
		name := strings.TrimPrefix(object.Key, directory)
		nameFile = append(nameFile, name)
	}

	return nameFile, nil
}
