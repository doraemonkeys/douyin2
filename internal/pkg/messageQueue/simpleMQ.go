package messagequeue

import (
	"sync"

	"github.com/Doraemonkeys/arrayQueue"
)

// This file contains the implementation of a SimpleMQ data structure in Go.
// SimpleMQ is a simple message queue that allows multiple workers to process messages concurrently.
// The implementation uses an array-based queue and two channels to communicate between the workers
// and the main thread. The code defines the SimpleMQ struct and its methods,
// including Push and Len, which allow messages to be added to the queue
// and the current length of the queue to be retrieved.

type SimpleMQ[T any] struct {
	que     *arrayQueue.Queue[T]
	queLock sync.Mutex
	msgChan chan T
	wait    chan struct{}
}

// NewSimpleMQ function creates a new SimpleMQ instance and starts the worker goroutines.
// The worker function processes messages from the message channel
// and calls the provided message handler function.
// The message handler must be a safe function that can be called concurrently.
func NewSimpleMQ[T any](workerNum int, msgHandler func(T) error) *SimpleMQ[T] {
	var Msg chan T = make(chan T)
	var ret = &SimpleMQ[T]{
		queLock: sync.Mutex{},
		que:     arrayQueue.New[T](),
		msgChan: Msg,
	}
	for i := 0; i < workerNum; i++ {
		go ret.worker(msgHandler)
	}

	go sendMsg(ret, msgHandler)
	return ret
}

func (mq *SimpleMQ[T]) worker(msgHandler func(T) error) {
	for {
		msg := <-mq.msgChan
		msgHandler(msg)
	}
}

// sendMsg function reads messages from the queue and sends them to the message channel.
// The implementation also includes a wait channel to notify the sendMsg function
// when the queue is not empty.
// 单线程读取队列中的消息，发送到消息通道中
func sendMsg[T any](mq *SimpleMQ[T], msgHandler func(T) error) {
	for {
		mq.queLock.Lock()
		msg := mq.que.Pop()
		if mq.que.Len() < mq.que.Cap()/2 {
			newCap := mq.que.Cap() / 2
			// 释放内存
			mq.que.Resize(newCap)
		}
		mq.queLock.Unlock()

		mq.msgChan <- msg
		if mq.que.Len() == 0 {
			// 等待队列中有消息
			<-mq.wait
		}
	}
}

func (mq *SimpleMQ[T]) Push(msg T) {
	mq.queLock.Lock()
	mq.que.Push(msg)
	mq.queLock.Unlock()
	if mq.Len() == 1 {
		// 通知发送消息的协程,队列中有消息了
		mq.wait <- struct{}{}
	}
}

func (mq *SimpleMQ[T]) Len() int {
	return mq.que.Len()
}
