package aggregator

import (
	// stdlib
	"testing"

	// 3p
	"github.com/stretchr/testify/assert"
)

const epsilon = 0.01

func TestGaugeSampling(t *testing.T) {
	// Initialize a new Gauge
	mGauge := Gauge{}

	// Add samples
	mGauge.addSample(1, 50)
	mGauge.addSample(2, 55)

	series, _ := mGauge.flush(60)
	// the last sample is flushed
	assert.Len(t, series, 1)
	assert.Len(t, series[0].Points, 1)
	assert.InEpsilon(t, 2, series[0].Points[0].Value, epsilon)
	assert.EqualValues(t, 60, series[0].Points[0].Ts)
}