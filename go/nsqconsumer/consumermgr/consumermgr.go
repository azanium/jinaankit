package consumermgr

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/azanium/jinaankit/go/nsqconsumer"
	"github.com/azanium/jinaankit/go/nsqconsumer/middleware"
	"github.com/nsqio/go-nsq"
)

type ConsumerMgr struct {
	stopTimeout time.Duration
	middlewares []middleware.Func
	consumers   []*nsq.Consumer
}

type Config struct {
	StopTimeout int
}

func NewConsumerMgr(cfg Config, middlewares ...middleware.Func) *ConsumerMgr {
	if cfg.StopTimeout <= 0 {
		cfg.StopTimeout = 30
	}

	return &ConsumerMgr{
		stopTimeout: time.Duration(cfg.StopTimeout) * time.Second,
		middlewares: middlewares,
	}
}

func (cm *ConsumerMgr) AddConsumer(consumer *nsq.Consumer) {
	cm.consumers = append(cm.consumers, consumer)
}

func (cm *ConsumerMgr) AddHandlerFunc(topic, channel string, cfg nsqconsumer.ConsumerConfig, handlerFn nsq.HandlerFunc) error {
	return cm.AddHandler(topic, channel, cfg, handlerFn)
}

func (cm *ConsumerMgr) AddHandler(topic, channel string, cfg nsqconsumer.ConsumerConfig, handler nsq.Handler) error {
	if cfg.Config == nil {
		cfg.Config = nsq.NewConfig()
	}

	if cfg.MaxAttempts > 0 {
		cfg.Config.MaxAttempts = cfg.MaxAttempts
	}

	if cfg.MaxInFlight > 0 {
		cfg.Config.MaxInFlight = cfg.MaxInFlight
	}

	con, err := nsq.NewConsumer(topic, channel, cfg.Config)
	if err != nil {
		return err
	}

	for _, mw := range cm.middlewares {
		handler = mw(topic, channel, handler)
	}

	log.Println("Config LookupAddresses: ", cfg.LookupdAddresses)
	con.AddConcurrentHandlers(handler, cfg.Concurrency)
	err = con.ConnectToNSQLookupds(cfg.LookupdAddresses)
	if err != nil {
		return err
	}

	cm.consumers = append(cm.consumers, con)
	return nil
}

func (cm *ConsumerMgr) Wait() {
	<-WaitTermSig(cm.Stop)
}

func WaitTermSig(handler func(context.Context) error) <-chan struct{} {
	stoppedCh := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)

		// wait for the sigterm
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-signals

		// We received an os signal, shut down.
		if err := handler(context.Background()); err != nil {
			log.Printf("graceful shutdown  failed: %v", err)
		} else {
			log.Println("graceful shutdown succeed")
		}

		close(stoppedCh)

	}()
	return stoppedCh
}

func (cm *ConsumerMgr) Stop(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, con := range cm.consumers {
		wg.Add(1)
		con := con
		go func() { // use goroutines to stop all of them ASAP
			defer wg.Done()
			con.Stop()

			select {
			case <-con.StopChan:
			case <-ctx.Done():
			case <-time.After(cm.stopTimeout):
			}
		}()
	}
	wg.Wait()
	return nil
}
