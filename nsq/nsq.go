package nsqs

import (
	"sync"
	"sync/atomic"

	"github.com/nsqio/go-nsq"
	"github.com/pyihe/go-pkg/errors"
)

var (
	m = mgr{
		consumers: map[string]map[string][]*nsq.Consumer{},
		producers: map[string]*nsq.Producer{},
	}
)

type mgr struct {
	config       *nsq.Config
	lookupdAddrs []string

	cMu       sync.RWMutex
	consumers map[string]map[string][]*nsq.Consumer

	pMu       sync.RWMutex
	producers map[string]*nsq.Producer

	stopTag int32
}

func NewConsumerMgr(config *nsq.Config, lookupdAddrs []string) {
	m.config = config
	m.lookupdAddrs = lookupdAddrs
}

func isStop() bool {
	return atomic.LoadInt32(&m.stopTag) == 1
}

func NewConfig() *nsq.Config {
	return nsq.NewConfig()
}

func AddConsumer(topic, channel string, handler nsq.Handler) error {
	if isStop() {
		return nil
	}
	c, err := nsq.NewConsumer(topic, channel, m.config)
	if err != nil {
		return err
	}
	c.AddHandler(handler)

	m.cMu.Lock()
	if m.consumers[topic] == nil {
		m.consumers[topic] = make(map[string][]*nsq.Consumer)
	}
	m.consumers[topic][channel] = append(m.consumers[topic][channel], c)
	m.cMu.Unlock()

	return c.ConnectToNSQLookupds(m.lookupdAddrs)
}

func AddHandler(topic, channel string, handler nsq.Handler) error {
	if isStop() {
		return nil
	}
	m.cMu.RLock()
	consumers := m.consumers[topic][channel]
	m.cMu.RUnlock()
	if len(consumers) == 0 {
		return errors.New("not exist consumer")
	}
	for _, c := range consumers {
		c.AddHandler(handler)
	}
	return nil
}

func AddConcurrentHandlers(topic, channel string, handler nsq.Handler, concurrency int) error {
	if isStop() {
		return nil
	}
	m.cMu.RLock()
	consumers := m.consumers[topic][channel]
	m.cMu.RUnlock()
	if len(consumers) == 0 {
		return errors.New("not exist consumer")
	}
	for _, c := range consumers {
		c.AddConcurrentHandlers(handler, concurrency)
	}
	return nil
}

func GetConsumer(topic, channel string) (consumers []*nsq.Consumer) {
	if isStop() {
		return nil
	}
	m.cMu.RLock()
	consumers = m.consumers[topic][channel]
	m.cMu.RUnlock()
	return
}

func Stop() {
	if isStop() {
		return
	}
	atomic.StoreInt32(&m.stopTag, 1)
	m.cMu.RLock()
	for _, channels := range m.consumers {
		for _, cs := range channels {
			for _, c := range cs {
				c.Stop()
			}
		}
	}
	m.cMu.RUnlock()

	m.pMu.RLock()
	for _, p := range m.producers {
		p.Stop()
	}
	m.pMu.RUnlock()
}

func AddProducer(addr string, config *nsq.Config) (*nsq.Producer, error) {
	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		return nil, err
	}
	m.pMu.Lock()
	m.producers[addr] = p
	m.pMu.Unlock()
	return p, nil
}
