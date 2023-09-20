package utils

import (
	"fmt"
	"math"
)

// Coordinates is a slice of float64
type Coordinates []float64

// Observations is a slice of observations
type Observations []Coordinates

// Distance returns the euclidean distance between two coordinates
func (c Coordinates) Distance(p2 Coordinates) float64 {
	var r float64
	for i, v := range c {
		r += math.Pow(v-p2[i], 2)
	}
	return r
}

// Center returns the center coordinates of a set of Observations
func (c Observations) Center() (Coordinates, error) {
	var l = len(c)
	if l == 0 {
		return nil, fmt.Errorf("there is no mean for an empty set of points")
	}

	cc := make([]float64, len(c[0]))
	for _, point := range c {
		for j, v := range point {
			cc[j] += v
		}
	}

	var mean Coordinates
	for _, v := range cc {
		mean = append(mean, v/float64(l))
	}
	return mean, nil
}

// Move obs from c1 to c2
func MoveObs(c1 *Coordinates, c2 *Coordinates, p float64) {
	for i := 0; i < len(*c1); i++ {
		val := (*c2)[i] / p * 0.01 // 1.001 is only for a practical implementation. To differenciate the clusters
		(*c1)[i] += val
		(*c2)[i] -= val
	}
}

// Sum two Coordinates
func (c1 *Coordinates) Sum(c2 Coordinates) {
	for i := 0; i < len(*c1); i++ {
		(*c1)[i] += c2[i]
	}
}

// Divide Coordinates
func (c *Coordinates) Divide(n float64) {
	for i := 0; i < len(*c); i++ {
		(*c)[i] = (*c)[i] / n
	}
}

// AverageDistance returns the average distance between o and all observations
func AverageDistance(o Coordinates, observations Observations) float64 {
	var d float64
	var l int

	for _, observation := range observations {
		dist := o.Distance(observation)
		if dist == 0 {
			continue
		}

		l++
		d += dist
	}

	if l == 0 {
		return 0
	}
	return d / float64(l)
}
