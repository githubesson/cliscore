package spinner

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Spinner struct {
	message  string
	stopChan chan struct{}
	doneChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	chars    []string
	delay    time.Duration
}

func New(message string) *Spinner {
	return &Spinner{
		message:  message,
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
		chars:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:    100 * time.Millisecond,
	}
}

func (s *Spinner) Start() {
	if s.running {
		return
	}

	s.running = true
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		i := 0
		for {
			select {
			case <-s.stopChan:
				// Clear the current line
				fmt.Print("\r\033[K")
				s.doneChan <- struct{}{}
				return
			default:
				char := s.chars[i%len(s.chars)]
				fmt.Printf("\r%s %s", char, s.message)
				i++
				time.Sleep(s.delay)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if !s.running {
		return
	}

	s.stopChan <- struct{}{}
	<-s.doneChan
	s.running = false
}

func (s *Spinner) SetMessage(message string) {
	s.message = message
}

func (s *Spinner) SetChars(chars []string) {
	s.chars = chars
}

func (s *Spinner) SetDelay(delay time.Duration) {
	s.delay = delay
}

func WithDots(message string) *Spinner {
	s := New(message)
	s.SetChars([]string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"})
	return s
}

func WithArrows(message string) *Spinner {
	s := New(message)
	s.SetChars([]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"})
	return s
}

func WithBounce(message string) *Spinner {
	s := New(message)
	s.SetChars([]string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"})
	return s
}

func WithSimple(message string) *Spinner {
	s := New(message)
	s.SetChars([]string{".", "..", "...", "...."})
	s.SetDelay(200 * time.Millisecond)
	return s
}

type ProgressBar struct {
	message    string
	total      int
	current    int
	width      int
	percentage int
	running    bool
	stopChan   chan struct{}
	doneChan   chan struct{}
	wg         sync.WaitGroup
}

func NewProgressBar(message string, total int) *ProgressBar {
	return &ProgressBar{
		message:  message,
		total:    total,
		width:    40,
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

func (p *ProgressBar) Start() {
	if p.running {
		return
	}

	p.running = true
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		for {
			select {
			case <-p.stopChan:
				p.render()
				fmt.Println()
				p.doneChan <- struct{}{}
				return
			default:
				p.render()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (p *ProgressBar) Increment() {
	p.current++
	if p.total > 0 {
		p.percentage = int(float64(p.current) / float64(p.total) * 100)
	}
}

func (p *ProgressBar) SetCurrent(current int) {
	p.current = current
	if p.total > 0 {
		p.percentage = int(float64(p.current) / float64(p.total) * 100)
	}
}

func (p *ProgressBar) render() {
	if p.total <= 0 {
		chars := []string{"=", "===", "====", "=====", "====", "===", "="}
		char := chars[time.Now().UnixMilli()/200%int64(len(chars))]
		fmt.Printf("\r%s [%s]", p.message, char)
		return
	}

	filled := int(float64(p.width) * float64(p.current) / float64(p.total))
	bar := strings.Repeat("=", filled) + strings.Repeat(" ", p.width-filled)
	fmt.Printf("\r%s [%s] %d/%d (%d%%)", p.message, bar, p.current, p.total, p.percentage)
}

func (p *ProgressBar) Stop() {
	if !p.running {
		return
	}

	p.stopChan <- struct{}{}
	<-p.doneChan
	p.running = false
}

func (p *ProgressBar) SetMessage(message string) {
	p.message = message
}

func ShowLoading(message string, duration time.Duration) {
	spinner := New(message)
	spinner.Start()
	time.Sleep(duration)
	spinner.Stop()
}
