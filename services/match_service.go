package services

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
	"virtual-cuppa-be/models"
	"virtual-cuppa-be/repositories"
)

var (
	ErrNoUsersToMatch        = errors.New("not enough users to create matches")
	ErrMatchNotFound         = errors.New("match not found")
	ErrUnauthorizedMatch     = errors.New("unauthorized to modify this match")
	ErrFeedbackAlreadyExists = errors.New("feedback already submitted for this match")
	ErrMatchNotAccepted      = errors.New("can only provide feedback for accepted matches")
	ErrInvalidRating         = errors.New("rating must be between 1 and 5")
)

// AvailabilitySlot represents a single availability time slot for email templates
type AvailabilitySlot struct {
	Day    string `json:"Day"`
	Period string `json:"Period"`
}

type MatchService interface {
	GenerateMatchesForOrganisation(organisationID uint) (int, error)
	TryGenerateMatchForUser(userID uint) error
	GetCurrentMatch(userID uint) (*models.Match, error)
	GetMatchHistory(userID uint) ([]*models.Match, error)
	AcceptMatch(userID uint, matchID uint) error
	AcceptMatchWithAvailability(userID uint, matchID uint, availability models.Availability) (*models.Match, error)
	RejectMatch(userID uint, matchID uint) error
	GetOrganisationMatches(organisationID uint) ([]*models.Match, error)
	GetMatchAvailabilities(userID uint, matchID uint) ([]*models.MatchAvailability, error)
	SubmitFeedback(userID uint, matchID uint, rating int, comment string) error
	GetMatchFeedbacks(userID uint, matchID uint) ([]*models.MatchFeedback, error)
	GetMatchesPendingFeedback(userID uint) ([]*models.Match, error)
}

type matchService struct {
	matchRepo         repositories.MatchRepository
	matchHistoryRepo  repositories.MatchHistoryRepository
	matchFeedbackRepo repositories.MatchFeedbackRepository
	userRepo          repositories.UserRepository
	availConfigRepo   repositories.UserAvailabilityConfigRepository
	emailSvc          EmailService
}

func NewMatchService(
	matchRepo repositories.MatchRepository,
	matchHistoryRepo repositories.MatchHistoryRepository,
	matchFeedbackRepo repositories.MatchFeedbackRepository,
	userRepo repositories.UserRepository,
	availConfigRepo repositories.UserAvailabilityConfigRepository,
	emailSvc EmailService,
) MatchService {
	return &matchService{
		matchRepo:         matchRepo,
		matchHistoryRepo:  matchHistoryRepo,
		matchFeedbackRepo: matchFeedbackRepo,
		userRepo:          userRepo,
		availConfigRepo:   availConfigRepo,
		emailSvc:          emailSvc,
	}
}

// calculateMatchScore calculates compatibility score based on common tags
func (s *matchService) calculateMatchScore(user1Tags, user2Tags []models.Tag) float64 {
	if len(user1Tags) == 0 && len(user2Tags) == 0 {
		return 0
	}

	tagMap := make(map[string]bool)
	for _, tag := range user1Tags {
		tagMap[tag.Name] = true
	}

	commonTags := 0
	for _, tag := range user2Tags {
		if tagMap[tag.Name] {
			commonTags++
		}
	}

	totalUniqueTags := len(tagMap)
	for _, tag := range user2Tags {
		if !tagMap[tag.Name] {
			totalUniqueTags++
		}
	}

	if totalUniqueTags == 0 {
		return 0
	}

	return float64(commonTags) / float64(totalUniqueTags) * 100
}

type userPair struct {
	user1ID uint
	user2ID uint
	score   float64
}

func (s *matchService) GenerateMatchesForOrganisation(organisationID uint) (int, error) {
	// Get all confirmed users in organisation
	users, err := s.userRepo.FindByOrganisation(organisationID)
	if err != nil {
		return 0, err
	}

	// Filter only confirmed users without pending matches (exclude Admins)
	// Also check if users have availability configuration
	var availableUsers []*models.User
	var userConfigs = make(map[uint]*models.UserAvailabilityConfig)
	
	for _, user := range users {
		// Skip admins - they should not be matched
		if user.AccountType == models.AccountTypeAdmin {
			continue
		}
		if !user.IsConfirmed {
			continue
		}
		hasPending, err := s.matchRepo.HasPendingMatch(user.ID)
		if err != nil {
			return 0, err
		}
		if hasPending {
			continue
		}
		
		// Check if user has availability configuration
		config, err := s.availConfigRepo.FindByUserID(user.ID)
		if err != nil {
			return 0, err
		}
		if config == nil {
			// User doesn't have availability config, skip
			continue
		}
		
		userConfigs[user.ID] = config
		availableUsers = append(availableUsers, user)
	}

	if len(availableUsers) < 2 {
		return 0, ErrNoUsersToMatch
	}

	// Calculate scores for all possible pairs
	var pairs []userPair
	for i := 0; i < len(availableUsers); i++ {
		for j := i + 1; j < len(availableUsers); j++ {
			user1 := availableUsers[i]
			user2 := availableUsers[j]

			// Check if they were ever matched before
			wasMatched, err := s.matchHistoryRepo.WasEverMatched(user1.ID, user2.ID)
			if err != nil {
				return 0, err
			}
			if wasMatched {
				continue
			}
			
			// Check if users have common availability
			config1 := userConfigs[user1.ID]
			config2 := userConfigs[user2.ID]
			if !models.HasCommonAvailability(config1, config2) {
				// No common time slots, skip this pair
				continue
			}

			score := s.calculateMatchScore(user1.Tags, user2.Tags)
			pairs = append(pairs, userPair{
				user1ID: user1.ID,
				user2ID: user2.ID,
				score:   score,
			})
		}
	}

	if len(pairs) == 0 {
		return 0, ErrNoUsersToMatch
	}

	// Sort pairs by score (highest first)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score > pairs[j].score
	})

	// Greedy algorithm: select best matches ensuring each user appears once
	matchedUsers := make(map[uint]bool)
	var matchesToCreate []userPair

	for _, pair := range pairs {
		if !matchedUsers[pair.user1ID] && !matchedUsers[pair.user2ID] {
			matchesToCreate = append(matchesToCreate, pair)
			matchedUsers[pair.user1ID] = true
			matchedUsers[pair.user2ID] = true
		}
	}

	if len(matchesToCreate) == 0 {
		return 0, ErrNoUsersToMatch
	}

	// Create matches with random dates in next week
	createdCount := 0
	for _, pair := range matchesToCreate {
		scheduledDate, scheduledTime := s.generateRandomDateTime()

		match := &models.Match{
			OrganisationID: organisationID,
			User1ID:        pair.user1ID,
			User2ID:        pair.user2ID,
			MatchScore:     pair.score,
			Status:         models.MatchStatusPending,
			ScheduledDate:  scheduledDate,
			ScheduledTime:  scheduledTime,
		}

		if err := s.matchRepo.Create(match); err != nil {
			continue
		}

		// Add to history
		history := &models.MatchHistory{
			User1ID:   pair.user1ID,
			User2ID:   pair.user2ID,
			MatchedAt: time.Now(),
		}
		s.matchHistoryRepo.Create(history)

		createdCount++
	}

	return createdCount, nil
}

// TryGenerateMatchForUser tries to find and create a match for a specific user
// Returns nil if no match can be found (no error)
func (s *matchService) TryGenerateMatchForUser(userID uint) error {
	// Get user and check if confirmed and not admin
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil // User not found, no error
	}

	if user.AccountType == models.AccountTypeAdmin || !user.IsConfirmed {
		return nil // Admins and unconfirmed users don't get matched
	}

	// Check if user has organisation
	if user.OrganisationID == nil {
		return nil // User not assigned to organisation
	}

	// Check if user already has a pending match
	hasPending, err := s.matchRepo.HasPendingMatch(userID)
	if err != nil {
		return nil
	}
	if hasPending {
		return nil // User already has a pending match
	}
	
	// Check if user has availability configuration
	userConfig, err := s.availConfigRepo.FindByUserID(userID)
	if err != nil {
		return nil
	}
	if userConfig == nil {
		return nil // User doesn't have availability config
	}

	// Get all confirmed users in the same organisation
	users, err := s.userRepo.FindByOrganisation(*user.OrganisationID)
	if err != nil {
		return nil
	}

	// Find available users (excluding the current user, admins, unconfirmed, and those with pending matches)
	var candidates []*models.User
	var candidateConfigs = make(map[uint]*models.UserAvailabilityConfig)
	
	for _, u := range users {
		if u.ID == userID || u.AccountType == models.AccountTypeAdmin || !u.IsConfirmed {
			continue
		}
		hasPending, err := s.matchRepo.HasPendingMatch(u.ID)
		if err != nil || hasPending {
			continue
		}

		// Check if they were ever matched before
		wasMatched, err := s.matchHistoryRepo.WasEverMatched(userID, u.ID)
		if err != nil || wasMatched {
			continue
		}
		
		// Check if candidate has availability configuration
		candidateConfig, err := s.availConfigRepo.FindByUserID(u.ID)
		if err != nil || candidateConfig == nil {
			continue
		}
		
		// Check if they have common availability
		if !models.HasCommonAvailability(userConfig, candidateConfig) {
			continue
		}
		
		candidateConfigs[u.ID] = candidateConfig
		candidates = append(candidates, u)
	}

	if len(candidates) == 0 {
		return nil // No available candidates, no error
	}

	// Find the best match based on score
	var bestCandidate *models.User
	var bestScore float64 = -1

	for _, candidate := range candidates {
		score := s.calculateMatchScore(user.Tags, candidate.Tags)
		if score > bestScore {
			bestScore = score
			bestCandidate = candidate
		}
	}

	if bestCandidate == nil {
		return nil // No suitable candidate found
	}

	// Create the match
	scheduledDate, scheduledTime := s.generateRandomDateTime()

	match := &models.Match{
		OrganisationID: *user.OrganisationID,
		User1ID:        userID,
		User2ID:        bestCandidate.ID,
		MatchScore:     bestScore,
		Status:         models.MatchStatusPending,
		ScheduledDate:  scheduledDate,
		ScheduledTime:  scheduledTime,
	}

	if err := s.matchRepo.Create(match); err != nil {
		return nil // Failed to create match, but no error to caller
	}

	// Add to history
	history := &models.MatchHistory{
		User1ID:   userID,
		User2ID:   bestCandidate.ID,
		MatchedAt: time.Now(),
	}
	s.matchHistoryRepo.Create(history)

	return nil
}

func (s *matchService) generateRandomDateTime() (time.Time, string) {
	// Random day in next 7 days (excluding weekends)
	daysAhead := rand.Intn(5) + 1 // 1-5 days ahead
	date := time.Now().AddDate(0, 0, daysAhead)

	// Ensure it's a weekday
	for date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		daysAhead++
		date = time.Now().AddDate(0, 0, daysAhead)
	}

	// Random time slot between 9 AM and 5 PM
	hours := []string{"9 AM", "10 AM", "11 AM", "2 PM", "3 PM", "4 PM", "5 PM"}
	timeSlot := hours[rand.Intn(len(hours))]

	return date, timeSlot
}

func (s *matchService) GetCurrentMatch(userID uint) (*models.Match, error) {
	match, err := s.matchRepo.FindCurrentByUserID(userID)
	if err != nil {
		return nil, ErrMatchNotFound
	}

	// Check if both users have already submitted feedback
	// If so, this match should not be considered "current"
	feedbackCount, err := s.matchFeedbackRepo.CountFeedbacksByMatch(match.ID)
	if err == nil && feedbackCount >= 2 {
		return nil, ErrMatchNotFound
	}

	return match, nil
}

func (s *matchService) GetMatchHistory(userID uint) ([]*models.Match, error) {
	return s.matchRepo.FindByUserID(userID)
}

func (s *matchService) AcceptMatch(userID uint, matchID uint) error {
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return ErrMatchNotFound
	}

	if match.User1ID != userID && match.User2ID != userID {
		return ErrUnauthorizedMatch
	}

	now := time.Now()

	// Mark the current user as accepted with timestamp
	if match.User1ID == userID {
		match.User1Accepted = true
		match.User1AcceptedAt = &now
	} else {
		match.User2Accepted = true
		match.User2AcceptedAt = &now
	}

	// If both users accepted, change status to waiting_for_feedback and set expiration (5 days from now)
	if match.User1Accepted && match.User2Accepted {
		match.Status = models.MatchStatusWaitingForFeedback
		expiresAt := now.AddDate(0, 0, 5) // 5 days from now
		match.ExpiresAt = &expiresAt
	}

	if err := s.matchRepo.Update(match); err != nil {
		return err
	}

	// Send email notification to the OTHER user when current user accepts
	// Reload match with user details for email
	updatedMatch, err := s.matchRepo.FindByID(matchID)
	if err == nil && updatedMatch != nil {
		var otherUser *models.User
		var currentUser *models.User
		
		if updatedMatch.User1ID == userID {
			otherUser = updatedMatch.User2
			currentUser = updatedMatch.User1
		} else {
			otherUser = updatedMatch.User1
			currentUser = updatedMatch.User2
		}

		if otherUser != nil && currentUser != nil {
			// Get current user's availability config
			availConfig, err := s.availConfigRepo.FindByUserID(currentUser.ID)
			if err == nil && availConfig != nil {
				// Convert availability config to Availability map format
				availability := availConfig.ToAvailability()
				availabilitySlots := s.formatAvailabilitySlots(availability)
				
				// If no slots, don't send email (config exists but all slots are false)
				if len(availabilitySlots) == 0 {
					fmt.Printf("WARNING: User %d has availability config but no slots enabled\n", currentUser.ID)
					return nil
				}
				
				currentUserName := fmt.Sprintf("%s %s", currentUser.FirstName, currentUser.LastName)
				otherUserName := fmt.Sprintf("%s %s", otherUser.FirstName, otherUser.LastName)
				
				err = s.emailSvc.SendMatchAccepted(
					otherUser.Email,
					otherUserName,
					currentUserName,
					currentUser.Email,
					availabilitySlots,
				)
				if err != nil {
					fmt.Printf("ERROR: Failed to send match accepted email: %v\n", err)
				}
			}
		}
	}

	return nil
}

func (s *matchService) AcceptMatchWithAvailability(userID uint, matchID uint, availability models.Availability) (*models.Match, error) {
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return nil, ErrMatchNotFound
	}

	if match.User1ID != userID && match.User2ID != userID {
		return nil, ErrUnauthorizedMatch
	}

	now := time.Now()

	// Mark the current user as accepted with timestamp
	if match.User1ID == userID {
		match.User1Accepted = true
		match.User1AcceptedAt = &now
	} else {
		match.User2Accepted = true
		match.User2AcceptedAt = &now
	}

	// If both users accepted, change status to waiting_for_feedback and set expiration (5 days from now)
	if match.User1Accepted && match.User2Accepted {
		match.Status = models.MatchStatusWaitingForFeedback
		expiresAt := now.AddDate(0, 0, 5) // 5 days from now
		match.ExpiresAt = &expiresAt
	}

	// Save or update availability
	existingAvailability, err := s.matchRepo.FindAvailabilityByMatchAndUser(matchID, userID)
	if err != nil || existingAvailability == nil {
		// Create new availability
		newAvailability := &models.MatchAvailability{
			MatchID:      matchID,
			UserID:       userID,
			Availability: availability,
		}
		if err := s.matchRepo.CreateAvailability(newAvailability); err != nil {
			return nil, err
		}
	} else {
		// Update existing availability
		existingAvailability.Availability = availability
		if err := s.matchRepo.UpdateAvailability(existingAvailability); err != nil {
			return nil, err
		}
	}

	// Update match
	if err := s.matchRepo.Update(match); err != nil {
		return nil, err
	}

	// Reload match with availabilities
	updatedMatch, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return nil, err
	}

	// Send email notification to the OTHER user when current user accepts
	// The other user gets notified about the current user's availability
	var otherUser *models.User
	var currentUserName string
	
	if updatedMatch.User1ID == userID {
		otherUser = updatedMatch.User2
		currentUserName = fmt.Sprintf("%s %s", updatedMatch.User1.FirstName, updatedMatch.User1.LastName)
	} else {
		otherUser = updatedMatch.User1
		currentUserName = fmt.Sprintf("%s %s", updatedMatch.User2.FirstName, updatedMatch.User2.LastName)
	}

	// Format the current user's availability as structured data for SendGrid
	if otherUser != nil {
		availabilitySlots := s.formatAvailabilitySlots(availability)
		
		var currentUserEmail string
		if updatedMatch.User1ID == userID {
			currentUserEmail = updatedMatch.User1.Email
		} else {
			currentUserEmail = updatedMatch.User2.Email
		}
		
		s.emailSvc.SendMatchAccepted(
			otherUser.Email,
			fmt.Sprintf("%s %s", otherUser.FirstName, otherUser.LastName),
			currentUserName,
			currentUserEmail,
			availabilitySlots,
		)
	}

	return updatedMatch, nil
}

// formatAvailabilitySlots converts availability map to structured array for SendGrid template
func (s *matchService) formatAvailabilitySlots(availability models.Availability) []AvailabilitySlot {
	if len(availability) == 0 {
		return []AvailabilitySlot{}
	}
	
	// Define weekday order for consistent display
	weekdayOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	
	// Map periods to readable names
	periodNames := map[string]string{
		"morning":   "Morning",
		"afternoon": "Afternoon",
	}
	
	var slots []AvailabilitySlot
	
	for _, weekday := range weekdayOrder {
		periods, exists := availability[weekday]
		if !exists || len(periods) == 0 {
			continue
		}
		
		for _, period := range periods {
			periodName := periodNames[period]
			if periodName == "" {
				periodName = period
			}
			
			slots = append(slots, AvailabilitySlot{
				Day:    weekday,
				Period: periodName,
			})
		}
	}
	
	return slots
}

func (s *matchService) GetMatchAvailabilities(userID uint, matchID uint) ([]*models.MatchAvailability, error) {
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return nil, ErrMatchNotFound
	}

	// Check if user is part of this match
	if match.User1ID != userID && match.User2ID != userID {
		return nil, ErrUnauthorizedMatch
	}

	return s.matchRepo.FindAvailabilitiesByMatch(matchID)
}

func (s *matchService) RejectMatch(userID uint, matchID uint) error {
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return ErrMatchNotFound
	}

	if match.User1ID != userID && match.User2ID != userID {
		return ErrUnauthorizedMatch
	}

	// If any user rejects, the whole match is rejected
	match.Status = models.MatchStatusRejected
	return s.matchRepo.Update(match)
}

func (s *matchService) GetOrganisationMatches(organisationID uint) ([]*models.Match, error) {
	return s.matchRepo.FindByOrganisation(organisationID)
}

func (s *matchService) SubmitFeedback(userID uint, matchID uint, rating int, comment string) error {
	// Validate rating
	if rating < 1 || rating > 5 {
		return ErrInvalidRating
	}

	// Get match
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return ErrMatchNotFound
	}

	// Check if user is part of this match
	if match.User1ID != userID && match.User2ID != userID {
		return ErrUnauthorizedMatch
	}

	// Check if match is waiting for feedback (both users accepted)
	if match.Status != models.MatchStatusWaitingForFeedback {
		return ErrMatchNotAccepted
	}

	// Check if feedback already exists
	hasFeedback, err := s.matchFeedbackRepo.HasFeedback(matchID, userID)
	if err != nil {
		return err
	}
	if hasFeedback {
		return ErrFeedbackAlreadyExists
	}

	// Create feedback
	feedback := &models.MatchFeedback{
		MatchID: matchID,
		UserID:  userID,
		Rating:  rating,
		Comment: comment,
	}

	if err := s.matchFeedbackRepo.Create(feedback); err != nil {
		return err
	}

	// Determine who is being rated (the other user in the match)
	var ratedUserID uint
	if match.User1ID == userID {
		ratedUserID = match.User2ID
	} else {
		ratedUserID = match.User1ID
	}

	// Update average rating for the rated user
	if err := s.updateUserAverageRating(ratedUserID); err != nil {
		return err
	}

	// Check if both users have now submitted feedback
	feedbackCount, err := s.matchFeedbackRepo.CountFeedbacksByMatch(matchID)
	if err == nil && feedbackCount >= 2 {
		// Both users gave feedback, mark match as completed
		match.Status = models.MatchStatusCompleted
		if updateErr := s.matchRepo.Update(match); updateErr != nil {
			// Log error but don't fail the feedback submission
			fmt.Printf("Warning: Failed to update match status to completed: %v\n", updateErr)
		}
		
		// Try to generate new matches for both users
		// These operations are fire-and-forget - errors are ignored
		go s.TryGenerateMatchForUser(match.User1ID)
		go s.TryGenerateMatchForUser(match.User2ID)
	}

	return nil
}

func (s *matchService) updateUserAverageRating(userID uint) error {
	// Get all feedbacks where this user was rated (they were the OTHER person in the match)
	matches, err := s.matchRepo.FindByUserID(userID)
	if err != nil {
		return err
	}

	var totalRating float64
	var feedbackCount int

	for _, match := range matches {
		// Get feedbacks for this match
		feedbacks, err := s.matchFeedbackRepo.FindByMatch(match.ID)
		if err != nil {
			continue
		}

		// Find feedback where the OTHER user rated this user
		for _, feedback := range feedbacks {
			// If this is a feedback from the other user (not from userID themselves)
			if feedback.UserID != userID {
				// Check if userID is part of this match
				if match.User1ID == userID || match.User2ID == userID {
					totalRating += float64(feedback.Rating)
					feedbackCount++
				}
			}
		}
	}

	// Calculate average
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if feedbackCount > 0 {
		avg := totalRating / float64(feedbackCount)
		user.AverageRating = &avg
	} else {
		user.AverageRating = nil
	}

	return s.userRepo.Update(user)
}

func (s *matchService) GetMatchFeedbacks(userID uint, matchID uint) ([]*models.MatchFeedback, error) {
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return nil, ErrMatchNotFound
	}

	// Check if user is part of this match
	if match.User1ID != userID && match.User2ID != userID {
		return nil, ErrUnauthorizedMatch
	}

	return s.matchFeedbackRepo.FindByMatch(matchID)
}

func (s *matchService) GetMatchesPendingFeedback(userID uint) ([]*models.Match, error) {
	// Get all user's matches
	matches, err := s.matchRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var pendingFeedback []*models.Match

	for _, match := range matches {
		// Only waiting_for_feedback matches (both users accepted, waiting for feedback)
		if match.Status != models.MatchStatusWaitingForFeedback {
			continue
		}

		// Check if user already submitted feedback
		hasFeedback, err := s.matchFeedbackRepo.HasFeedback(match.ID, userID)
		if err != nil {
			continue
		}

		if !hasFeedback {
			pendingFeedback = append(pendingFeedback, match)
		}
	}

	return pendingFeedback, nil
}

