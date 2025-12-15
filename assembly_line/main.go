package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Employee struct {
	ID    int
	Count int
}

type Item1 struct{}

type Item2 struct{}

type Item3 struct{}

type Item interface {
	// Process 這是一個耗時操作
	Process()
	ItemType() string
}

func (i Item1) ItemType() string {
	return "Item1"
}

func (i Item1) Process() {
	time.Sleep(100 * time.Millisecond)
}

func (i Item2) ItemType() string {
	return "Item2"
}
func (i Item2) Process() {
	time.Sleep(200 * time.Millisecond)
}

func (i Item3) ItemType() string {
	return "Item3"
}
func (i Item3) Process() {
	time.Sleep(300 * time.Millisecond)
}

func (e *Employee) Work(
	ch <-chan interface{},
	wg *sync.WaitGroup,
	logMu *sync.Mutex,
	startTime time.Time,
) {
	fmt.Printf("Employee-%d Start Work\n", e.ID)
	defer wg.Done()

	for v := range ch {
		item, ok := v.(Item)
		if !ok {
			continue
		}

		logMu.Lock()
		fmt.Printf("[%v]Employee-%d START processing %s\n",
			time.Since(startTime),
			e.ID,
			item.ItemType())
		logMu.Unlock()

		item.Process()
		e.Count++

		logMu.Lock()
		fmt.Printf("[%v]Employee-%d END processing %s\n",
			time.Since(startTime),
			e.ID,
			item.ItemType())
		logMu.Unlock()

	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const (
		numEmployees = 5
		numItems     = 10
	)

	startTime := time.Now()
	logMu := &sync.Mutex{}
	taskCh := make(chan interface{}, numItems)

	var employees []*Employee
	for i := 1; i <= numEmployees; i++ {
		employees = append(employees, &Employee{ID: i})
	}

	var wg sync.WaitGroup
	for _, emp := range employees {
		wg.Add(1)
		go emp.Work(taskCh, &wg, logMu, startTime)
	}

	var tasks []interface{}
	for i := 1; i <= numItems; i++ {
		tasks = append(tasks,
			Item1{},
			Item2{},
			Item3{},
		)
	}

	rand.Shuffle(len(tasks), func(i, j int) {
		tasks[i], tasks[j] = tasks[j], tasks[i]
	})

	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()

	fmt.Printf("總耗時: %v\n", time.Since(startTime))
	for _, e := range employees {
		fmt.Printf("Employee-%d processed %d items\n", e.ID, e.Count)
	}
}
