package s3

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
)

func (r *repository) Ensure(ctx context.Context) error {
	exist, err := r.client.BucketExists(ctx, r.conf.BucketName)
	if err != nil {
		return fmt.Errorf("Ensure: bucket %s existence check failed: %w", r.conf.BucketName, err)
	}

	if !exist {
		fmt.Printf("[info] bucket %s not found, creating...\n", r.conf.BucketName)

		err = r.client.MakeBucket(ctx, r.conf.BucketName, minio.MakeBucketOptions{
			Region: r.conf.Region, // bisa kosong atau "us-east-1"
		})
		if err != nil {
			return fmt.Errorf("connect: failed to create bucket %s: %w", r.conf.BucketName, err)
		}

		fmt.Printf("[info] bucket %s created successfully\n", r.conf.BucketName)
	}

	return nil
}

func (r *repository) Clean(ctx context.Context) error {
	objectsCh := make(chan minio.ObjectInfo)

	// Kirim semua object ke channel
	go func() {
		defer close(objectsCh)
		for object := range r.client.ListObjects(ctx, r.conf.BucketName, minio.ListObjectsOptions{Recursive: true}) {
			fmt.Printf("[info] removing object: %s\n", object.Key)
			objectsCh <- object
		}
	}()

	// Terima error dari RemoveObjects
	errorCh := r.client.RemoveObjects(ctx, r.conf.BucketName, objectsCh, minio.RemoveObjectsOptions{})

	// Tunggu semua selesai, dan catat error kalau ada
	for removeErr := range errorCh {
		if removeErr.Err != nil {
			log.Printf("[warn] failed to delete object %s: %v", removeErr.ObjectName, removeErr.Err)
		}
	}

	// Setelah semua object selesai dihapus, baru hapus bucket
	if err := r.client.RemoveBucket(ctx, r.conf.BucketName); err != nil {
		return fmt.Errorf("RemoveBucket: failed to remove bucket %s: %w", r.conf.BucketName, err)
	}

	fmt.Printf("[info] bucket %s removed successfully\n", r.conf.BucketName)
	return nil
}
