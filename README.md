gonukeS3
=============

This is a script for clearing out an s3 bucket.
Its quick and dirty, but useful, because you cannot delete an S3 bucket that has gotten too large.


Installation
============

1. [Install Go](https://golang.org/doc/install)
1. [Setup your AWS CLI credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-config-files)
1. Clone this repo somewhere on your `$GOPATH`
1. Then, from this repo's directory, run:

  ```sh
  go get github.com/aws/aws-sdk-go
  go build
  ```

Usage
=====

```
./gonukeS3 -bucket="bucket-name" -region="region-name"
```

```
Usage of ./gonukeS3:
  -bucket string
    	AWS S3 Bucket Name
  -delim string
    	AWS S3 Delimiter
  -max-delete int
    	Delim (default 1000)
  -prefix string
    	The Prefix
  -region string
    	AWS Region (default "us-west-2")
  -v	verbose mode?
  -workers int
    	Number of workers (default 10)
```
