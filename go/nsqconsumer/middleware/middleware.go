package middleware

import "github.com/nsqio/go-nsq"

type Func func(topic, channel string, next nsq.Handler) nsq.Handler
