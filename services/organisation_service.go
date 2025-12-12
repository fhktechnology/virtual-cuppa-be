package services

import (
	"log"
	"virtual-cuppa-be/models"
	"virtual-cuppa-be/repositories"
)

type OrganisationService interface {
	UpsertOrganisation(input *models.UpsertOrganisationInput) (*models.Organisation, error)
	GetOrganisationByID(id uint) (*models.Organisation, error)
}

type organisationService struct {
	orgRepo repositories.OrganisationRepository
}

func NewOrganisationService(orgRepo repositories.OrganisationRepository) OrganisationService {
	return &organisationService{
		orgRepo: orgRepo,
	}
}

func (s *organisationService) UpsertOrganisation(input *models.UpsertOrganisationInput) (*models.Organisation, error) {
	organisation := &models.Organisation{
		ID:         input.ID,
		Name:       input.Name,
		CompanyUrl: input.CompanyUrl,
	}

	if err := s.orgRepo.Upsert(organisation); err != nil {
		return nil, err
	}

	log.Printf("Organisation created/updated with ID: %d", organisation.ID)

	return organisation, nil
}

func (s *organisationService) GetOrganisationByID(id uint) (*models.Organisation, error) {
	return s.orgRepo.FindByID(id)
}
