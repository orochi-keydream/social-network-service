package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var createPostAttempts = promauto.NewCounter(prometheus.CounterOpts{
	Name: "create_post_attempts",
	Help: "Number of attempts to create post",
})

var createPostSuccessful = promauto.NewCounter(prometheus.CounterOpts{
	Name: "create_post_successful",
	Help: "Number of created posts",
})

var getMessagesErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "read_messages_errors",
	Help: "Number of errors when reading messages from database",
})

var addMessageErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "send_message_errors",
	Help: "Number of errors when writing messages to database",
})

func IncCreatePostAttempts() {
	createPostAttempts.Inc()
}

func IncCreatePostSuccessful() {
	createPostSuccessful.Inc()
}

func IncGetMessagesErrors() {
	getMessagesErrors.Inc()
}

func IncAddMessageErrors() {
	addMessageErrors.Inc()
}
