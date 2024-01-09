package fallback

import . "github.com/CharLemAznable/gogo/fn"

func DecorateConsumer[T any, E error](
	consumer Consumer[T],
	fallback func(Context[T, any, E]) error,
	predicate FailurePredicate[T, any, E]) Consumer[T] {
	return ConsumerCast(func(t T) error {
		ctx := execute[T, any](t, func() (any, error) {
			return nil, consumer.CheckedAccept(t)
		})
		if ok, failCtx := predicate(ctx); ok {
			return fallback(failCtx)
		}
		return ctx.returnError()
	})
}

func DecorateConsumerWithFailure[T any, E error](
	consumer Consumer[T], fallback func(T, E) error) Consumer[T] {
	return DecorateConsumer(consumer, func(ctx Context[T, any, E]) error {
		return fallback(ctx.Param(), ctx.Err())
	}, DefaultFailurePredicate[T, any, E]())
}

func DecorateConsumerByType[T any, E error](
	consumer Consumer[T], fallback func(T) error) Consumer[T] {
	return DecorateConsumerWithFailure(consumer, func(t T, _ E) error { return fallback(t) })
}

func DecorateConsumerDefault[T any](
	consumer Consumer[T], fallback func(T) error) Consumer[T] {
	return DecorateConsumerByType[T, error](consumer, fallback)
}

func DecorateCheckedAcceptWithFailure[T any, E error](
	fn func(T) error, fallback func(T, E) error) func(T) error {
	return DecorateConsumerWithFailure(ConsumerCast(fn), fallback).CheckedAccept
}

func DecorateCheckedAcceptByType[T any, E error](
	fn func(T) error, fallback func(T) error) func(T) error {
	return DecorateConsumerByType[T, E](ConsumerCast(fn), fallback).CheckedAccept
}

func DecorateCheckedAcceptDefault[T any](
	fn func(T) error, fallback func(T) error) func(T) error {
	return DecorateConsumerDefault(ConsumerCast(fn), fallback).CheckedAccept
}
