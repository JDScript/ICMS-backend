package model

import (
	"bytes"
	"encoding/gob"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FacialDescriptor struct {
	DescriptorID uuid.UUID `gorm:"primaryKey" json:"descriptor_id"`
	UserID       int32     `gorm:"primaryKey" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`

	BinaryFacialDescriptor []byte       `gorm:"type:blob;not null" json:"-"`
	FacialDescriptor       [128]float32 `gorm:"-" json:"-"`
}

func (descriptor *FacialDescriptor) BeforeCreate(tx *gorm.DB) (err error) {
	descriptor.DescriptorID = uuid.New()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(descriptor.FacialDescriptor); err != nil {
		return err
	}
	descriptor.BinaryFacialDescriptor = buf.Bytes()
	return
}

func (descriptor *FacialDescriptor) AfterFind(tx *gorm.DB) (err error) {
	if descriptor.BinaryFacialDescriptor != nil {
		reader := bytes.NewReader(descriptor.BinaryFacialDescriptor)
		dec := gob.NewDecoder(reader)
		if err := dec.Decode(&descriptor.FacialDescriptor); err != nil {
			return err
		}
	}
	return
}
