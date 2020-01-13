package display

import "time"

type Scheduler struct {
  do   func()
  tick <-chan time.Time
}

func (s *Scheduler) Do(f func()) *Scheduler {
  s.do = f
  return s
}

func (s *Scheduler) Every(d time.Duration) *Scheduler {
  s.tick = time.Tick(d)
  return s
}

func (s *Scheduler) Run() {
  if s.tick != nil && s.do != nil {
    for range s.tick {
      s.do()
    }
  }
}