package repositories

import (
	"virtual-cuppa-be/models"

	"gorm.io/gorm"
)

type OrganisationRepository interface {
	FindByID(id uint) (*models.Organisation, error)
	FindByName(name string) (*models.Organisation, error)
	FindAll() ([]*models.Organisation, error)
	Create(organisation *models.Organisation) error
	Update(organisation *models.Organisation) error
	Upsert(organisation *models.Organisation) error
}

type organisationRepository struct {
	db *gorm.DB
}

func NewOrganisationRepository(db *gorm.DB) OrganisationRepository {
	return &organisationRepository{
		db: db,
	}
}

func (r *organisationRepository) FindByID(id uint) (*models.Organisation, error) {
	var organisation models.Organisation
	err := r.db.First(&organisation, id).Error
	if err != nil {
		return nil, err
	}
	return &organisation, nil
}

func (r *organisationRepository) FindByName(name string) (*models.Organisation, error) {
	var organisation models.Organisation
	err := r.db.Where("name = ?", name).First(&organisation).Error
	if err != nil {
		return nil, err
	}
	return &organisation, nil
}

func (r *organisationRepository) FindAll() ([]*models.Organisation, error) {
	var organisations []*models.Organisation
	err := r.db.Find(&organisations).Error
	return organisations, err
}

func (r *organisationRepository) Create(organisation *models.Organisation) error {
	return r.db.Create(organisation).Error
}

func (r *organisationRepository) Update(organisation *models.Organisation) error {
	return r.db.Save(organisation).Error
}

func (r *organisationRepository) Upsert(organisation *models.Organisation) error {
	var existing models.Organisation
	err := r.db.Where("id = ?", organisation.ID).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(organisation).Error
	}
	
	if err != nil {
		return err
	}
	
	return r.db.Model(&existing).Updates(organisation).Error
}
