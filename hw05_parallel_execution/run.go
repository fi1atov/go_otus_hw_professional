package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(functions []Task, workersCount int, maxErrors int) error {
	if maxErrors <= 0 {
		return ErrErrorsLimitExceeded
	}
	tasksChan := make(chan Task) // канал для всех задач
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

	errorsSlice := func() []error {
		var errors int
		var inProgress int
		var counter int
		var errorsSlice []error

		// отправка задач в канал
		// отправляем не все задачи - а то количество которое равно кол-ву созданных горутин
		for i := 0; i < workersCount; i++ {
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

			// если счетчик заданий или ошибок достиг предела - закрыть канал
			if counter == len(functions) || errors == maxErrors {
				close(closeChan)
				return errorsSlice
			} else if len(functions)-counter-inProgress > 0 { // а тут отправляем все остальные задачи
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
			resultsChan <- task() // исполнение задачи
		case <-closeChan:
			return // канал закрылся - выходим из функции-исполнителя
		}
	}
}
