package pkg

import (
	"sync"
)

//SenderCollection is a collection of Senders
type SenderCollection struct {
	procs  []*Sender
	config *Config

	m sync.Mutex
}

//NewSenderCollection create SenderCollection of Senders
func NewSenderCollection(config *Config) *SenderCollection {
	subProcesses := new(SenderCollection)

	processes := []*Sender{}
	for _, f := range config.Filters {
		sp, err := NewSender(f, config)
		checkOrFail(err)
		processes = append(processes, sp)
	}

	subProcesses.procs = processes
	subProcesses.config = config
	subProcesses.resetData()
	return subProcesses
}

func (s *SenderCollection) resetData() {
	for _, proc := range s.procs {
		proc.resetData()
	}
}

//appendData appends RowEntry to every filter instance from config
func (s *SenderCollection) appendData(row *RowEntry) {
	s.m.Lock()
	var wg sync.WaitGroup
	wg.Add(len(s.procs))
	for _, proc := range s.procs {
		go func(proc *Sender) {
			defer wg.Done()
			proc.appendIfOk(row)
		}(proc)
	}
	wg.Wait()
	s.m.Unlock()
}

func (s *SenderCollection) sendStats() {
	s.m.Lock()
	var wg sync.WaitGroup
	wg.Add(len(s.procs))

	for _, proc := range s.procs {
		go func(proc *Sender) {
			defer wg.Done()
			proc.sendStats()
		}(proc)
	}
	wg.Wait()

	s.resetData()
	s.m.Unlock()
}
