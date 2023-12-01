package core

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	ID_BUCKET_KEY_FORMAT     = "id_bucket_%d%d"
	ID_BUCKET_RESERVED_VALUE = "reserved"
	ID_BUCKET_SIZE           = 1000000
	ALLOCATE_ID_BUCKET_RETRY = 10
)

type idGenerator struct {
	bucket struct {
		end   int64
		index int64
	}
	prefix    string
	milestone int64
	locker    sync.Mutex
}

var ID *idGenerator

/*
* Init the id generator
 */
func initIdGenerator() {
	ID = &idGenerator{}
	// 1. Get the id bucket from redis
	flag := false
	bucket := (time.Now().Unix() % (60 * 60 * 24)) / 60
	for i := 0; i < ALLOCATE_ID_BUCKET_RETRY; i++ {
		log.Printf("Try to allocate bucket: %d", bucket)
		if flag = allocateIdBucket(ID, int64(bucket)); flag {
			break
		}
		bucket++
	}

	if !flag {
		log.Panic("Cannot allocate id generator!")
	}
}

/*
* Generate the id
* @params: void
* @return: string
 */
func (generator *idGenerator) GenerateID() string {
	generator.locker.Lock()
	defer generator.locker.Unlock()
	// Check if the bucket is full
	if generator.bucket.index >= generator.bucket.end {
		bucket := (time.Now().Unix() % (60 * 60 * 24)) / 60
		flag := false
		for i := 0; i < ALLOCATE_ID_BUCKET_RETRY; i++ {
			log.Printf("Try to allocate bucket: %d", bucket)
			if flag = allocateIdBucket(ID, int64(bucket)); flag {
				break
			}
			bucket++
		}

		if flag {
			log.Panic("Cannot reallocate id generator!")
		}
	}
	generator.bucket.index++
	// Generate the id
	return fmt.Sprintf("%s%06d", generator.prefix, generator.bucket.index)
}

/*
 * Allocate id bucket when id generator is initialized
 * or after used all ids is allcated
 * @params: void
 * @return: void
 */
func allocateIdBucket(generator *idGenerator, bucketId int64) bool {
	// 1. Get the id bucket
	now := time.Now().Unix()
	secondPrefix := now / (60 * 60 * 24)

	// Save bucket to local generator
	generator.bucket.index = 0
	generator.bucket.end = ID_BUCKET_SIZE
	generator.prefix = fmt.Sprintf("%d%d", secondPrefix, bucketId)
	generator.milestone = secondPrefix * (60 * 60 * 24)

	return true
}
