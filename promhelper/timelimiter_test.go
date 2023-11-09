package promhelper_test

import (
	"errors"
	"github.com/CharLemAznable/resilience4go/promhelper"
	"github.com/CharLemAznable/resilience4go/timelimiter"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestTimeLimiterRegistry(t *testing.T) {
	entry := timelimiter.NewTimeLimiter("test") // Create a new timelimiter entry for testing
	registerer := &testRegisterer{
		testingT: t,
		testCases: []*metricTestCase{
			{
				name: "TestSuccessCount",
				desc: `Desc{fqName: "resilience4go_timelimiter_calls", help: "The number of successful calls", constLabels: {kind="successful",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("successful")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Counter: &dto.Counter{
						Value: proto.Float64(1),
					},
				},
			},
			{
				name: "TestTimeoutCount",
				desc: `Desc{fqName: "resilience4go_timelimiter_calls", help: "The number of timed out calls", constLabels: {kind="timeout",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("timeout")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Counter: &dto.Counter{
						Value: proto.Float64(1),
					},
				},
			},
			{
				name: "TestFailureCount",
				desc: `Desc{fqName: "resilience4go_timelimiter_calls", help: "The number of failed calls", constLabels: {kind="failed",name="test"}, variableLabels: {}}`,
				metric: &dto.Metric{
					Label: []*dto.LabelPair{
						{Name: proto.String("kind"), Value: proto.String("failed")},
						{Name: proto.String("name"), Value: proto.String(entry.Name())},
					},
					Counter: &dto.Counter{
						Value: proto.Float64(1),
					},
				},
			},
		},
	}
	fn := timelimiter.DecorateRunnable(entry, func() error {
		panic("panic error")
	})
	func() {
		defer func() {
			if r := recover(); r != nil {
				// ignored
			}
		}()
		_ = fn()
	}()
	fn = timelimiter.DecorateRunnable(entry, func() error {
		time.Sleep(time.Second * 2)
		return nil
	})
	_ = fn()
	fn = timelimiter.DecorateRunnable(entry, func() error {
		time.Sleep(time.Millisecond * 500)
		return errors.New("error")
	})
	_ = fn()
	registerFn, unregisterFn := promhelper.TimeLimiterRegistry(entry)
	_ = registerFn(registerer)
	unregisterFn(registerer)
}
