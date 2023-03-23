package messageQueue

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
	// 队列的最小容量，防止队列中的消息过少，导致频繁的扩容后又缩容，影响性能
	queMinCap int
	workerNum int
	// 用于传递消息给worker的通道,容量为buf的长度(同样是为了尽快从消息队列中读取消息)
	msgChan chan T
	// 用于通知发送消息的协程，队列中有消息了,waitChan的容量为1，保证Push时不会阻塞
	waitChan chan struct{}
	//由于只有一个goroutine读取队列中的消息，可能发生读饥饿的情况，
	// buf是用于存储消息的缓冲区,保证在低概率抢到锁的情况下，也能读取到足够的队列中的消息
	buf []T
}

// NewSimpleMQ function creates a new SimpleMQ instance and starts the worker goroutines.
// The worker function processes messages from the message channel
// and calls the provided message handler function.
// The message handler must be a safe function that can be called concurrently.
func NewSimpleMQ[T any](workerNum int, msgHandler func(T)) *SimpleMQ[T] {
	var buf []T = make([]T, workerNum*2)
	var Msg chan T = make(chan T, len(buf))
	var Wait chan struct{} = make(chan struct{}, 1)
	var ret = &SimpleMQ[T]{
		queLock:   sync.Mutex{},
		que:       arrayQueue.New[T](),
		workerNum: workerNum,
		msgChan:   Msg,
		waitChan:  Wait,
		buf:       buf,
	}
	ret.queMinCap = 200
	for i := 0; i < workerNum; i++ {
		go ret.worker(msgHandler)
	}

	go sendMsg(ret, msgHandler)
	return ret
}

func (mq *SimpleMQ[T]) worker(msgHandler func(T)) {
	for {
		msg := <-mq.msgChan
		msgHandler(msg)
	}
}

// sendMsg function reads messages from the queue and sends them to the message channel.
// The implementation also includes a wait channel to notify the sendMsg function
// when the queue is not empty.
// 单线程读取队列中的消息，发送到消息通道中
func sendMsg[T any](mq *SimpleMQ[T], msgHandler func(T)) {
	for {
		var empty bool = false
		var msgNum int = 0
		mq.queLock.Lock()
		// 读取队列中的消息，放到buf中，保证在低概率抢到锁的情况下，也能读取到足够的队列中的消息
		for msgNum < len(mq.buf) && !mq.que.Empty() {
			mq.buf[msgNum] = mq.que.Pop()
			msgNum++
		}
		mq.tryShrink()
		if mq.que.Empty() {
			empty = true
		}
		mq.queLock.Unlock()
		//log.Printf("msgNum: %v, empty: %v,msgChan len: %v\n", msgNum, empty, len(mq.msgChan))
		// 发送消息到消息通道中
		for i := 0; i < msgNum; i++ {
			mq.msgChan <- mq.buf[i]
		}
		// 读取完队列中的消息后，若队列为空，
		// 其他goroutine再次调用Push时，必然会给waitChan发送消息
		if empty {
			// 等待队列中有消息
			<-mq.waitChan
		}
	}
}

// 缩小队列的容量(调用需要加锁)
func (mq *SimpleMQ[T]) tryShrink() {
	if mq.que.Len() < mq.que.Cap()/2 && mq.que.Cap() > mq.queMinCap {
		newCap := mq.que.Cap() / 2
		// 释放内存
		mq.que.Resize(newCap)
	}
}

func (mq *SimpleMQ[T]) Push(msg T) {
	mq.queLock.Lock()
	mq.que.Push(msg)
	if mq.Len() == 1 {
		// 通知发送消息的协程,队列中有消息了
		mq.waitChan <- struct{}{}
	}
	mq.queLock.Unlock()
}

func (mq *SimpleMQ[T]) Len() int {
	return mq.que.Len()
}
