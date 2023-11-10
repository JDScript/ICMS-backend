package user

import (
	"icms/internal/model"

	"github.com/Kagami/go-face"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) GetFacialDescriptors() (userIds []int32, facialDescriptors []face.Descriptor, err error) {
	descriptors := []model.FacialDescriptor{}
	err = repo.db.Select("user_id", "binary_facial_descriptor").Find(&descriptors).Error
	if err != nil {
		return nil, nil, err
	}

	userIds = make([]int32, len(descriptors))
	facialDescriptors = make([]face.Descriptor, len(descriptors))

	for idx := range descriptors {
		userIds[idx] = descriptors[idx].UserID
		facialDescriptors[idx] = face.Descriptor(descriptors[idx].FacialDescriptor)
	}

	return
}

func (repo *UserRepository) GetByID(userId string) (user *model.User) {
	repo.db.Omit("binary_facial_descriptor").Where("id", userId).Find(&user)
	return
}

func (repo *UserRepository) Create(user *model.User) error {
	return repo.db.Create(user).Error
}
