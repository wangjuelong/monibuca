package engine

import (
	"context"
	"time"
)

// Publisher 发布者实体定义
type Publisher struct {
	context.Context
	cancel        context.CancelFunc
	AutoUnPublish bool //	当无人订阅时自动停止发布
	*Stream       `json:"-"`
	Type          string      //类型，用来区分不同的发布者
	timeout       *time.Timer //更新时间用来做超时处理
}

// Close 关闭发布者
func (p *Publisher) Close() {
	if p.Running() {
		p.Stream.Close()
	}
}

func (p *Publisher) Update() {
	if p.timeout != nil {
		p.timeout.Reset(config.PublishTimeout)
	} else {
		p.timeout = time.AfterFunc(config.PublishTimeout, p.Close)
	}
}

// Running 发布者是否正在发布
func (p *Publisher) Running() bool {
	return p.Stream != nil && p.Err() == nil
}

// Publish 发布者进行发布操作
func (p *Publisher) Publish(streamPath string) bool {
	p.Stream = GetStream(streamPath)
	//检查是否已存在发布者
	if p.Publisher != nil {
		return false
	}
	p.Context, p.cancel = context.WithCancel(p.Stream)
	p.Publisher = p
	p.StartTime = time.Now()
	p.Update()
	//触发钩子
	TriggerHook(Hook{HOOK_PUBLISH, p.Stream})
	return true
}
