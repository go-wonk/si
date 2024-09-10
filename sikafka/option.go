package sikafka

func WithSyncProducerOptionRetyMax(retryMax uint16) SyncProducerOptionFunc {
	return SyncProducerOptionFunc(func(o *SyncProducer) error {
		o.retryMax = retryMax
		return nil
	})
}

type SyncProducerOption interface {
	apply(o *SyncProducer) error
}

type SyncProducerOptionFunc func(o *SyncProducer) error

func (s SyncProducerOptionFunc) apply(o *SyncProducer) error {
	return s(o)
}
