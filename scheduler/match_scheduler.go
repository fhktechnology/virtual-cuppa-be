package scheduler

import (
	"log"
	"time"

	"virtual-cuppa-be/repositories"
	"virtual-cuppa-be/services"
)

type MatchScheduler struct {
	matchService services.MatchService
	orgRepo      repositories.OrganisationRepository
	stopChan     chan bool
	ticker       *time.Ticker
}

func NewMatchScheduler(
	matchService services.MatchService,
	orgRepo repositories.OrganisationRepository,
) *MatchScheduler {
	return &MatchScheduler{
		matchService: matchService,
		orgRepo:      orgRepo,
		stopChan:     make(chan bool),
	}
}

// Start begins the scheduler that generates matches every Monday at 9 AM
func (s *MatchScheduler) Start() {
	log.Println("Match scheduler started - will generate matches every Monday at 9 AM")

	// Start immediate check
	go s.checkAndGenerateMatches()

	// Calculate time until next Monday 9 AM
	s.scheduleNextRun()

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.checkAndGenerateMatches()
				s.scheduleNextRun()
			case <-s.stopChan:
				s.ticker.Stop()
				log.Println("Match scheduler stopped")
				return
			}
		}
	}()
}

func (s *MatchScheduler) scheduleNextRun() {
	now := time.Now()
	
	// Find next Monday at 9 AM
	daysUntilMonday := (8 - int(now.Weekday())) % 7
	if daysUntilMonday == 0 {
		// It's Monday - check if it's before 9 AM
		targetTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
		if now.Before(targetTime) {
			// Run today at 9 AM
			daysUntilMonday = 0
		} else {
			// Run next Monday
			daysUntilMonday = 7
		}
	}
	
	nextMonday := now.AddDate(0, 0, daysUntilMonday)
	nextRun := time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(), 9, 0, 0, 0, nextMonday.Location())
	
	duration := nextRun.Sub(now)
	
	if s.ticker != nil {
		s.ticker.Stop()
	}
	
	s.ticker = time.NewTicker(duration)
	
	log.Printf("Next match generation scheduled for: %s (in %v)", nextRun.Format("2006-01-02 15:04:05"), duration.Round(time.Minute))
}

func (s *MatchScheduler) checkAndGenerateMatches() {
	log.Println("Starting automatic match generation for all organisations...")

	// Get all organisations
	orgs, err := s.orgRepo.FindAll()
	if err != nil {
		log.Printf("Error fetching organisations: %v", err)
		return
	}

	totalMatches := 0
	successfulOrgs := 0

	for _, org := range orgs {
		count, err := s.matchService.GenerateMatchesForOrganisation(org.ID)
		if err != nil {
			if err == services.ErrNoUsersToMatch {
				log.Printf("Organisation %s (ID: %d): Not enough users to generate matches", org.Name, org.ID)
			} else {
				log.Printf("Error generating matches for organisation %s (ID: %d): %v", org.Name, org.ID, err)
			}
			continue
		}

		totalMatches += count
		successfulOrgs++
		log.Printf("Organisation %s (ID: %d): Generated %d matches", org.Name, org.ID, count)
	}

	log.Printf("Match generation completed: %d matches created across %d organisations (total: %d organisations)", 
		totalMatches, successfulOrgs, len(orgs))
}

func (s *MatchScheduler) Stop() {
	s.stopChan <- true
}

// RunNow triggers immediate match generation (for testing or manual triggers)
func (s *MatchScheduler) RunNow() {
	log.Println("Manual match generation triggered")
	go s.checkAndGenerateMatches()
}
