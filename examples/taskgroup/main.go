package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lei006/go-assist/taskgroup"
)

type TaskTest struct {
	Id int64
}

func (test *TaskTest) Run(task_id int64) {

	num := time.Duration(rand.Intn(1000))
	time.Sleep(num * time.Millisecond)

	fmt.Println("task(", task_id, ")  data=", test.GetID())

}

func (test *TaskTest) GetID() int64 {
	return test.Id
}

func main() {
	taskgroup.Default.SetTaskCount(2)

	for i := 0; i < 100; i++ {

		newTask := TaskTest{
			Id: int64(i),
		}
		taskgroup.Default.Push(&newTask)

	}

}
