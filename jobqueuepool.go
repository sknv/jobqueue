package jobqueue

import (
	"sync"
)

// JobQueuePool manages a pool of JobQueues.
type JobQueuePool struct {
	pool map[string]*JobQueue
	mx   sync.Mutex
}

// NewJobQueuePool allocates a new JobQueuePool.
func NewJobQueuePool() *JobQueuePool {
	return &JobQueuePool{pool: make(map[string]*JobQueue)}
}

// GetJobQueue creates or returns a job queue with the provided limit and hash function by id.
func (j *JobQueuePool) GetJobQueue(id string, limit int, hash Hasher) *JobQueue {
	j.mx.Lock()
	defer j.mx.Unlock()

	if jq, found := j.pool[id]; found { // try to find an existing job queue
		return jq
	}

	// Store and return a new job queue
	jq := NewJobQueue(limit, hash)
	j.pool[id] = jq
	return jq
}
