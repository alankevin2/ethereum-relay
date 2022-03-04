package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Storage_client struct {
	Ctx    context.Context
	Client *storage.Client
}
type ObjectAttrs struct {
	Name      string
	Size      int64
	CreatedAt time.Time
}

func CreateClient() (Storage_client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return Storage_client{}, fmt.Errorf("storage.NewClient: %v\n", err)
	}
	return Storage_client{Ctx: ctx, Client: client}, nil
}

func (s *Storage_client) Close() {
	s.Client.Close()
	_, cancel := context.WithTimeout(s.Ctx, time.Second*50)
	cancel()
}

// downloadFile downloads an object.
func (s *Storage_client) DownloadFile(bucket, location, object string) error {
	outputFile, err := os.Create(location)
	outputFile.Close()

	rc, err := s.Client.Bucket(bucket).Object(object).NewReader(s.Ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	err = ioutil.WriteFile(location, data, 0666)

	if err != nil {
		return fmt.Errorf("ioutil.WriteFile: %v", err)
	}

	return nil
}

func (s *Storage_client) UploadFile(bucket, location, object string) error {
	// Open local file.
	f, err := os.Open(location)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	// Upload an object with storage.Writer.
	wc := s.Client.Bucket(bucket).Object(object).NewWriter(s.Ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Printf("Blob %s uploaded.\n", object)
	return nil
}

func (s *Storage_client) DeleteFile(bucket, object string) error {
	o := s.Client.Bucket(bucket).Object(object)
	if err := o.Delete(s.Ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	fmt.Printf("Blob %s deleted.\n", object)
	return nil
}

func (s *Storage_client) ListFilesWithPrefix(bucket, prefix string) ([]ObjectAttrs, error) {
	list := make([]ObjectAttrs, 0)

	it := s.Client.Bucket(bucket).Objects(s.Ctx, &storage.Query{
		Prefix: prefix,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects(): %v", bucket, err)
		}
		list = append(list, ObjectAttrs{attrs.Name, attrs.Size, attrs.Created})
	}
	return list, nil
}

// copyFile copies an object into specified bucket.
func (s *Storage_client) CopyFile(dstBucket, srcBucket, srcObject string) error {
	// dstBucket := "bucket-1"
	// srcBucket := "bucket-2"
	// srcObject := "object"

	dstObject := srcObject + "-copy"
	src := s.Client.Bucket(srcBucket).Object(srcObject)
	dst := s.Client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(s.Ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", dstObject, srcObject, err)
	}

	fmt.Printf("Blob %v in bucket %v copied to blob %v in bucket %v.\n", srcObject, srcBucket, dstObject, dstBucket)
	return nil
}

func (s *Storage_client) MoveFile(bucket, object, target string) error {
	// bucket := "bucket-name"
	// object := "object-name"

	dstName := target
	src := s.Client.Bucket(bucket).Object(object)
	dst := s.Client.Bucket(bucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(s.Ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", dstName, object, err)
	}
	if err := src.Delete(s.Ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}

	fmt.Printf("Blob %v moved to %v.\n", object, dstName)
	return nil
}
