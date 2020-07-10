package hashmap

import (
	"fmt"
)

type KV struct {
	Key   string
	Value string
}

// link
type LinkNode struct {
	Data     KV
	NextNode *LinkNode
}

func CreateLinkNode() *LinkNode {
	var linkNode = &LinkNode{KV{"", ""}, nil}
	return linkNode
}

func (link *LinkNode) AddNode(data KV) int {
	var count = 0

	// find the tail
	tail := link
	for {
		count += 1
		if tail.NextNode == nil {
			break
		} else {
			tail = tail.NextNode
		}
	}

	var newNode = &LinkNode{data, nil}
	tail.NextNode = newNode

	return count + 1
}

// HashMap木桶（数组）的个数
const BucketCount = 16

type HashMap struct {
	Buckets [BucketCount]*LinkNode
}

func CreateHashMap() *HashMap {
	myMap := &HashMap{}

	for i := 0; i < BucketCount; i++ {
		myMap.Buckets[i] = CreateLinkNode()
	}

	return myMap
}

// 散列算法
func HashCode(key string) int {
	var sum = 0

	for i := 0; i < len(key); i++ {
		sum += int(key[i])
	}

	return (sum % BucketCount)
}

func (myMap *HashMap) AddKeyValue(key string, value string) {
	var mapIndex = HashCode(key)
	var link = myMap.Buckets[mapIndex]
	if link.Data.Key == "" && link.NextNode == nil {
		link.Data.Key = key
		link.Data.Value = value

		fmt.Printf("node key: %v add to buckets %d first node\n", key, mapIndex)
	} else {
		index := link.AddNode(KV{key, value})
		fmt.Printf("node key: %v add to buckets %d %dth node\n", key, mapIndex, index)
	}
}

func (myMap *HashMap) GetValueForKey(key string) string {
	var mapIndex = HashCode(key)
	var link = myMap.Buckets[mapIndex]
	var value string

	head := link
	for {
		if head.Data.Key == key {
			value = head.Data.Value
			break
		} else {
			head = head.NextNode
		}
	}

	return value
}
