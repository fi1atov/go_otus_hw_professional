package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run will run slice of functions on workersCount workers while get maxErrors errors.
func Run(functions []Task, workersCount int, maxErrors int) error {
	tasksChan := make(chan Task) // чан для всех задач
	resultsChan := make(chan error, workersCount-1)
	closeChan := make(chan bool) // сигнальный канал для закрытия

	wg := sync.WaitGroup{}

	// Запускаем наши горутины в числе workersCount
	// они зависают на <-tasksChan - ждут пока там задачи появятся
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go startWorker(tasksChan, resultsChan, closeChan, &wg)
	}

	// далее основная горутина должна закиндать в tasksChan задачи

	counter, errorsSlice := func() (int, []error) {
		var errors int
		var inProgress int
		var counter int
		var errorsSlice []error

		// отправка задач в канал - почему то отправятся не все задачи
		fmt.Println("У нас всего заданий: - ", len(functions))
		for i := 0; i < len(functions); i++ {
			fmt.Println("Кладем задание в канал - ", i)
			inProgress++
			tasksChan <- functions[i]
		}

		for {
			err := <-resultsChan
			inProgress--
			counter++ // счетчик исполненных заданий

			if err != nil {
				errorsSlice = append(errorsSlice, err)
				errors++ // счетчик ошибок
			}

			if counter == len(functions) || errors == maxErrors {
				close(closeChan)
				return counter, errorsSlice
			} else if len(functions)-counter-inProgress > 0 {
				inProgress++
				tasksChan <- functions[workersCount-1+counter]
			}
		}
	}()

	wg.Wait()
	fmt.Println(counter)
	fmt.Println(errorsSlice)
	return ErrErrorsLimitExceeded
}

func startWorker(tasksChan <-chan Task, resultsChan chan<- error, closeChan <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		// ждет пока задача появятся
		case task := <-tasksChan:
			resultsChan <- task() // исполнение задачи
		case <-closeChan:
			return // канал закрылся - выходим из функции-исполнителя
		}
	}
}
