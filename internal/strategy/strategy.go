package strategy

import "sync"

type RoundRobin struct {
	Mu     sync.Mutex
	CurAdd int
}

func (r *RoundRobin) RoundRobin(address []string) string {
	r.Mu.Lock()
	selectedAddress := address[r.CurAdd]
	r.CurAdd = (r.CurAdd + 1) % len(address)
	r.Mu.Unlock()

	return selectedAddress
}
