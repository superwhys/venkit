package vhttp

import "time"

type Config struct {
	RequestTimeOut time.Duration
	Proxy          string
}
