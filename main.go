package main

import (
	"math/rand"
	"time"

	"github.com/yosisa/fluxion/buffer"
	"github.com/yosisa/fluxion/message"
	"github.com/yosisa/fluxion/plugin"
)

type Config struct {
	Tag        string
	Interval   buffer.Duration
	RandFactor float64 `toml:"rand_factor"`
}

type InIpsum struct {
	env      *plugin.Env
	conf     Config
	interval time.Duration
	chain    *Chain
	stop     chan struct{}
}

func (p *InIpsum) Init(env *plugin.Env) (err error) {
	p.env = env
	if err = env.ReadConfig(&p.conf); err != nil {
		return
	}
	if p.conf.RandFactor == 0 {
		p.conf.RandFactor = 0.5
	}
	p.interval = time.Duration(p.conf.Interval)
	p.chain = NewChain(LoremIpsum)
	p.stop = make(chan struct{})
	return
}

func (p *InIpsum) Start() error {
	go p.loop()
	return nil
}

func (p *InIpsum) Close() error {
	close(p.stop)
	return nil
}

func (p *InIpsum) loop() {
	for {
		m := map[string]interface{}{"message": p.chain.Generate(100)}
		p.env.Emit(message.NewEvent(p.conf.Tag, m))
		delta := float64(p.interval) * p.conf.RandFactor
		wait := float64(p.interval) - delta + delta*2*rand.Float64()
		select {
		case <-time.After(time.Duration(wait)):
		case <-p.stop:
			return
		}
	}
}

func main() {
	plugin.New("in-ipsum", func() plugin.Plugin { return &InIpsum{} }).Run()
}
