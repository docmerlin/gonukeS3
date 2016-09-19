gonukeS3
=============

This is a script for clearing out an s3 bucket.
Its quick and dirty, but useful, because you cannot delete an S3 bucket that has gotten too large.

Usage
=====

```
go build
./gonukeS3 -bucket="bucket-name -region="region-name"
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
