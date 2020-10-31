package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// ExecutePipeline is
func ExecutePipeline(jobs ...job) {
	// fmt.Println("JOBS: ", jobs)
	in := make(chan interface{})
	wg := &sync.WaitGroup{}
	for _, jobFunc := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go WorkerPipeline(wg, jobFunc, in, out)
		in = out
		// runtime.Gosched()
	}
	wg.Wait()

}

// WorkerPipeline is
func WorkerPipeline(wg *sync.WaitGroup, jobFunc job, in, out chan interface{}) {
	defer wg.Done()
	defer close(out)
	jobFunc(in, out)
}

// SingleHash is
func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for i := range in {
		data := fmt.Sprintf("%v", i)
		fmt.Println(data)
		crcMd5 := DataSignerMd5(data)
		fmt.Println(crcMd5)
		wg.Add(1)
		go WorkerSingleHash(wg, data, crcMd5, out)
	}
	wg.Wait()
}

// WorkerSingleHash is
func WorkerSingleHash(wg *sync.WaitGroup, data, crcMd5 string, out chan interface{}) {
	defer wg.Done()
	crc32Chan := make(chan string)
	crc32Md5Chan := make(chan string)
	go calculateHash(crc32Chan, data, DataSignerCrc32)
	go calculateHash(crc32Md5Chan, crcMd5, DataSignerCrc32)
	crc32Hash := <-crc32Chan
	crc32Md5Hash := <-crc32Md5Chan
	out <- crc32Hash + "~" + crc32Md5Hash
}

// после хэширования передаем обратно
func calculateHash(ch chan string, data string, f func(string) string) {
	result := f(data)
	ch <- result
}

// MultiHash is
func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for i := range in {
		wg.Add(1)
		go WorkerMultiHash(wg, i, out)
	}
	wg.Wait()
}

// WorkerMultiHash is
func WorkerMultiHash(wg *sync.WaitGroup, h interface{}, ch chan interface{}) {
	defer wg.Done()
	wgInternal := &sync.WaitGroup{}
	hashArray := make([]string, 6)
	for i := 0; i < 6; i++ {
		wgInternal.Add(1)
		data := fmt.Sprintf("%v%v", i, h)
		go calculateMultiHash(wgInternal, i, data, hashArray)
	}
	wgInternal.Wait()
	multiHash := strings.Join(hashArray, "")
	ch <- multiHash
}

func calculateMultiHash(wg *sync.WaitGroup, id int, data string, arr []string) {
	defer wg.Done()
	crc32Hash := DataSignerCrc32(data)
	arr[id] = crc32Hash
}

// CombineResults is
func CombineResults(in, out chan interface{}) {
	var hashArray []string
	for i := range in {
		hashArray = append(hashArray, i.(string))
	}
	sort.Strings(hashArray)
	combineR := strings.Join(hashArray, "_")
	out <- combineR
}
