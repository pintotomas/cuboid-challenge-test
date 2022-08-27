package models

import (
	"encoding/json"
	"fmt"
)

type Bag struct {
	Model

	Title    string `validate:"required,max=255"`
	Volume   uint   `validate:"gt=0"`
	Disabled bool

	Cuboids []Cuboid
}

func (b *Bag) PayloadVolume() uint {
	var usedVolume uint = 0
	for i := range b.Cuboids {
		usedVolume += b.Cuboids[i].PayloadVolume()
	}
	return usedVolume
}

func (b *Bag) AvailableVolume() uint {
	if len(b.Cuboids) == 0 {
		return b.Volume
	}

	return b.Volume - b.PayloadVolume()
}

func (b *Bag) SetDisabled(value bool) {
	b.Disabled = value
}

func (b *Bag) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		ID              uint     `json:"id"`
		Title           string   `json:"title"`
		Volume          uint     `json:"volume"`
		PayloadVolume   uint     `json:"payloadVolume"`
		AvailableVolume uint     `json:"availableVolume"`
		Disabled        bool     `json:"disabled"`
		Cuboids         []Cuboid `json:"cuboids"`
	}{
		b.ID, b.Title, b.Volume, b.PayloadVolume(), b.AvailableVolume(), b.Disabled, b.Cuboids,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Bag. %w", err)
	}

	return j, nil
}
