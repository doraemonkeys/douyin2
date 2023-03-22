package messageQueue

type MQ[T any] interface {
	// Push push a message to queue
	Push(T)
	// Pop pop a message from queue
	//Pop()
	// PopWithTimeout pop a message from queue with timeout
	//PopWithTimeout(timeout int) (T, error)

	// Len get the length of queue
	Len() int
}
