package main

import (
	"atp/storage/s3"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	flag.Usage = func() {
		log.Println("Usage: ")
		log.Printf("      go run . methode file")
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	methode := flag.Args()[0]
	fmt.Printf("methode:%s\n", methode)

	file := ""
	if len(flag.Args()) > 1 {
		file = flag.Args()[1]
		fmt.Printf("file   :%s\n", file)
	}

	conf := s3.Cloud{
		Region:          "",
		Endpoint:        "localhost:9000", //api url
		Token:           "",
		Secure:          false, // if use https set true
		AccessKeyID:     "admin",
		SecretAccessKey: "admin@123",
		BucketName:      "minio-bucket",
		ContentType:     "application/octet-stream",
		Expire:          15 * time.Minute, //expire url download
	}

	ctx := context.Background()

	// Create S3 repository
	repo, err := s3.NewCloud(conf)
	if err != nil {
		log.Fatalf("[main] NewCloud:%v", err)
	}

	// New to S3
	if err := repo.New(ctx); err != nil {
		log.Fatalf("[main] New:%v", err)
	}

	// Ensure bucket exists
	if err := repo.Ensure(ctx); err != nil {
		log.Fatalf("[main] Ensure:%v", err)
	}

	upload := "./example/upload/"
	download := "./example/download/"
	images := "images/"

	switch methode {
	case "list": // go run . list
		list, err := repo.FileList(ctx, images)
		if err != nil {
			log.Fatalf("[main] list->%v", err)
		}
		if len(list) == 0 {
			log.Fatal("[main] file is empty")
		}
		for i, name := range list {
			fmt.Printf("[%d] name:%s\n", i+1, name)
		}

	case "upload":
		path := upload + file
		object := images + file
		fmt.Printf("path     :%s\n", path)
		fmt.Printf("object   :%s\n", object)

		hash, err := hasher(path)
		if err != nil {
			log.Fatalf("[main] Upload hash:%v", err)
		}
		hash256 := hex.EncodeToString(hash.Sum(nil))
		fmt.Printf("Hash [%s]\n\n", hash256)

		stat, err := repo.Upload(ctx, path, object, hash)
		if err != nil {
			log.Fatalf("[main] Upload:%v", err)
		}

		if hash256 != stat.UserMetadata["Checksum-Sha256"] {
			log.Fatalf("[main] Hash [%s] is NOT match", hash256)
		}
		fmt.Println("[info] Hash is match")

		fmt.Printf("Key :%s\n", stat.Key)
		fmt.Printf("ETag:%s\n", stat.ETag)
		fmt.Printf("Size: %.2f KB\n", float64(stat.Size)/1024)
		fmt.Printf("Last:%s\n", stat.LastModified.String())

	case "download":
		path := download + file
		object := images + file
		fmt.Printf("path     :%s\n", path)
		fmt.Printf("object   :%s\n\n", object)

		stat, err := repo.Download(ctx, path, object)
		if err != nil {
			log.Fatalf("[main] Download:%v", err)
		}

		fmt.Printf("Name :%s\n", stat.Key)
		fmt.Printf("ETag:%s\n", stat.ETag)
		fmt.Printf("Size: %.2f KB\n", float64(stat.Size)/1024)
		fmt.Printf("Last:%s\n", stat.LastModified.String())
		hash, err := hasher(path)
		if err != nil {
			log.Fatalf("[main] Upload hash:%v", err)
		}
		hash256 := hex.EncodeToString(hash.Sum(nil))
		fmt.Printf("Hash [%s]\n", hash256)

		if hash256 != stat.UserMetadata["Checksum-Sha256"] {
			log.Fatalf("[main] Hash [%s] is NOT match", hash256)
		}

		fmt.Println("[info] Hash is match")
		fmt.Printf("check file at here:%s\n\n", path)

	case "url":
		object := images + file
		url, err := repo.URLDownload(ctx, object)
		if err != nil {
			log.Fatalf("[main] URLDownload:%v", err)
		}
		fmt.Printf("url:[%s]\n", url)

	case "clean":
		if err := repo.Clean(ctx); err != nil {
			log.Fatalf("[main] Clean:%v", err)
		}

		fmt.Println("✅ Bucket deleted successfully")
	}
}

func hasher(path string) (hash.Hash, error) {
	h := sha256.New()

	f, err := os.Open(path)
	if err != nil {
		return h, fmt.Errorf("open->%w", err)
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return h, fmt.Errorf("copy->%w", err)
	}

	return h, nil
}
