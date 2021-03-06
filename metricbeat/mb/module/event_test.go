// +build !integration

package module

import (
	"errors"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

const (
	moduleName    = "mymodule"
	metricSetName = "mymetricset"
	host          = "localhost"
	elapsed       = time.Duration(500 * time.Millisecond)
	tag           = "alpha"
)

var (
	startTime = time.Now()
	errFetch  = errors.New("error fetching data")
	tags      = []string{tag}
)

var builder = EventBuilder{
	ModuleName:    moduleName,
	MetricSetName: metricSetName,
	// host
	StartTime:     startTime,
	FetchDuration: elapsed,
	// event
	// fetchErr
	// processors
	// metadata
}

func TestEventBuilder(t *testing.T) {
	b := builder
	b.Host = host
	event, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, startTime, event.Timestamp)

	module := event.Fields[moduleName].(common.MapStr)
	metricset := event.Fields["metricset"].(common.MapStr)
	assert.Equal(t, moduleName, metricset["module"])
	assert.Equal(t, metricSetName, metricset["name"])
	assert.Equal(t, int64(500000), metricset["rtt"])
	assert.Equal(t, host, metricset["host"])
	assert.Equal(t, common.MapStr{}, module[metricSetName])
	assert.Nil(t, event.Fields["error"])
}

func TestEventBuilderError(t *testing.T) {
	b := builder
	b.fetchErr = errFetch
	event, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}

	errDoc := event.Fields["error"].(common.MapStr)
	assert.Equal(t, errFetch.Error(), errDoc["message"])
}

func TestEventBuilderNoHost(t *testing.T) {
	b := builder
	event, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}

	_, found := event.Fields["metricset-host"]
	assert.False(t, found)
}

func TestEventBuildNoRTT(t *testing.T) {
	b := builder
	b.FetchDuration = 0

	event, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}

	metricset := event.Fields["metricset"].(common.MapStr)
	_, found := metricset["rtt"]
	assert.False(t, found, "found rtt")
}
