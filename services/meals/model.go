// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package meals

import "time"

type Stats struct {
	Weight Weight `json:"weight,omitempty"`
	Volume Volume `json:"volume,omitempty"`
	Energy Energy `json:"energy,omitempty"`

	Fat           Weight `json:"fat,omitempty"`
	Carbohydrates Weight `json:"carbohydrates,omitempty"`
	Sugar         Weight `json:"sugar,omitempty"`
	Protein       Weight `json:"protein,omitempty"`
	Salt          Weight `json:"salt,omitempty"`
	Water         Volume `json:"water,omitempty"`
}

func (s *Stats) Multiply(size float64) {
	s.Weight *= Weight(size)
	s.Energy *= Energy(size)

	s.Fat *= Weight(size)
	s.Carbohydrates *= Weight(size)
	s.Sugar *= Weight(size)
	s.Protein *= Weight(size)
	s.Salt *= Weight(size)
	s.Water *= Volume(size)
}

func (s *Stats) ScaleToWeight(w Weight) {
	if s.Weight != 0 {
		s.Multiply(float64(w / s.Weight))
	}
}

func (s *Stats) ScaleToVolume(v Volume) {
	if s.Volume != 0 {
		s.Multiply(float64(v / s.Volume))
	}
}

func (s *Stats) Normalize() {
	if s.Volume != 0 {
		s.ScaleToVolume(100 * Millilitre)
	} else {
		s.ScaleToWeight(100 * Gram)
	}
}

func (s *Stats) Add(o Stats) {
	s.Weight += o.Weight
	s.Energy += o.Energy

	s.Carbohydrates += o.Carbohydrates
	s.Sugar += o.Sugar
	s.Fat += o.Fat
	s.Protein += o.Protein
	s.Salt += o.Salt
	s.Water += o.Water
}

type Product struct {
	Id    int64  `json:"-"`
	RefId int64  `json:"ref_id,omitempty"`
	Name  string `json:"name,omitempty"`
	Code  string `json:"code,omitempty"`

	ServingWeight Weight `json:"serving_weight,omitempty"`
	ServingVolume Volume `json:"serving_volume,omitempty"`
	Stats

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (p Product) TableName() string {
	return "meals_products"
}

func (p Product) Servings(n float64) Stats {
	stats := p.Stats
	if p.ServingVolume != 0 {
		stats.ScaleToVolume(p.ServingVolume)
	} else {
		stats.ScaleToWeight(p.ServingWeight)
	}
	stats.Multiply(n)
	return stats
}

type Serving struct {
	Id           int64     `json:"-"`
	RefId        int64     `json:"ref_id,omitempty"`
	Name         string    `json:"name"`
	AmountWeight Weight    `json:"amount_weight"`
	AmountVolume Volume    `json:"amount_volume"`
	Time         time.Time `json:"time" sql:"index"`

	ProductId int64    `json:"-"`
	Product   *Product `json:"product,omitempty"`

	Size            float64 `json:"size,omitempty" sql:"-"`
	CalculatedStats *Stats  `json:"stats,omitempty"`
}

func (s Serving) TableName() string {
	return "meals_servings"
}

func (s *Serving) Stats() Stats {
	if s.Product == nil {
		return Stats{}
	}
	stats := s.Product.Stats
	if s.AmountWeight > 0 {
		stats.ScaleToWeight(s.AmountWeight)
	} else if s.AmountVolume > 0 {
		stats.ScaleToVolume(s.AmountVolume)
	}
	s.CalculatedStats = &stats
	return stats
}
