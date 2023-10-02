package nsqconsumer

import (
	"time"

	"github.com/nsqio/go-nsq"
)

type PublisherConfig struct {
	*nsq.Config
	NsqdAddress  string `json:"nsqd_address" yaml:"nsqd_address" validate:"nonzero"`
	nsqdHostname string // Stores hostname from parsed NSQd address to avoid parsing on each client refresh.
	nsqdPort     string // Stores port from parsed NSQd address to avoid parsing on each client refresh.
	// ClientRefreshInterval defines how long the NSQd clients should be refreshed.
	ClientRefreshInterval time.Duration `json:"client_refresh_interval" yaml:"client_refresh_interval" default:"10m"`
	// EnableLoadBalancer decides the publisher to use load balancer & auto-refresh clients.
	EnableLoadBalancer bool `json:"enable_load_balancer" yaml:"enable_load_balancer"`
}

type ConsumerConfig struct {
	*nsq.Config
	LookupdAddresses []string `json:"lookupd_addresses" yaml:"lookupd_addresses" validate:"nonzero"`
	Concurrency      int      `json:"concurrency" yaml:"concurrency" validate:"nonzero"`
	MaxInFlight      int      `json:"max_in_flight" yaml:"max_in_flight"`
	MaxAttempts      uint16   `json:"max_attempts" yaml:"max_attempts"`
}
