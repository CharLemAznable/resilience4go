package circuitbreaker

import (
	"github.com/CharLemAznable/resilience4go/utils"
	"sync"
	"sync/atomic"
	"time"
)

type Metrics interface {
	FailureRate() float64
	SlowCallRate() float64
	NumberOfCalls() int64
	NumberOfSuccessfulCalls() int64
	NumberOfFailedCalls() int64
	NumberOfSlowCalls() int64
	NumberOfSlowSuccessfulCalls() int64
	NumberOfSlowFailedCalls() int64
	NumberOfNotPermittedCalls() int64

	onCallNotPermitted()
	onSuccess(time.Duration) metricsResult
	onError(time.Duration) metricsResult
}

func newMetrics(slidingWindowSize int64, slidingWindowType SlidingWindowType, config *Config) Metrics {
	m := &metrics{
		failureRateThreshold:      config.failureRateThreshold,
		slowCallRateThreshold:     config.slowCallRateThreshold,
		slowCallDurationThreshold: config.slowCallDurationThreshold,
	}
	if CountBased == slidingWindowType {
		m.recorder = newFixedSizeSlidingWindowRecorder(slidingWindowSize)
		m.minimumNumberOfCalls = utils.Min(config.minimumNumberOfCalls, slidingWindowSize)
	} else {
		m.recorder = newSlidingTimeWindowRecorder(slidingWindowSize)
		m.minimumNumberOfCalls = config.minimumNumberOfCalls
	}
	return m
}

func forClosed(config *Config) Metrics {
	return newMetrics(config.slidingWindowSize, config.slidingWindowType, config)
}

func forHalfOpen(permittedNumberOfCallsInHalfOpenState int64, config *Config) Metrics {
	return newMetrics(permittedNumberOfCallsInHalfOpenState, CountBased, config)
}

func forDisabled(config *Config) Metrics {
	return newMetrics(0, CountBased, config)
}

func forForcedOpen(config *Config) Metrics {
	return newMetrics(0, CountBased, config)
}

type metricsResult string

const (
	belowThresholds             metricsResult = "BELOW_THRESHOLDS"
	failureRateAboveThresholds  metricsResult = "FAILURE_RATE_ABOVE_THRESHOLDS"
	slowCallRateAboveThresholds metricsResult = "SLOW_CALL_RATE_ABOVE_THRESHOLDS"
	aboveThresholds             metricsResult = "ABOVE_THRESHOLDS"
	belowMinimumCallsThreshold  metricsResult = "BELOW_MINIMUM_CALLS_THRESHOLD"
)

func failureRateExceededThreshold(result metricsResult) bool {
	return result == aboveThresholds || result == failureRateAboveThresholds
}

func slowCallRateExceededThreshold(result metricsResult) bool {
	return result == aboveThresholds || result == slowCallRateAboveThresholds
}

func exceededThresholds(result metricsResult) bool {
	return failureRateExceededThreshold(result) || slowCallRateExceededThreshold(result)
}

type metrics struct {
	recorder                  recorder
	minimumNumberOfCalls      int64
	failureRateThreshold      float64
	slowCallRateThreshold     float64
	slowCallDurationThreshold time.Duration
	numberOfNotPermittedCalls atomic.Int64
}

func (m *metrics) FailureRate() float64 {
	return m.failureRate(m.recorder.snapshot())
}

func (m *metrics) SlowCallRate() float64 {
	return m.slowCallRate(m.recorder.snapshot())
}

func (m *metrics) NumberOfCalls() int64 {
	return m.recorder.snapshot().totalNumberOfCalls
}

func (m *metrics) NumberOfSuccessfulCalls() int64 {
	return m.recorder.snapshot().totalNumberOfSuccessfulCalls()
}

func (m *metrics) NumberOfFailedCalls() int64 {
	return m.recorder.snapshot().totalNumberOfFailedCalls
}

func (m *metrics) NumberOfSlowCalls() int64 {
	return m.recorder.snapshot().totalNumberOfSlowCalls
}

func (m *metrics) NumberOfSlowSuccessfulCalls() int64 {
	return m.recorder.snapshot().totalNumberOfSlowSuccessfulCalls()
}

func (m *metrics) NumberOfSlowFailedCalls() int64 {
	return m.recorder.snapshot().totalNumberOfSlowFailedCalls
}

func (m *metrics) NumberOfNotPermittedCalls() int64 {
	return m.numberOfNotPermittedCalls.Load()
}

func (m *metrics) onCallNotPermitted() {
	m.numberOfNotPermittedCalls.Add(1)
}

func (m *metrics) onSuccess(duration time.Duration) metricsResult {
	calcOutcome := m.calcOutcome(duration, slowSuccessOutcome, successOutcome)
	snap := m.recorder.record(duration, calcOutcome)
	return m.checkIfThresholdsExceeded(snap)
}

func (m *metrics) onError(duration time.Duration) metricsResult {
	calcOutcome := m.calcOutcome(duration, slowErrorOutcome, errorOutcome)
	snap := m.recorder.record(duration, calcOutcome)
	return m.checkIfThresholdsExceeded(snap)
}

func (m *metrics) calcOutcome(duration time.Duration, slow, normal outcome) outcome {
	if duration > m.slowCallDurationThreshold {
		return slow
	}
	return normal
}

func (m *metrics) checkIfThresholdsExceeded(snap *snapshot) metricsResult {
	failureRateInPercentage := m.failureRate(snap)
	slowCallsInPercentage := m.slowCallRate(snap)

	if failureRateInPercentage == -1 || slowCallsInPercentage == -1 {
		return belowMinimumCallsThreshold
	}
	if failureRateInPercentage >= m.failureRateThreshold &&
		slowCallsInPercentage >= m.slowCallRateThreshold {
		return aboveThresholds
	}
	if failureRateInPercentage >= m.failureRateThreshold {
		return failureRateAboveThresholds
	}
	if slowCallsInPercentage >= m.slowCallRateThreshold {
		return slowCallRateAboveThresholds
	}
	return belowThresholds
}

func (m *metrics) failureRate(snap *snapshot) float64 {
	bufferedCalls := snap.totalNumberOfCalls
	if bufferedCalls == 0 || bufferedCalls < m.minimumNumberOfCalls {
		return -1.0
	}
	return snap.failureRate()
}

func (m *metrics) slowCallRate(snap *snapshot) float64 {
	bufferedCalls := snap.totalNumberOfCalls
	if bufferedCalls == 0 || bufferedCalls < m.minimumNumberOfCalls {
		return -1.0
	}
	return snap.slowCallRate()
}

type outcome string

const (
	successOutcome     outcome = "SUCCESS"
	errorOutcome       outcome = "ERROR"
	slowSuccessOutcome outcome = "SLOW_SUCCESS"
	slowErrorOutcome   outcome = "SLOW_ERROR"
)

type recorder interface {
	record(duration time.Duration, outcome outcome) *snapshot
	snapshot() *snapshot
}

func newFixedSizeSlidingWindowRecorder(windowSize int64) recorder {
	return &fixedSizeSlidingWindowRecorder{
		windowSize:       windowSize,
		totalAggregation: totalAggregation{},
		measurements:     make([]measurement, windowSize),
		headIndex:        0,
	}
}

func newSlidingTimeWindowRecorder(timeWindowSizeInSeconds int64) recorder {
	s := &slidingTimeWindowRecorder{
		timeWindowSizeInSeconds: timeWindowSizeInSeconds,
		totalAggregation:        totalAggregation{},
		partialAggregations:     make([]partialAggregation, timeWindowSizeInSeconds),
		headIndex:               0,
	}
	epochSecond := time.Now().Unix()
	for _, p := range s.partialAggregations {
		p.epochSecond = epochSecond
		epochSecond++
	}
	return s
}

type fixedSizeSlidingWindowRecorder struct {
	mutex            sync.Mutex
	windowSize       int64
	totalAggregation totalAggregation
	measurements     []measurement
	headIndex        int64
}

func (r *fixedSizeSlidingWindowRecorder) record(duration time.Duration, outcome outcome) *snapshot {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.totalAggregation.record(duration, outcome)
	r.moveWindowByOne().record(duration, outcome)
	return newSnapshot(&r.totalAggregation)
}

func (r *fixedSizeSlidingWindowRecorder) snapshot() *snapshot {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return newSnapshot(&r.totalAggregation)
}

func (r *fixedSizeSlidingWindowRecorder) moveWindowByOne() *measurement {
	r.moveHeadIndexByOne()
	latestMeasurement := r.latestMeasurement()
	r.totalAggregation.removeBucket(&latestMeasurement.aggregation)
	latestMeasurement.reset()
	return latestMeasurement
}

func (r *fixedSizeSlidingWindowRecorder) moveHeadIndexByOne() {
	r.headIndex = (r.headIndex + 1) % r.windowSize
}

func (r *fixedSizeSlidingWindowRecorder) latestMeasurement() *measurement {
	return &r.measurements[r.headIndex]
}

type slidingTimeWindowRecorder struct {
	mutex                   sync.Mutex
	timeWindowSizeInSeconds int64
	totalAggregation        totalAggregation
	partialAggregations     []partialAggregation
	headIndex               int64
}

func (s *slidingTimeWindowRecorder) record(duration time.Duration, outcome outcome) *snapshot {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.totalAggregation.record(duration, outcome)
	s.moveWindowToCurrentEpochSecond(s.latestPartialAggregation()).record(duration, outcome)
	return newSnapshot(&s.totalAggregation)
}

func (s *slidingTimeWindowRecorder) snapshot() *snapshot {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.moveWindowToCurrentEpochSecond(s.latestPartialAggregation())
	return newSnapshot(&s.totalAggregation)
}

func (s *slidingTimeWindowRecorder) moveWindowToCurrentEpochSecond(latestPartialAggregation *partialAggregation) *partialAggregation {
	currentEpochSecond := time.Now().Unix()
	differenceInSeconds := currentEpochSecond - latestPartialAggregation.epochSecond
	if differenceInSeconds == 0 {
		return latestPartialAggregation
	}
	secondsToMoveTheWindow := utils.Min(differenceInSeconds, s.timeWindowSizeInSeconds)
	var currentPartialAggregation *partialAggregation
	for secondsToMoveTheWindow > 0 {
		secondsToMoveTheWindow--
		s.moveHeadIndexByOne()
		currentPartialAggregation = s.latestPartialAggregation()
		s.totalAggregation.removeBucket(&currentPartialAggregation.aggregation)
		currentPartialAggregation.reset(currentEpochSecond - secondsToMoveTheWindow)
	}
	return currentPartialAggregation
}

func (s *slidingTimeWindowRecorder) moveHeadIndexByOne() {
	s.headIndex = (s.headIndex + 1) % s.timeWindowSizeInSeconds
}

func (s *slidingTimeWindowRecorder) latestPartialAggregation() *partialAggregation {
	return &s.partialAggregations[s.headIndex]
}

type snapshot struct {
	totalDuration                time.Duration
	totalNumberOfCalls           int64
	totalNumberOfFailedCalls     int64
	totalNumberOfSlowCalls       int64
	totalNumberOfSlowFailedCalls int64
}

func newSnapshot(total *totalAggregation) *snapshot {
	return &snapshot{
		totalDuration:                total.totalDuration,
		totalNumberOfCalls:           total.numberOfCalls,
		totalNumberOfFailedCalls:     total.numberOfFailedCalls,
		totalNumberOfSlowCalls:       total.numberOfSlowCalls,
		totalNumberOfSlowFailedCalls: total.numberOfSlowFailedCalls,
	}
}

func (s *snapshot) totalNumberOfSuccessfulCalls() int64 {
	return s.totalNumberOfCalls - s.totalNumberOfFailedCalls
}

func (s *snapshot) totalNumberOfSlowSuccessfulCalls() int64 {
	return s.totalNumberOfSlowCalls - s.totalNumberOfSlowFailedCalls
}

func (s *snapshot) failureRate() float64 {
	if s.totalNumberOfCalls == 0 {
		return 0
	}
	return float64(s.totalNumberOfFailedCalls) * 100.0 / float64(s.totalNumberOfCalls)
}

func (s *snapshot) slowCallRate() float64 {
	if s.totalNumberOfCalls == 0 {
		return 0
	}
	return float64(s.totalNumberOfSlowCalls) * 100.0 / float64(s.totalNumberOfCalls)
}

type aggregation struct {
	totalDuration           time.Duration
	numberOfCalls           int64
	numberOfFailedCalls     int64
	numberOfSlowCalls       int64
	numberOfSlowFailedCalls int64
}

func (agg *aggregation) record(duration time.Duration, outcome outcome) {
	agg.numberOfCalls++
	agg.totalDuration += duration
	switch outcome {
	case slowSuccessOutcome:
		agg.numberOfSlowCalls++
	case slowErrorOutcome:
		agg.numberOfSlowCalls++
		agg.numberOfFailedCalls++
		agg.numberOfSlowFailedCalls++
	case errorOutcome:
		agg.numberOfFailedCalls++
	}
}

type totalAggregation struct {
	aggregation
}

func (t *totalAggregation) removeBucket(bucket *aggregation) {
	t.totalDuration -= bucket.totalDuration
	t.numberOfCalls -= bucket.numberOfCalls
	t.numberOfFailedCalls -= bucket.numberOfFailedCalls
	t.numberOfSlowCalls -= bucket.numberOfSlowCalls
	t.numberOfSlowFailedCalls -= bucket.numberOfSlowFailedCalls
}

type measurement struct {
	aggregation
}

func (m *measurement) reset() {
	m.totalDuration = 0
	m.numberOfCalls = 0
	m.numberOfFailedCalls = 0
	m.numberOfSlowCalls = 0
	m.numberOfSlowFailedCalls = 0
}

type partialAggregation struct {
	aggregation
	epochSecond int64
}

func (p *partialAggregation) reset(epochSecond int64) {
	p.epochSecond = epochSecond
	p.totalDuration = 0
	p.numberOfCalls = 0
	p.numberOfFailedCalls = 0
	p.numberOfSlowCalls = 0
	p.numberOfSlowFailedCalls = 0
}
