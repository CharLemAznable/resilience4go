### resilience4go

[![Build](https://github.com/CharLemAznable/resilience4go/actions/workflows/go.yml/badge.svg)](https://github.com/CharLemAznable/resilience4go/actions/workflows/go.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/CharLemAznable/resilience4go)

[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)
![GitHub code size](https://img.shields.io/github/languages/code-size/CharLemAznable/resilience4go)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=alert_status)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)

[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=bugs)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)

[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=security_rating)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)

[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=sqale_index)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=code_smells)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)

[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=ncloc)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=coverage)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=CharLemAznable_resilience4go&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=CharLemAznable_resilience4go)

Golang实现的弹性调用工具包, 参考 [resilience4j](https://github.com/resilience4j/resilience4j) 实现.

#### 舱壁隔离(Bulkhead)

用于限制并发调用的最大次数.

```
// 添加依赖
import "github.com/CharLemAznable/resilience4go/bulkhead"

// 初始化组件
entry := bulkhead.NewBulkhead(string,
	bulkhead.WithMaxConcurrentCalls(int64), // 设置最大并发数量
	bulkhead.WithMaxWaitDuration(time.Duration)) // 设置当舱壁满时的goroutine等待时长

// 可监测指标
metrics := entry.Metrics()
metrics.MaxAllowedConcurrentCalls() // 允许的最大并发数量
metrics.AvailableConcurrentCalls() // 当前可用的并发余量

// 事件监听
listener := entry.EventListener()
listener.OnPermittedFunc(func(PermittedEvent) {
	// goroutine被允许进入并发
})
listener.OnRejectedFunc(func(RejectedEvent) {
	// goroutine被拒绝进入并发
})
listener.OnFinishedFunc(func(FinishedEvent) {
	// goroutine执行完成
})

// 包装执行
bulkhead.DecorateRunnable(entry, func() error {
})
bulkhead.DecorateSupplier(entry, func() (any, error) {
})
bulkhead.DecorateConsumer(entry, func(any) error {
})
bulkhead.DecorateFunction(entry, func(any) (any, error) {
})
```

#### 时长限制(TimeLimiter)

用于限制调用的最大耗时.

```
// 添加依赖
import "github.com/CharLemAznable/resilience4go/timelimiter"

// 初始化组件
entry := timelimiter.NewTimeLimiter(string,
	timelimiter.WithTimeoutDuration(time.Duration)) // 设置执行时限

// 可监测指标
metrics := entry.Metrics()
metrics.SuccessCount() // 执行完成计数
metrics.TimeoutCount() // 执行超时计数
metrics.PanicCount() // 执行发生panic计数

// 事件监听
listener := entry.EventListener()
listener.OnSuccessFunc(func(SuccessEvent) {
	// 执行完成
})
listener.OnTimeoutFunc(func(TimeoutEvent) {
	// 执行超时
})
listener.OnPanicFunc(func(PanicEvent) {
	// 执行发生panic
})

// 包装执行
timelimiter.DecorateRunnable(entry, func() error {
})
timelimiter.DecorateSupplier(entry, func() (any, error) {
})
timelimiter.DecorateConsumer(entry, func(any) error {
})
timelimiter.DecorateFunction(entry, func(any) (any, error) {
})
```

#### 速率限制(RateLimiter)

用于限制并发调用的速率.

```
// 添加依赖
import "github.com/CharLemAznable/resilience4go/ratelimiter"

// 初始化组件
entry := ratelimiter.NewRateLimiter(string,
	ratelimiter.WithTimeoutDuration(time.Duration), // 设置等待并发的时长
	ratelimiter.WithLimitRefreshPeriod(time.Duration), // 设置并发数量的刷新时间
	ratelimiter.WithLimitForPeriod(int64)) // 设置刷新时间内允许的并发数量

// 可监测指标
metrics := entry.Metrics()
metrics.NumberOfWaitingThreads() // 等待中的goroutine数量
metrics.AvailablePermissions() // 剩余可用的并发数量

// 事件监听
listener := entry.EventListener()
listener.OnSuccessFunc(func(SuccessEvent) {
	// goroutine被允许并发
})
listener.OnFailureFunc(func(FailureEvent) {
	// goroutine被限制并发
})

// 包装执行
ratelimiter.DecorateRunnable(entry, func() error {
})
ratelimiter.DecorateSupplier(entry, func() (any, error) {
})
ratelimiter.DecorateConsumer(entry, func(any) error {
})
ratelimiter.DecorateFunction(entry, func(any) (any, error) {
})
```

#### 断路器(CircuitBreaker)

对调用进行熔断，避免因持续的失败或拒绝而消耗资源.

```
// 添加依赖
import "github.com/CharLemAznable/resilience4go/circuitbreaker"

// 初始化组件
entry := circuitbreaker.NewCircuitBreaker(string,
	circuitbreaker.WithSlidingWindow(SlidingWindowType, int64, int64), // 设置滑动窗口类型(时间/调用数量), 大小, 和断路判断的最小调用次数
	circuitbreaker.WithFailureRateThreshold(float64), // 设置断路开启的失败调用率阈值
	circuitbreaker.WithSlowCallRateThreshold(float64), // 设置断路开启的慢调用率阈值
	circuitbreaker.WithSlowCallDurationThreshold(time.Duration), // 设置慢调用判断的时长阈值
	circuitbreaker.WithFailureResultPredicate(func(any, error) bool), // 设置失败调用的判断断言
	circuitbreaker.WithAutomaticTransitionFromOpenToHalfOpenEnabled(bool), // 设置是否自动从断路开启转换为断路半开
	circuitbreaker.WithWaitIntervalFunctionInOpenState(func(int64) time.Duration), // 设置自动从断路开启转换为断路半开的等待时长函数
	circuitbreaker.WithPermittedNumberOfCallsInHalfOpenState(int64), // 设置断路半开时允许通过的调用次数
	circuitbreaker.WithMaxWaitDurationInHalfOpenState(time.Duration)) // 设置断路半开时的最大等待时长

// 可监测指标
metrics := entry.Metrics()
metrics.FailureRate() // 失败调用率
metrics.SlowCallRate() // 慢调用率
metrics.NumberOfCalls() // 调用总量计数
metrics.NumberOfSuccessfulCalls() // 成功调用量计数
metrics.NumberOfFailedCalls() // 失败调用量计数
metrics.NumberOfSlowCalls() // 慢调用总量计数
metrics.NumberOfSlowSuccessfulCalls() // 成功慢调用量计数
metrics.NumberOfSlowFailedCalls() // 失败慢调用量计数
metrics.NumberOfNotPermittedCalls() // 断路调用量计数

// 事件监听
listener := entry.EventListener()
listener.OnSuccessFunc(func(SuccessEvent) {
	// 成功调用
})
listener.OnErrorFunc(func(ErrorEvent) {
	// 失败调用
})
listener.OnNotPermittedFunc(func(NotPermittedEvent) {
	// 断路调用
})
listener.OnStateTransitionFunc(func(StateTransitionEvent) {
	// 断路器状态转换
})
listener.OnFailureRateExceededFunc(func(FailureRateExceededEvent) {
	// 失败调用率到达阈值
})
listener.OnSlowCallRateExceededFunc(func(SlowCallRateExceededEvent) {
	// 慢调用率到达阈值
})

// 包装执行
circuitbreaker.DecorateRunnable(entry, func() error {
})
circuitbreaker.DecorateSupplier(entry, func() (any, error) {
})
circuitbreaker.DecorateConsumer(entry, func(any) error {
})
circuitbreaker.DecorateFunction(entry, func(any) (any, error) {
})
```

#### 重试(Retry)

在调用失败后, 自动尝试重试.

```
// 添加依赖
import "github.com/CharLemAznable/resilience4go/retry"

// 初始化组件
entry := retry.NewRetry(string,
	retry.WithMaxAttempts(int), // 设置最大重试次数
	retry.WithFailAfterMaxAttempts(bool), // 设置是否在最后一次重试失败后返回错误
	retry.WithFailureResultPredicate(func(any, error) bool), // 设置失败调用的判断断言
	retry.WithWaitIntervalFunction(func(int) time.Duration)) // 设置重试的等待时长函数

// 可监测指标
metrics := entry.Metrics()
metrics.NumberOfSuccessfulCallsWithoutRetryAttempt() // 未重试成功调用计数
metrics.NumberOfSuccessfulCallsWithRetryAttempt() // 重试成功调用计数
metrics.NumberOfFailedCallsWithoutRetryAttempt() // 未重试失败调用计数
metrics.NumberOfFailedCallsWithRetryAttempt() // 重试失败调用计数

// 事件监听
listener := entry.EventListener()
listener.OnSuccessFunc(func(SuccessEvent) {
	// 重试成功调用
})
listener.OnRetryFunc(func(RetryEvent) {
	// 即将重试调用
})
listener.OnErrorFunc(func(ErrorEvent) {
	// 失败调用
})

// 包装执行
retry.DecorateRunnable(entry, func() error {
})
retry.DecorateSupplier(entry, func() (any, error) {
})
retry.DecorateConsumer(entry, func(any) error {
})
retry.DecorateFunction(entry, func(any) (any, error) {
})
```

#### 故障恢复(Fallback)

在调用失败后, 根据失败信息进行补偿操作.

```
// 添加依赖
import "github.com/CharLemAznable/resilience4go/fallback"

// 包装执行
fallback.DecorateRunnable(
	func() error {},
	func(ctx fallback.Context[any, any, E]) error {}, // 恢复操作
	func(ctx fallback.Context[any, any, error]) (bool, fallback.Context[any, any, E]) {}) // 根据调用上下文判断是否需要恢复
fallback.DecorateSupplier(
	func() (R, error) {},
	func(ctx fallback.Context[any, R, E]) (R, error) {}, // 恢复操作
	func(ctx fallback.Context[any, R, error]) (bool, fallback.Context[any, R, E]) {}) // 根据调用上下文判断是否需要恢复
fallback.DecorateConsumer(
	func(T) error {},
	func(ctx fallback.Context[T, any, E]) error {}, // 恢复操作
	func(ctx fallback.Context[T, any, error]) (bool, fallback.Context[T, any, E]) {}) // 根据调用上下文判断是否需要恢复
fallback.DecorateFunction(
	func(T) (R, error) {},
	func(ctx fallback.Context[T, R, E]) (R, error) {}, // 恢复操作
	func(ctx fallback.Context[T, R, error]) (bool, fallback.Context[T, R, E]) {}) // 根据调用上下文判断是否需要恢复

// 包装执行, 恢复操作函数接受失败上下文参数, 可限定error类型
fallback.DecorateRunnableWithFailure(
	func() error {},
	func(E) error {}) // 恢复操作
fallback.DecorateSupplierWithFailure(
	func() (R, error) {},
	func(R, E) (R, error) {}) // 恢复操作
fallback.DecorateConsumerWithFailure(
	func(T) error {},
	func(T, E) error {}) // 恢复操作
fallback.DecorateFunctionWithFailure(
	func(T) (R, error) {},
	func(T, R, E) (R, error) {}) // 恢复操作

// 包装执行, 当发生限定类型的error时执行恢复操作函数
fallback.DecorateRunnableByType[E](
	func() error {},
	func() error {}) // 恢复操作
fallback.DecorateSupplierByType[R, E](
	func() (R, error) {},
	func() (R, error) {}) // 恢复操作
fallback.DecorateConsumerByType[T, E](
	func(T) error {},
	func(T) error {}) // 恢复操作
fallback.DecorateFunctionByType[T, R, E](
	func(T) (R, error) {},
	func(T) (R, error) {}) // 恢复操作

// 包装执行, 当发生error时执行恢复操作函数
fallback.DecorateRunnableDefault(
	func() error {},
	func() error {}) // 恢复操作
fallback.DecorateSupplierDefault(
	func() (R, error) {},
	func() (R, error) {}) // 恢复操作
fallback.DecorateConsumerDefault(
	func(T) error {},
	func(T) error {}) // 恢复操作
fallback.DecorateFunctionDefault(
	func(T) (R, error) {},
	func(T) (R, error) {}) // 恢复操作
```

#### 缓存(Cache)

缓存调用结果, 仅支持Function类型的函数包装.

```
// 添加依赖
import "github.com/CharLemAznable/resilience4go/cache"

// 初始化组件
entry := cache.NewCache[K, V](string,
	cache.WithCapacity(int64), // 设置缓存容量
	cache.WithItemTTL(time.Duration), // 设置缓存有效时间
	cache.WithKeyToHash(func(any) (uint64, uint64)), // 设置缓存key的hash函数
	cache.WithCacheResultPredicate(func(any, error) bool)) // 设置是否缓存调用结果的判断断言

// 可选设置缓存值的编解码器
entry = entry.WithMarshalFn(func(V) any, func(any) V)

// 可监测指标
metrics := entry.Metrics()
metrics.NumberOfCacheHits() // 缓存命中计数
metrics.NumberOfCacheMisses() // 缓存未命中计数

// 事件监听
listener := entry.EventListener()
listener.OnCacheHitFunc(func(HitEvent) {
	// 缓存命中
})
listener.OnCacheMissFunc(func(MissEvent) {
	// 缓存未命中
})

// 包装执行
cache.DecorateFunction[K, V](entry, func(K) (V, error) {
})
```

#### 对如下四种类型的函数进行包装

```
import "github.com/CharLemAznable/resilience4go/decorator"

// Runnable: func() error
runnableFn := decorator.
	OfRunnable(func() error {}).
	WithBulkhead(bulkhead.Bulkhead).
	WhenFull(func() error).
	WithTimeLimiter(timelimiter.TimeLimiter).
	WhenTimeout(func() error).
	WithRateLimiter(ratelimiter.RateLimiter).
	WhenOverRate(func() error).
	WithCircuitBreaker(circuitbreaker.CircuitBreaker).
	WhenOverLoad(func() error).
	WithRetry(retry.Retry).
	WhenMaxRetries(func() error).
	WithFallback(func() error, func(error, any) bool).
	Decorate()

// Supplier: func() (any, error)
supplierFn := decorator.
	OfSupplier(func() (any, error) {}).
	WithBulkhead(bulkhead.Bulkhead).
	WhenFull(func() (any, error)).
	WithTimeLimiter(timelimiter.TimeLimiter).
	WhenTimeout(func() (any, error)).
	WithRateLimiter(ratelimiter.RateLimiter).
	WhenOverRate(func() (any, error)).
	WithCircuitBreaker(circuitbreaker.CircuitBreaker).
	WhenOverLoad(func() (any, error)).
	WithRetry(retry.Retry).
	WhenMaxRetries(func() (any, error)).
	WithFallback(func() (any, error), func(any, error, any) bool).
	Decorate()

// Consumer: func(any) error
consumerFn := decorator.
	OfConsumer(func(any) error {}).
	WithBulkhead(bulkhead.Bulkhead).
	WhenFull(func(any) error).
	WithTimeLimiter(timelimiter.TimeLimiter).
	WhenTimeout(func(any) error).
	WithRateLimiter(ratelimiter.RateLimiter).
	WhenOverRate(func(any) error).
	WithCircuitBreaker(circuitbreaker.CircuitBreaker).
	WhenOverLoad(func(any) error).
	WithRetry(retry.Retry).
	WhenMaxRetries(func(any) error).
	WithFallback(func(any) error, func(any, error, any) bool).
	Decorate()

// Function: func(any) (any, error)
functionFn := decorator.
	OfFunction(func() error {}).
	WithBulkhead(bulkhead.Bulkhead).
	WhenFull(func(any) (any, error)).
	WithTimeLimiter(timelimiter.TimeLimiter).
	WhenTimeout(func(any) (any, error)).
	WithRateLimiter(ratelimiter.RateLimiter).
	WhenOverRate(func(any) (any, error)).
	WithCircuitBreaker(circuitbreaker.CircuitBreaker).
	WhenOverLoad(func(any) (any, error)).
	WithRetry(retry.Retry).
	WhenMaxRetries(func(any) (any, error)).
	WithFallback(func(any) (any, error), func(any, any, error, any) bool).
	WithCache(cache.Cache[any, any]).
	Decorate()
```

#### 使用Prometheus监控弹性组件的指标

```
import "github.com/CharLemAznable/resilience4go/promhelper"

promhelper.BulkheadRegistry(bulkheadEntry)
promhelper.TimeLimiterRegistry(timelimiterEntry)
promhelper.RateLimiterRegistry(ratelimiterEntry)
promhelper.CircuitBreakerRegistry(circuitbreakerEntry)
promhelper.RetryRegistry(retryEntry)
promhelper.CacheRegistry(cacheEntry)
```

以上方法返回两个函数, 分别为注册到Prometheus的函数和反注册Prometheus的函数.

注册的指标如下:

```
// bulkhead
gauge: 
  name:  resilience4go_bulkhead_max_allowed_concurrent_calls
  label: {name: bulkhead-name}
gauge:
  name:  resilience4go_bulkhead_available_concurrent_calls
  label: {name: bulkhead-name}

// timelimiter
counter:
  name:  resilience4go_timelimiter_calls
  label: {name: timelimiter-name}, {kind: successful}
counter:
  name:  resilience4go_timelimiter_calls
  label: {name: timelimiter-name}, {kind: timeout}
counter:
  name:  resilience4go_timelimiter_calls
  label: {name: timelimiter-name}, {kind: panicked}

// ratelimiter
gauge:
  name:  resilience4go_ratelimiter_waiting_threads
  label: {name: ratelimiter-name}
gauge:
  name:  resilience4go_ratelimiter_available_permissions
  label: {name: ratelimiter-name}

// circuitbreaker
gauge:
  name:  resilience4go_circuitbreaker_state
  label: {name: circuitbreaker-name}, {state: closed}
gauge:
  name:  resilience4go_circuitbreaker_state
  label: {name: circuitbreaker-name}, {state: open}
gauge:
  name:  resilience4go_circuitbreaker_state
  label: {name: circuitbreaker-name}, {state: half_open}
gauge:
  name:  resilience4go_circuitbreaker_state
  label: {name: circuitbreaker-name}, {state: disabled}
gauge:
  name:  resilience4go_circuitbreaker_state
  label: {name: circuitbreaker-name}, {state: forced_open}
gauge:
  name:  resilience4go_circuitbreaker_buffered_calls
  label: {name: circuitbreaker-name}, {kind: successful}
gauge:
  name:  resilience4go_circuitbreaker_buffered_calls
  label: {name: circuitbreaker-name}, {kind: failed}
gauge:
  name:  resilience4go_circuitbreaker_slow_calls
  label: {name: circuitbreaker-name}, {kind: successful}
gauge:
  name:  resilience4go_circuitbreaker_slow_calls
  label: {name: circuitbreaker-name}, {kind: failed}
gauge:
  name:  resilience4go_circuitbreaker_failure_rate
  label: {name: circuitbreaker-name}
gauge:
  name:  resilience4go_circuitbreaker_slow_call_rate
  label: {name: circuitbreaker-name}
histogram:
  name:  resilience4go_circuitbreaker_calls
  label: {name: circuitbreaker-name}, {kind: successful}
histogram:
  name:  resilience4go_circuitbreaker_calls
  label: {name: circuitbreaker-name}, {kind: failed}
counter:
  name:  resilience4go_circuitbreaker_not_permitted_calls
  label: {name: circuitbreaker-name}, {kind: not_permitted}

// retry
counter:
  name:  resilience4go_retry_calls
  label: {name: retry-name}, {kind: successful_without_retry}
counter:
  name:  resilience4go_retry_calls
  label: {name: retry-name}, {kind: successful_with_retry}
counter:
  name:  resilience4go_retry_calls
  label: {name: retry-name}, {kind: failed_without_retry}
counter:
  name:  resilience4go_retry_calls
  label: {name: retry-name}, {kind: failed_with_retry}

// cache
gauge: 
  name:  resilience4go_cache_hits
  label: {name: cache-name}
gauge:
  name:  resilience4go_cache_misses
  label: {name: cache-name}
```
