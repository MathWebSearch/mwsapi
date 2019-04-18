package utils

import "sync"

// SyncWorker represents a worker running syncronous code
type SyncWorker struct {
	done chan bool // channel that gets triggered when done

	worker           *sync.Mutex
	hasActiveChannel bool
	activeChannel    int

	messages chan *SyncMessage // for receiving new messages
	backlog  []*SyncMessage    //backlog of messages
}

// NewSyncWorker makes a new SyncWorker
func NewSyncWorker(capacity int) (worker *SyncWorker) {
	worker = &SyncWorker{
		done:     make(chan bool),
		messages: make(chan *SyncMessage, capacity),

		worker: &sync.Mutex{},
	}
	go worker.workerThread()
	return
}

// SyncMessage represents a single syncronized message
type SyncMessage struct {
	channel int
	close   bool
	code    *func()
}

// Work performs some syncronous work
func (worker *SyncWorker) Work(channel int, code func()) {
	worker.messages <- &SyncMessage{
		channel: channel,
		code:    &code,
	}
}

// Close closes a channel and informs the caller that there are no more messages
func (worker *SyncWorker) Close(channel int) {
	worker.messages <- &SyncMessage{
		channel: channel,
		close:   true,
	}
}

// worker thread starts reading all messages
func (worker *SyncWorker) workerThread() {
	// aquire the lock
	worker.worker.Lock()
	defer worker.worker.Unlock()

	// grab the current backlog
	backlog := worker.backlog
	worker.backlog = nil

	for i, m := range backlog {
		closed := worker.processMessage(m)
		if closed {
			worker.backlog = append(worker.backlog, backlog[(i+1):]...)
			go worker.workerThread()
			return
		}
	}

	for m := range worker.messages {
		closed := worker.processMessage(m)
		if closed {
			go worker.workerThread()
			return
		}
	}

	worker.done <- true
	close(worker.done)
}

// Wait waits for processing to complete
func (worker *SyncWorker) Wait() {
	close(worker.messages)
	<-worker.done
	return
}

func (worker *SyncWorker) processMessage(message *SyncMessage) bool {

	// if we don't have an active channel set it to the message
	if !(worker.hasActiveChannel) {
		worker.hasActiveChannel = true
		worker.activeChannel = message.channel
	}

	// if it is a different channel, append it to the backlog
	if message.channel != worker.activeChannel {
		worker.backlog = append(worker.backlog, message)
		return false
	}

	// if we have some code, run it
	if message.code != nil {
		(*message.code)()
	}

	// if we are asked to close the channel, close it
	if message.close {
		worker.hasActiveChannel = false
		return true
	}

	return false
}
