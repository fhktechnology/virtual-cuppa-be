package repositories

import (
	"virtual-cuppa-be/models"

	"gorm.io/gorm"
)

type TagRepository interface {
	Create(tag *models.Tag) error
	FindByID(id uint) (*models.Tag, error)
	FindByUserID(userID uint) ([]models.Tag, error)
	FindByOrganisation(organisationID uint) ([]models.Tag, error)
	FindOrCreateByName(name string, organisationID uint) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(id uint) error
	AssignTagToUser(userID uint, tagID uint) error
	RemoveTagFromUser(userID uint, tagID uint) error
	ClearUserTags(userID uint) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindByUserID(userID uint) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("user_id = ?", userID).Find(&tags).Error
	return tags, err
}

func (r *tagRepository) FindByOrganisation(organisationID uint) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("organisation_id = ?", organisationID).Find(&tags).Error
	return tags, err
}

func (r *tagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

func (r *tagRepository) AssignTagToUser(userID uint, tagID uint) error {
	return r.db.Exec("INSERT INTO user_tags (user_id, tag_id) VALUES (?, ?) ON CONFLICT DO NOTHING", userID, tagID).Error
}

func (r *tagRepository) RemoveTagFromUser(userID uint, tagID uint) error {
	return r.db.Exec("DELETE FROM user_tags WHERE user_id = ? AND tag_id = ?", userID, tagID).Error
}

func (r *tagRepository) FindOrCreateByName(name string, organisationID uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ? AND organisation_id = ?", name, organisationID).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		tag = models.Tag{
			Name:           name,
			OrganisationID: organisationID,
		}
		if err := r.db.Create(&tag).Error; err != nil {
			return nil, err
		}
		return &tag, nil
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) ClearUserTags(userID uint) error {
	return r.db.Exec("DELETE FROM user_tags WHERE user_id = ?", userID).Error
}
