package elasticp

import "time"

type ElasticPoolOption func(*ElasticPoolConfiguration)

type ElasticPoolConfiguration struct {
	workBuffer    int
	workThreshold float32
	growthRatio   int

	minWorkers    int
	maxWorkers    int
	idleTolerance time.Duration
	isDebug       bool
}

func NewConfig(options ...ElasticPoolOption) ElasticPoolConfiguration {
	config := &ElasticPoolConfiguration{
		workBuffer:    10,
		workThreshold: 0.5,
		growthRatio:   2,
		minWorkers:    5,
		maxWorkers:    20,
		idleTolerance: time.Second,
		isDebug:       true,
	}

	for _, option := range options {
		option(config)
	}

	return *config
}

func WithWorkBuffer(buffer int) ElasticPoolOption {
	return func(config *ElasticPoolConfiguration) {
		config.workBuffer = buffer
	}
}

func WithWorkThreshold(threshold float32) ElasticPoolOption {
	return func(config *ElasticPoolConfiguration) {
		config.workThreshold = threshold
	}
}

func WithGrowthRatio(ratio int) ElasticPoolOption {
	return func(config *ElasticPoolConfiguration) {
		config.growthRatio = ratio
	}
}

func WithMinWorkers(minWorkers int) ElasticPoolOption {
	return func(config *ElasticPoolConfiguration) {
		config.minWorkers = minWorkers
	}
}

func WithMaxWorkers(maxWorkers int) ElasticPoolOption {
	return func(config *ElasticPoolConfiguration) {
		config.maxWorkers = maxWorkers
	}
}

func WithIdleTolerance(idleTolerance time.Duration) ElasticPoolOption {
	return func(config *ElasticPoolConfiguration) {
		config.idleTolerance = idleTolerance
	}
}

func WithDebug(debug bool) ElasticPoolOption {
	return func(config *ElasticPoolConfiguration) {
		config.isDebug = debug
	}
}
