package grpc

import (
	"context"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cloudeventscontext "github.com/cloudevents/sdk-go/v2/context"

	"open-cluster-management.io/sdk-go/pkg/cloudevents/generic/types"
)

func TestSourceContext(t *testing.T) {
	cases := []struct {
		name          string
		event         cloudevents.Event
		expectedTopic string
		assertError   func(error)
	}{
		{
			name: "unsupported event",
			event: func() cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetType("unsupported")
				return evt
			}(),
			assertError: func(err error) {
				if err == nil {
					t.Errorf("expected error, but failed")
				}
			},
		},
		{
			name: "resync status",
			event: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceStatus,
					Action:              types.ResyncRequestAction,
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				return evt
			}(),
			expectedTopic: "sources/hub1/clusters/statusresync",
			assertError: func(err error) {
				if err != nil {
					t.Errorf("unexpected error %v", err)
				}
			},
		},
		{
			name: "unsupported send resource no cluster name",
			event: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              "test",
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				return evt
			}(),
			assertError: func(err error) {
				if err == nil {
					t.Errorf("expected error, but failed")
				}
			},
		},
		{
			name: "send spec",
			event: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              "test",
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				evt.SetExtension("clustername", "cluster1")
				return evt
			}(),
			expectedTopic: "sources/hub1/clusters/cluster1/spec",
			assertError: func(err error) {
				if err != nil {
					t.Errorf("unexpected error %v", err)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sourceOptions := &gRPCSourceOptions{sourceID: "hub1"}
			ctx, err := sourceOptions.WithContext(context.TODO(), c.event.Context)
			c.assertError(err)

			topic := func(ctx context.Context) string {
				if ctx == nil {
					return ""
				}

				return cloudeventscontext.TopicFrom(ctx)
			}(ctx)

			if topic != c.expectedTopic {
				t.Errorf("expected %s, but got %s", c.expectedTopic, topic)
			}
		})
	}
}
