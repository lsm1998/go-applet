package main

import (
	"fmt"
	"strings"
	"time"
)

type Bar interface {
	Done(val int64)

	Finish()
}

type progressBar struct {
	total        int64         // 总进度
	current      int64         // 当前进度
	filler       string        // 进度填充字符
	fillerLength int64         // 进度条长度
	timeFormat   string        // 进度条时间格式
	interval     time.Duration // 打印时间间隔
	begin        time.Time     // 任务开始时间
}

func (p *progressBar) Percent() int64 {
	return p.current * 100 / p.total
}

func (p *progressBar) Done(val int64) {
	p.current += val
}

func (p *progressBar) Finish() {
	fmt.Println(p.ProgressString())
}

// getEta 获取eta时间
func (p *progressBar) getEta(now time.Time) string {
	eta := (now.Unix() - p.begin.Unix()) * 100 / (p.Percent() + 1)
	return p.begin.Add(time.Second * time.Duration(eta)).Format(p.timeFormat)
}

func (p *progressBar) ProgressString() string {
	fills := p.Percent() * p.fillerLength / 100
	chunks := make([]string, p.fillerLength, p.fillerLength)
	for i := int64(0); i < p.fillerLength; i++ {
		switch {
		case i < fills:
			chunks[i] = p.filler
		default:
			chunks[i] = " "
		}
	}
	now := time.Now()
	eta := p.getEta(now)
	qps := p.current / (now.Unix() - p.begin.Unix() + 1)
	return fmt.Sprintf("\r[%s]%d/%d [eta]%s [qps]%d ", strings.Join(chunks, ""), p.current, p.total, eta, qps)
}

// NewProgressBar
func NewProgressBar(total int64, opts ...func(Bar)) Bar {
	bar := &progressBar{
		total:        total,
		filler:       "=",
		fillerLength: 25,
		timeFormat:   "15:04:05", // 2006-01-02T15:04:05
		interval:     time.Second,
		begin:        time.Now(),
	}
	for _, opt := range opts {
		opt(bar)
	}
	// 定时打印
	ticker := time.NewTicker(bar.interval)
	go func() {
		for bar.current < bar.total {
			fmt.Print(bar.ProgressString())
			<-ticker.C
		}
	}()
	return bar
}
