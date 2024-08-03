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

func IncCreatePostAttempts() {
	createPostAttempts.Inc()
}

func IncCreatePostSuccessful() {
	createPostSuccessful.Inc()
}
