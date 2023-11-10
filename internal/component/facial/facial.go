package facial

import (
	"icms/internal/repository/user"
	"icms/pkg/console"

	"github.com/Kagami/go-face"
	"github.com/go-kratos/kratos/v2/log"
)

type Config struct {
	ModelPath         string
	ClassifyThreshold float32
}

type Facial struct {
	recognizer        *face.Recognizer
	ClassifyThreshold float32
	userIds           []int32
	descriptors       []face.Descriptor
}

func (f *Facial) Recognize(imgData []byte) (faces *face.Face, err error) {
	return f.recognizer.RecognizeSingle(imgData)
}

func (f *Facial) RecognizeCNN(imgData []byte) (faces *face.Face, err error) {
	return f.recognizer.RecognizeSingleCNN(imgData)
}

func (f *Facial) Classify(descriptor face.Descriptor) (userId int) {
	return f.recognizer.ClassifyThreshold(descriptor, f.ClassifyThreshold)
}

func (f *Facial) AddSample(descriptors []face.Descriptor, userIds []int32) {
	f.descriptors = append(f.descriptors, descriptors...)
	f.userIds = append(f.userIds, userIds...)
	f.recognizer.SetSamples(f.descriptors, f.userIds)
}

func New(config *Config, logger log.Logger, userRepo *user.UserRepository) (*Facial, error) {
	rec, err := face.NewRecognizer(config.ModelPath)
	if err != nil {
		return nil, err
	}

	logger.Log(log.LevelInfo, "msg", "[Facial Recognition] start fetching facial descriptors")

	userIds, descriptors, err := userRepo.GetFacialDescriptors()

	if err != nil {
		console.Exit(err.Error())
	}

	rec.SetSamples(descriptors, userIds)

	logger.Log(log.LevelInfo, "msg", "[Facial Recognition] loaded")

	return &Facial{
		recognizer:        rec,
		ClassifyThreshold: config.ClassifyThreshold,
		userIds:           userIds,
		descriptors:       descriptors,
	}, nil
}
