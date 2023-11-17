### resilience4go

[![Build](https://github.com/CharLemAznable/gofn/actions/workflows/go.yml/badge.svg)](https://github.com/CharLemAznable/resilience4go/actions/workflows/go.yml)
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

```go
import "github.com/CharLemAznable/resilience4go/bulkhead"

entry := bulkhead.NewBulkhead("name")

decoratedFn := bulkhead.DecorateRunnable(entry, func() error {
	// do something
	return nil
})
```

#### 时长限制(TimeLimiter)

用于限制调用的最大耗时.

```go
import "github.com/CharLemAznable/resilience4go/timelimiter"

entry := timelimiter.NewTimeLimiter("name")

decoratedFn := timelimiter.DecorateRunnable(entry, func() error {
	// do something
	return nil
})
```

#### 速率限制(RateLimiter)

用于限制并发调用的速率.

```go
import "github.com/CharLemAznable/resilience4go/ratelimiter"

entry := ratelimiter.NewRateLimiter("name")

decoratedFn := ratelimiter.DecorateRunnable(entry, func() error {
	// do something
	return nil
})
```

#### 断路器(CircuitBreaker)

对调用进行熔断，避免因持续的失败或拒绝而消耗资源.

```go
import "github.com/CharLemAznable/resilience4go/circuitbreaker"

entry := circuitbreaker.NewCircuitBreaker("name")

decoratedFn := circuitbreaker.DecorateRunnable(entry, func() error {
	// do something
	return nil
})
```

#### 重试(Retry)

在调用失败后, 自动尝试重试.

```go
import "github.com/CharLemAznable/resilience4go/retry"

entry := retry.NewRetry("name")

decoratedFn := retry.DecorateRunnable(entry, func() error {
	// do something
	return nil
})
```

#### 故障恢复(Fallback)

在调用失败后, 根据失败信息进行补偿操作.

```go
import "github.com/CharLemAznable/resilience4go/fallback"

decoratedFn := fallback.DecorateRunnable(func() error {
	// do something
	return nil
}, func(ctx fallback.Context[any, any, error]) error {
	// fallback if has error
	return nil
}, func(ctx fallback.Context[any, any, error]) (bool, fallback.Context[any, any, error]) {
	// 根据调用上下文判断是否需要恢复
})
```

#### 缓存(Cache)

缓存调用结果, 仅支持Function类型的函数包装.

```go
import "github.com/CharLemAznable/resilience4go/cache"

entry := cache.NewCache[keyType, valueType]("name")

decoratedFn := cache.DecorateFunction(entry, func(key keyType) (valueType, error) {
	// do something
	return value, nil
})
```

#### 对如下四种类型的函数进行包装

* Runnable

```go
import "github.com/CharLemAznable/resilience4go/decorator"

decorator.OfRunnable(func() error {
	// ...
})
```

* Supplier

```go
import "github.com/CharLemAznable/resilience4go/decorator"

decorator.OfSupplier[T any](func() (T, error) {
	// ...
})
```

* Consumer

```go
import "github.com/CharLemAznable/resilience4go/decorator"

decorator.OfConsumer[T any](func(t T) error {
	// ...
})
```

* Function

```go
import "github.com/CharLemAznable/resilience4go/decorator"

decorator.OfFunction[T any, R any](func(t T) (R, error) {
	// ...
})
```

#### 使用Prometheus监控弹性组件的指标

```go
import "github.com/CharLemAznable/resilience4go/promhelper"

promhelper.BulkheadRegistry(bulkheadEntry)
promhelper.TimeLimiterRegistry(timelimiterEntry)
promhelper.RateLimiterRegistry(ratelimiterEntry)
promhelper.CircuitBreakerRegistry(circuitbreakerEntry)
promhelper.RetryRegistry(retryEntry)
promhelper.CacheRegistry(cacheEntry)
```

以上方法返回两个函数, 分别为注册到Prometheus的函数和反注册Prometheus的函数.
