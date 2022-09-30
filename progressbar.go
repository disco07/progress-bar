package progressbar

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Bar struct {
	state  state
	option option
	theme  Theme
}

type state struct {
	percent          float64
	current          int64
	currentGraphRate int
}

type Theme struct {
	rate       string
	GraphType  string
	GraphStart string
	GraphEnd   string
	GraphWidth int64
}

type option struct {
	total     int64
	startTime time.Time
	sync.Mutex
}

func (b *Bar) SetTheme(t Theme) {
	if t.GraphType != "" {
		b.theme.GraphType = t.GraphType
	}
	if t.GraphWidth != 0 {
		b.theme.GraphWidth = t.GraphWidth
	}
	if t.GraphStart != "" {
		b.theme.GraphStart = t.GraphStart
	}
	if t.GraphEnd != "" {
		b.theme.GraphEnd = t.GraphEnd
	}
}

func New(end int64) *Bar {
	return &Bar{
		state: state{
			percent: getPercent(int64(0), end),
			current: int64(0),
		},
		theme: Theme{
			GraphType:  "█",
			GraphStart: "[",
			GraphEnd:   "]",
			GraphWidth: 50,
		},
		option: option{
			total:     end,
			startTime: time.Now(),
		},
	}
}

func getPercent(current, total int64) float64 {
	return 100 * (float64(current) / float64(total))
}

func (b *Bar) view() error {
	last := b.state.percent
	b.state.percent = getPercent(b.state.current, b.option.total)
	lastGraphRate := b.state.currentGraphRate
	b.state.currentGraphRate = int(b.state.percent / 100.0 * float64(b.theme.GraphWidth))
	if b.state.percent != last {
		b.theme.rate += strings.Repeat(b.theme.GraphType, b.state.currentGraphRate-lastGraphRate)
	}
	secondsLeft := time.Since(b.option.startTime).Seconds() / float64(b.state.current) * (float64(b.option.total) - float64(b.state.current))
	fmt.Printf(
		"\r%s%-*s%s%3d%% %4d/%d (%v-%v)",
		b.theme.GraphStart,
		b.theme.GraphWidth,
		b.theme.rate,
		b.theme.GraphEnd,
		int(b.state.percent),
		b.state.current,
		b.option.total,
		time.Since(b.option.startTime).Round(time.Second),
		time.Duration(secondsLeft)*time.Second,
	)

	return nil
}

// Add is a func who add the number passed as a parameter to the progress bar.
func (b *Bar) Add(num int) error {
	b.option.Lock()
	defer b.option.Unlock()
	if b.option.total == 0 {
		return errors.New("the end must be greater than 0")
	}

	currentNum := int64(num)
	b.state.current += currentNum
	if b.state.current > b.option.total {
		return errors.New("current exceeds total")
	}
	b.view()
	return nil
}

// Default is a basic usage of progress bar.
// In parameter, the max size of things you want to view progress.
// It returns a pointer of Bar.
func Default(end int64) *Bar {
	return New(end)
}
