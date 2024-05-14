package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksChan := make(chan Task)
	resultsChan := make(chan error, n-1)
	closeChan := make(chan bool)
	exitChan := make(chan bool)

	// запускаются n горутин
	for i := 0; i < n; i++ {
		go startWorker(tasksChan, resultsChan, closeChan, exitChan)
	}
	return nil
}

func startWorker(tasksChan <-chan Task, resultsChan chan<- error, closeChan <-chan bool, exitChan chan<- bool) {
	// Функция startWorker
	// также отсылает сигнал в exitChan, когда завершает работу.
	defer func() {
		exitChan <- true
	}()

	// В каждой горутине происходит выбор из нескольких возможных операций
	for {
		select {
		// Принятие задачи из tasksChan и выполнение её.
		// Результат выполнения отправляется в resultsChan.
		case task := <-tasksChan:
			resultsChan <- task()
		// Ожидание сигнала о закрытии (closeChan).
		// Если такой сигнал приходит, горутина завершает свою работу.
		case <-closeChan:
			return
		}
	}
}
