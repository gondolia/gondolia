package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/auth"
	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/repository"
)

// UserService handles user operations
type UserService struct {
	userRepo        repository.UserRepository
	userCompanyRepo repository.UserCompanyRepository
	roleRepo        repository.RoleRepository
	authLogRepo     repository.AuthLogRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	userCompanyRepo repository.UserCompanyRepository,
	roleRepo repository.RoleRepository,
	authLogRepo repository.AuthLogRepository,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		userCompanyRepo: userCompanyRepo,
		roleRepo:        roleRepo,
		authLogRepo:     authLogRepo,
	}
}

// List returns a list of users
func (s *UserService) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error) {
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.userRepo.List(ctx, filter)
}

// GetByID returns a user by ID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load companies
	companies, err := s.userCompanyRepo.ListByUser(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Companies = companies

	return user, nil
}

// GetByEmail returns a user by email
func (s *UserService) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, tenantID, email)
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, tenantID uuid.UUID, req domain.CreateUserRequest) (*domain.User, error) {
	// Check if email already exists
	_, err := s.userRepo.GetByEmail(ctx, tenantID, req.Email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Create user
	user := domain.NewUser(tenantID, req.Email, req.FirstName, req.LastName)
	user.Phone = req.Phone
	user.Mobile = req.Mobile

	// If password provided, set it
	if req.Password != nil {
		if err := auth.ValidatePasswordStrength(*req.Password); err != nil {
			return nil, domain.ErrPasswordTooWeak
		}
		hash, err := auth.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = &hash
		user.IsActive = true
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates a user
func (s *UserService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Mobile != nil {
		user.Mobile = req.Mobile
	}
	if req.DefaultLanguage != nil {
		user.DefaultLanguage = *req.DefaultLanguage
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete soft deletes a user
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

// Activate activates a user
func (s *UserService) Activate(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.IsActive = true
	return s.userRepo.Update(ctx, user)
}

// Deactivate deactivates a user
func (s *UserService) Deactivate(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.IsActive = false
	return s.userRepo.Update(ctx, user)
}

// SetPassword sets a user's password
func (s *UserService) SetPassword(ctx context.Context, id uuid.UUID, password string) error {
	if err := auth.ValidatePasswordStrength(password); err != nil {
		return domain.ErrPasswordTooWeak
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	return s.userRepo.SetPassword(ctx, id, hash)
}

// VerifyPassword verifies a user's password
func (s *UserService) VerifyPassword(ctx context.Context, id uuid.UUID, password string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}

	if user.PasswordHash == nil {
		return false, nil
	}

	return auth.VerifyPassword(password, *user.PasswordHash), nil
}

// InviteUserToCompany invites a new user to a company
func (s *UserService) InviteUserToCompany(ctx context.Context, tenantID uuid.UUID, req domain.InviteUserRequest) (*domain.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, tenantID, req.Email)
	if err == nil {
		// User exists, add to company instead
		return s.addExistingUserToCompany(ctx, existingUser, req.CompanyID, req.RoleID)
	}

	// Validate role exists
	_, err = s.roleRepo.GetByID(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}

	// Generate invitation token
	token, err := auth.GenerateInvitationToken()
	if err != nil {
		return nil, err
	}

	// Create user with invitation
	user := domain.NewUser(tenantID, req.Email, req.FirstName, req.LastName)
	user.Phone = req.Phone
	user.Mobile = req.Mobile
	user.InvitationToken = &token
	now := time.Now()
	user.InvitedAt = &now
	user.IsActive = true // Active but can't login until invitation accepted

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Add to company
	uc := domain.NewUserCompany(user.ID, req.CompanyID, &req.RoleID, domain.UserTypeInvited)
	if err := s.userCompanyRepo.Create(ctx, uc); err != nil {
		return nil, err
	}

	// Log invitation
	s.logAuthEvent(ctx, tenantID, &user.ID, domain.AuthEventInvitationSent, "", "", map[string]any{
		"company_id": req.CompanyID,
	})

	return user, nil
}

// AcceptInvitation accepts an invitation and sets password
func (s *UserService) AcceptInvitation(ctx context.Context, token, password string) (*domain.User, error) {
	// Find user by invitation token
	user, err := s.userRepo.GetByInvitationToken(ctx, token)
	if err != nil {
		return nil, domain.ErrInvitationInvalid
	}

	// Check invitation validity
	if !user.HasOpenInvitation() {
		return nil, domain.ErrInvitationInvalid
	}

	if user.HasExpiredInvitation() {
		return nil, domain.ErrInvitationExpired
	}

	// Validate password
	if err := auth.ValidatePasswordStrength(password); err != nil {
		return nil, domain.ErrPasswordTooWeak
	}

	// Hash password
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Update user
	if err := s.userRepo.SetPassword(ctx, user.ID, hash); err != nil {
		return nil, err
	}

	if err := s.userRepo.ClearInvitation(ctx, user.ID); err != nil {
		return nil, err
	}

	// Log acceptance
	s.logAuthEvent(ctx, user.TenantID, &user.ID, domain.AuthEventInvitationAccepted, "", "", nil)

	// Return updated user
	return s.userRepo.GetByID(ctx, user.ID)
}

// ValidateInvitationToken validates an invitation token
func (s *UserService) ValidateInvitationToken(ctx context.Context, token string) (*domain.User, error) {
	user, err := s.userRepo.GetByInvitationToken(ctx, token)
	if err != nil {
		return nil, domain.ErrInvitationInvalid
	}

	if !user.HasOpenInvitation() {
		return nil, domain.ErrInvitationInvalid
	}

	if user.HasExpiredInvitation() {
		return nil, domain.ErrInvitationExpired
	}

	return user, nil
}

// addExistingUserToCompany adds an existing user to a company
func (s *UserService) addExistingUserToCompany(ctx context.Context, user *domain.User, companyID, roleID uuid.UUID) (*domain.User, error) {
	// Check if already in company
	_, err := s.userCompanyRepo.GetByUserAndCompany(ctx, user.ID, companyID)
	if err == nil {
		return nil, domain.ErrUserAlreadyInCompany
	}

	// Validate role exists
	_, err = s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Add to company
	uc := domain.NewUserCompany(user.ID, companyID, &roleID, domain.UserTypeInvited)
	if err := s.userCompanyRepo.Create(ctx, uc); err != nil {
		return nil, err
	}

	// Activate user if not active
	if !user.IsActive {
		user.IsActive = true
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) logAuthEvent(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, eventType, ipAddress, userAgent string, metadata map[string]any) {
	log := &domain.AuthenticationLog{
		ID:        uuid.New(),
		TenantID:  tenantID,
		UserID:    userID,
		EventType: eventType,
		IPAddress: nilIfEmpty(ipAddress),
		UserAgent: nilIfEmpty(userAgent),
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}
	_ = s.authLogRepo.Create(ctx, log)
}
