package main // github.com/docmerlin/gonukes3
/*
* Nuke the s3 bucket, its the only way to be sure.
*
 */

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"log"
	"sync"
)

var region string
var bucketName string
var prefix string
var delim string
var maxSizeBucketList int64 = 10000000
var workers int
var maxDelete int
var verbose bool
var total = &Total{}

func init() {
	flag.StringVar(&region, "region", "us-west-2", "AWS Region")
	flag.StringVar(&bucketName, "bucket", "", "AWS S3 Bucket Name")
	flag.StringVar(&prefix, "prefix", "", "The Prefix")
	flag.StringVar(&delim, "delim", "", "AWS S3 Delimiter")
	flag.IntVar(&workers, "workers", 10, "Number of workers")
	flag.IntVar(&maxDelete, "max-delete", 1000, "Delim")
	flag.BoolVar(&verbose, "v", false, "verbose mode?")
	flag.Parse()
}

// Total is a Mutex for Totals
type Total struct {
	sync.RWMutex
	count int
}

// Add adds numbers to a Total
func (t *Total) Add(n int) {
	t.Lock()
	t.count = t.count + n
	t.Unlock()
}

// Count is a getter for total.
func (t *Total) Count() (c int) {
	t.RLock()
	c = t.count
	t.RUnlock()
	return c
}

func work(objs chan []*s3.ObjectIdentifier, bucket *string, s3Instance *s3.S3) {
	for true {
		objs := <-objs
		if _, err := s3Instance.DeleteObjects(&s3.DeleteObjectsInput{Bucket: bucket, Delete: &s3.Delete{Objects: objs}}); err != nil {
			log.Printf("WE HAVE ENCOUNTERED AN ERROR: %s", err)
			return
		}
		total.Add(len(objs))
		log.Printf("Deleted a total of %d items.", total.Count())
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	log.Printf("Nuking Bucket: %s", bucketName)
	s3Instance := s3.New(session.New(), &aws.Config{Region: aws.String("us-west-2")})

	keys := make(chan []*s3.ObjectIdentifier)
	wg := &sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			work(keys, &bucketName, s3Instance)
			wg.Done()
		}()
	}
	marker := ""
	for true {
		resp, err := s3Instance.ListObjectVersions(&s3.ListObjectVersionsInput{Bucket: &bucketName, MaxKeys: &maxSizeBucketList})
		if err != nil {
			log.Printf("Error in getting bucket info: %s", err.Error())
			return
		}
		accumulator := make([]*s3.ObjectIdentifier, 0, maxDelete)
		for i, obj := range resp.Versions {
			versionID := *obj.VersionId
			key := *obj.Key
			accumulator = append(accumulator, &s3.ObjectIdentifier{Key: &key, VersionId: &versionID})
			if verbose {
				log.Printf("Queued object for deletion: %s", *obj.Key)
			}
			if (i+1)%maxDelete == 0 || i+1 == len(resp.Versions) {
				objectsToDelete := make([]*s3.ObjectIdentifier, len(accumulator), maxDelete)
				copy(objectsToDelete, accumulator)
				keys <- objectsToDelete
			}
			marker = *obj.Key
		}
		if ((resp.IsTruncated == nil || !*(resp.IsTruncated)) && int64(len(resp.Versions)) < maxSizeBucketList) || len(resp.Versions) == 0 {
			if len(marker) == 0 {
				return
			}
			marker = ""
		}
	}
	log.Print("Waiting for deleters to finish")
	wg.Wait()
}
