// Package pubsub implements pub/sub interfaces defined in package kolide.
package pubsub

// NoSubscriber is an interface implemented by errors in the pubsub package
// that allows handlers to determine whether the error is due to no subscribers
type NoSubscriber interface {
	// NoSubscriber returns true if the error occurred because there are no
	// subscribers on the channel
	NoSubscriber() bool
}

// NoSubscriberError can be returned when channel operations fail because there
// are no subscribers. Its NoSubscriber() method always returns true.
type NoSubscriberError struct {
	Channel string
}

func (e NoSubscriberError) Error() string {
	return "no subscriber for channel " + e.Channel
}

func (e NoSubscriberError) NoSubscriber() bool {
	return true
}
