package quotas

import (
	"context"
	"fmt"
	"sort"
	"time"
)

type (
	// PriorityRateLimiterImpl is a wrapper around the golang rate limiter
	PriorityRateLimiterImpl struct {
		requestPriorityFn      RequestPriorityFn
		priorityToRateLimiters map[int]RequestRateLimiter

		// priority value 0 means highest priority
		// sorted rate limiter from low priority value to high priority value
		priorityToIndex map[int]int
		rateLimiters    []RequestRateLimiter
	}
)

var _ RequestRateLimiter = (*PriorityRateLimiterImpl)(nil)

// NewPriorityRateLimiter returns a new rate limiter that can handle dynamic
// configuration updates
func NewPriorityRateLimiter(
	requestPriorityFn RequestPriorityFn,
	priorityToRateLimiters map[int]RequestRateLimiter,
) *PriorityRateLimiterImpl {
	priorities := make([]int, 0, len(priorityToRateLimiters))
	for priority := range priorityToRateLimiters {
		priorities = append(priorities, priority)
	}
	sort.Slice(priorities, func(i, j int) bool {
		return priorities[i] < priorities[j]
	})
	priorityToIndex := make(map[int]int, len(priorityToRateLimiters))
	rateLimiters := make([]RequestRateLimiter, 0, len(priorityToRateLimiters))
	for index, priority := range priorities {
		priorityToIndex[priority] = index
		rateLimiters = append(rateLimiters, priorityToRateLimiters[priority])
	}

	return &PriorityRateLimiterImpl{
		requestPriorityFn:      requestPriorityFn,
		priorityToRateLimiters: priorityToRateLimiters,

		priorityToIndex: priorityToIndex,
		rateLimiters:    rateLimiters,
	}
}

func (p *PriorityRateLimiterImpl) Allow(
	now time.Time,
	request Request,
) bool {
	decidingRateLimiter, consumeRateLimiters := p.getRateLimiters(request)

	allow := decidingRateLimiter.Allow(now, request)
	if !allow {
		return false
	}

	for _, limiter := range consumeRateLimiters {
		_ = limiter.Reserve(now, request)
	}
	return allow
}

func (p *PriorityRateLimiterImpl) Reserve(
	now time.Time,
	request Request,
) Reservation {
	decidingRateLimiter, consumeRateLimiters := p.getRateLimiters(request)

	decidingReservation := decidingRateLimiter.Reserve(now, request)
	if !decidingReservation.OK() {
		return decidingReservation
	}

	otherReservations := make([]Reservation, len(consumeRateLimiters))
	for index, limiter := range consumeRateLimiters {
		otherReservations[index] = limiter.Reserve(now, request)
	}
	return NewPriorityReservation(decidingReservation, otherReservations)
}

func (p *PriorityRateLimiterImpl) Wait(
	ctx context.Context,
	request Request,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	now := time.Now().UTC()
	reservation := p.Reserve(now, request)
	if !reservation.OK() {
		return fmt.Errorf("rate: Wait(n=%d) would exceed context deadline", request.Token)
	}

	delay := reservation.DelayFrom(now)
	if delay == 0 {
		return nil
	}
	waitLimit := InfDuration
	if deadline, ok := ctx.Deadline(); ok {
		waitLimit = deadline.Sub(now)
	}
	if waitLimit < delay {
		reservation.CancelAt(now)
		return fmt.Errorf("rate: Wait(n=%d) would exceed context deadline", request.Token)
	}

	t := time.NewTimer(delay)
	defer t.Stop()
	select {
	case <-t.C:
		return nil

	case <-ctx.Done():
		reservation.CancelAt(time.Now())
		return ctx.Err()
	}
}

func (p *PriorityRateLimiterImpl) getRateLimiters(
	request Request,
) (RequestRateLimiter, []RequestRateLimiter) {
	priority := p.requestPriorityFn(request)
	if _, ok := p.priorityToRateLimiters[priority]; !ok {
		panic("Request to priority & priority to rate limiter does not match")
	}

	rateLimiterIndex := p.priorityToIndex[priority]
	return p.rateLimiters[rateLimiterIndex], p.rateLimiters[rateLimiterIndex+1:]
}
