package BucketUtils

import (
	"context"
	"errors"
	"fmt"
	"goScan/utilityFunctions"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func processGCPFiles(bucket string) error {
	ctx := context.Background()
	var files []string

	projectID, err := utilityFunctions.ReadFileEnvs(".env")
	if err != nil {
		log.Fatal("Cannot read project ID from .env file", err.Error())
	}

	// Creates a client from gcloud cli account on machine
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Leaving the error off since this will be program end anyway
	utilityFunctions.SafeClose(client)

	list, err := ListBuckets(ctx, *client, projectID)
	if err != nil {
		log.Fatalf("Failed to list buckets: %v", err)
	}

	for _, bucket := range list {
		files, err = ListFiles(ctx, *client, bucket)
		if err != nil {
			log.Fatalf("Failed to list files: %v", err)
		}
		fmt.Println(files)
	}

	return nil
}

func ListBuckets(ctx context.Context, client storage.Client, projectID string) ([]string, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	var buckets []string
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
	}

	return buckets, nil
}

// ListFiles lists objects within specified bucket.
func ListFiles(ctx context.Context, client storage.Client, bucket string) ([]string, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var files []string

	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			fmt.Println("Can't list the bucket", bucket, err.Error())
			continue
		}
		files = append(files, attrs.Name+"\n")
	}
	return files, nil
}
