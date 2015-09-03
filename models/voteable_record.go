package models

import (
	"math"
)

// NOTE: Not even sure if I'm using this anymore

type voteableRecord struct {
	Ups        uint `json:"ups" bson:"ups,inline"`
	Downs      uint `json:"downs" bson:"downs,inline"`
	Confidence int  `json:"confidence" bson:"confidence,inline"`
	Heat       int  `json:"heat" bson:"heat,inline"`
	datedRecord
}

func (v *voteableRecord) UpVote() {
	v.Ups++
	v.calculateConfidence()
	v.calculateHeat()
}

func (v *voteableRecord) DownVote() {
	v.Downs++
	v.calculateConfidence()
	v.calculateHeat()
}

func (v *voteableRecord) calculateConfidence() {
	var n, z, phat, confidence float64

	n = float64(v.Ups + v.Downs)

	if n != 0 {
		z = 1.0 //1.0 = 85%, 1.6 = 95%
		phat = float64(v.Ups) / n

		// reddit algorithm... yeah, it's craziness
		confidence = math.Sqrt(phat+z*z/(2*n)-z*((phat*(1-phat)+z*z/(4*n))/n)) / (1 + z*z/n)
		v.Confidence = int(math.Floor(confidence * 1000))
	} else {
		v.Confidence = 0
	}
}

func (v *voteableRecord) calculateHeat() {
	var sign, score int

	score = int(v.Ups) - int(v.Downs)

	order := math.Log10(math.Max(math.Abs(float64(score)), 1))
	seconds := float64(v.Updated.Second() - 1383973200) // The crazy number is kind of arbitrary. Really just to make the seconds smaller.

	if score > 0 {
		sign = 1
	} else if score < 0 {
		sign = -1
	} else {
		sign = 0
	}

	buf := int(order+float64(sign)*seconds/45000.0) * int(math.Pow10(7))
	v.Heat = int(math.Floor(float64(buf) / math.Pow10(7) * 1000))
}
