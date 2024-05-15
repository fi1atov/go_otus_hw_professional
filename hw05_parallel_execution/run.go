package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(functions []Task, workersCount int, maxErrors int) error {
	tasksChan := make(chan Task) // канал для всех задач
	resultsChan := make(chan error, workersCount-1)
	closeChan := make(chan bool) // сигнальный канал для закрытия

	wg := sync.WaitGroup{}

	// Запускаем наши горутины в числе workersCount
	// они зависают на <-tasksChan - ждут пока там задачи появятся
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		fmt.Println("Запуск горутины-исполнителя")
		go startWorker(tasksChan, resultsChan, closeChan, &wg)
	}

	// далее основная горутина должна закиндать в tasksChan задачи

	_, errorsSlice := func() (int, []error) {
		var errors int
		var inProgress int
		var counter int
		var errorsSlice []error

		// отправка задач в канал - почему то отправятся не все задачи
		fmt.Println("У нас всего заданий: - ", len(functions))
		for i := 0; i < workersCount; i++ {
			fmt.Println("Кладем задание в канал - ", i)
			// time.Sleep(time.Millisecond * 500)
			inProgress++
			tasksChan <- functions[i]
			fmt.Println("I HERE")
		}

		for {
			fmt.Println("I NOT HERE")
			err := <-resultsChan
			inProgress--
			counter++ // счетчик исполненных заданий

			if err != nil {
				errorsSlice = append(errorsSlice, err)
				errors++ // счетчик ошибок
			}

			// если счетчик заданий или ошибок достиг предела - закрыть канал
			if counter == len(functions) || errors == maxErrors {
				fmt.Println("Закрываем канал")
				close(closeChan)
				return counter, errorsSlice
			} else if len(functions)-counter-inProgress > 0 {
				inProgress++
				tasksChan <- functions[workersCount-1+counter]
			}
		}
	}()

	wg.Wait()

	// Если счетчик ошибок достиг предела - вернуть ошибку
	if len(errorsSlice) == maxErrors {
		return ErrErrorsLimitExceeded
	}
	// Иначе все задачи были выполнены
	return nil
}

func startWorker(tasksChan <-chan Task, resultsChan chan<- error, closeChan <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		// ждет пока задача появятся
		case task := <-tasksChan:
			res := task()
			fmt.Println("Результат исполнения - ", res)
			resultsChan <- res // исполнение задачи
		case <-closeChan:
			return // канал закрылся - выходим из функции-исполнителя
		}
	}
}
