package entity

type PollResult struct {
	RetryCount    int
	IsHealthy     bool
	PolledService Service
}
