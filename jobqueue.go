// jobqueue package reduces parallelism providing an ability to execute the code sequentially
// depending on the provided job id.

package jobqueue

import (
	"sync"
)

const (
	// DefaultLimit is the default queue limit (prime number is preferable).
	DefaultLimit = 47
)

// Job represents some job to be executed.
type Job func() error

// Hasher calculates an integer hash for the provided string.
type Hasher func(string) int

// JobQueue object.
type JobQueue struct {
	limit   int
	hash    Hasher
	tickets []*sync.Mutex
}

// NewJobQueue allocates a new JobQueue.
func NewJobQueue(limit int, hash Hasher) *JobQueue {
	if limit <= 0 {
		limit = DefaultLimit
	}

	// Allocate a worker queue instance
	jobs := JobQueue{
		limit:   limit,
		tickets: make([]*sync.Mutex, limit),
		hash:    hash,
	}

	// Allocate the tickets
	for i := 0; i < jobs.limit; i++ {
		jobs.tickets[i] = &sync.Mutex{}
	}

	// Provide the default hasher if needed
	if jobs.hash == nil {
		jobs.hash = FNV
	}

	return &jobs
}

// Execute adds a job to the sequentially execution queue.
func (j *JobQueue) Execute(id string, job Job) error {
	ticket := j.ticketByJobID(id)
	return j.executeWithTicket(ticket, job)
}

// ExecuteByHash adds a job with provided hash to the sequentially execution queue.
func (j *JobQueue) ExecuteByHash(hash int, job Job) error {
	ticket := j.ticketByHash(hash)
	return j.executeWithTicket(ticket, job)
}

// executeWithTicket actually executes a job with provided ticket.
func (j *JobQueue) executeWithTicket(ticket *sync.Mutex, job Job) error {
	ticket.Lock()
	defer ticket.Unlock()

	return job() // run the job
}

// ticketByJobID returns a ticket by a provided job id.
func (j *JobQueue) ticketByJobID(id string) *sync.Mutex {
	hash := j.hash(id)
	return j.ticketByHash(hash)
}

// ticketByHash returns a ticket by a provided hash.
func (j *JobQueue) ticketByHash(hash int) *sync.Mutex {
	return j.tickets[hash%j.limit]
}
