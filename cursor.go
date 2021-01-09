package rouge

import "math"

// Cursor seeks to the start of the next peak/valley.
type Cursor struct {
	seekingOuter bool

	InnerThreshold float64
	OuterThreshold float64
	Position       uint64
}

// Progress walks a cursor forward along the signal by one sample.
func (o *Cursor) Progress() {
	o.Position++
}

// CheckPeak examines the current sample amplitude,
// returning whether the sample represents the start of an extreme peak/valley.
func (o *Cursor) CheckPeak(sample float64) bool {
	sampleAbs := math.Abs(sample)

	if o.seekingOuter && sampleAbs > o.OuterThreshold {
		o.seekingOuter = !o.seekingOuter
		return true
	}

	if !o.seekingOuter && sampleAbs < o.InnerThreshold {
		o.seekingOuter = !o.seekingOuter
	}

	return false
}

// CheckStem examines the current sample amplitude,
// returning whether the represents a
func (o *Cursor) CheckStem(sample float64) bool {
	return math.Abs(sample) > o.InnerThreshold
}
