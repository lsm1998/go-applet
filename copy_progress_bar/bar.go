package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Bar interface {
	// Done job part
	Done(val float64) bool

	// Cancel job
	Cancel()
}

type progressBar struct {
	total        float64       // 总进度
	current      float64       // 当前进度
	filler       string        // 进度填充字符
	fillerLength int64         // 进度条长度
	timeFormat   string        // 进度条时间格式
	interval     time.Duration // 打印时间间隔
	begin        time.Time     // 任务开始时间
	tail         string
	cancelC      chan struct{} // 取消channel
	ticker       *time.Ticker
	done         bool
	*bufio.Writer
}

type BarOption = func(*progressBar)

func (p *progressBar) Percent() float64 {
	return p.current * 100 / p.total
}

func (p *progressBar) Done(val float64) bool {
	if p.current+val < p.total {
		p.current += val
		return false
	} else if p.done {
		return true
	} else {
		p.current, p.done = p.total, true
		fmt.Println(p.progressString())
		return true
	}
}

func (p *progressBar) Cancel() {
	p.cancelC <- struct{}{}
}

// getEta 获取eta时间
func (p *progressBar) getEta(now time.Time) string {
	eta := float64(now.Unix()-p.begin.Unix()) * 100 / (p.Percent() + 1)
	return p.begin.Add(time.Second * time.Duration(eta)).Format(p.timeFormat)
}

func (p *progressBar) tryProgressString() {
	if p.current < p.total {
		//if _, err := p.Writer.WriteString(p.progressString()); err != nil {
		//	panic(err)
		//}
		//if err := p.Flush(); err != nil {
		//	panic(err)
		//}
		fmt.Printf("%s", p.progressString())
	}
}

func (p *progressBar) progressString() string {
	fills := p.Percent() * float64(p.fillerLength) / 100
	chunks := make([]string, p.fillerLength, p.fillerLength)
	var flag bool
	for i := int64(0); i < p.fillerLength; i++ {
		switch {
		case float64(i) < fills:
			chunks[i] = p.filler
		case !flag:
			flag = !flag
			chunks[i] = p.tail
		default:
			chunks[i] = " "
		}
	}
	now := time.Now()
	eta := p.getEta(now)
	qps := p.current / float64(now.Unix()-p.begin.Unix()+1)
	return fmt.Sprintf("\r[%s]%.2f/%.2f [eta]%s [qps]%.2f ", strings.Join(chunks, ""), p.current, p.total, eta, qps)
}

func WithFiller(filler string) BarOption {
	return func(bar *progressBar) {
		bar.filler = filler
	}
}

func WithFillerLength(fillerLength int64) BarOption {
	return func(bar *progressBar) {
		bar.fillerLength = fillerLength
	}
}

func WithInterval(interval time.Duration) BarOption {
	return func(bar *progressBar) {
		bar.interval = interval
	}
}

func WithTail(tail string) BarOption {
	return func(bar *progressBar) {
		bar.tail = tail
	}
}

// NewProgressBar new
func NewProgressBar(total float64, opts ...BarOption) Bar {
	bar := &progressBar{
		total:        total,
		filler:       "=",
		fillerLength: 25,
		timeFormat:   "15:04:05", // 2006-01-02T15:04:05
		interval:     time.Second,
		begin:        time.Now(),
		cancelC:      make(chan struct{}),
		tail:         ">",
		Writer:       bufio.NewWriter(os.Stdout),
	}
	for _, opt := range opts {
		opt(bar)
	}
	bar.ticker = time.NewTicker(bar.interval)
	go func() {
		for bar.current < bar.total {
			select {
			case <-bar.ticker.C:
				bar.tryProgressString()
			case <-bar.cancelC:
				bar.ticker.Stop()
				return
			}
		}
	}()
	return bar
}
