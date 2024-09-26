package telemetry

import (
	"time"
)

type TraceData struct {
	TraceID string     `json:"traceID"`
	Spans   []SpanData `json:"spans"`
}

type RecentSummaries struct {
	TraceSummaries []TraceSummary `json:"traceSummaries"`
}

type TraceSummary struct {
	HasRootSpan bool `json:"hasRootSpan"`

	RootServiceName string    `json:"rootServiceName"`
	RootName        string    `json:"rootName"`
	RootStartTime   time.Time `json:"rootStartTime"`
	RootEndTime     time.Time `json:"rootEndTime"`

	SpanCount uint32 `json:"spanCount"`
	TraceID   string `json:"traceID"`
}

func (trace *TraceData) GetTraceSummary() TraceSummary {
	rootSpan, err := trace.getRootSpan()

	if err == ErrMissingRootSpan {
		return TraceSummary{
			HasRootSpan:     false,
			RootServiceName: "",
			RootName:        "",
			RootStartTime:   time.Time{},
			RootEndTime:     time.Time{},
			SpanCount:       uint32(len(trace.Spans)),
			TraceID:         trace.TraceID,
		}
	}

	return TraceSummary{
		HasRootSpan:     true,
		RootServiceName: rootSpan.GetServiceName(),
		RootName:        rootSpan.Name,
		RootStartTime:   rootSpan.StartTime,
		RootEndTime:     rootSpan.EndTime,
		SpanCount:       uint32(len(trace.Spans)),
		TraceID:         trace.TraceID,
	}
}

func (trace *TraceData) getRootSpan() (SpanData, error) {
	for i := 0; i < len(trace.Spans); i++ {
		if trace.Spans[i].ParentSpanID == "" {
			return trace.Spans[i], nil
		}
	}
	return SpanData{}, ErrMissingRootSpan
}