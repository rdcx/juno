package service

import (
	"juno/pkg/api/user"
	"juno/pkg/util"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger   logrus.FieldLogger
	userRepo user.Repository
}

func New(logger logrus.FieldLogger, userRepo user.Repository) *Service {
	return &Service{
		logger:   logger,
		userRepo: userRepo,
	}
}

func validateEmail(email string) error {
	if !util.IsValidEmail(email) {
		return user.ErrInvalidEmail
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return user.ErrInvalidPassword
	}

	return nil
}

func (s *Service) Get(id uuid.UUID) (*user.User, error) {
	return s.userRepo.Get(id)
}

func (s *Service) Create(email, password string) error {

	var errs []error

	if err := validateEmail(email); err != nil {
		errs = append(errs, err)
	}

	if err := validatePassword(password); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return util.ValidationErrs(errs)
	}

	if u, err := s.userRepo.FirstWhereEmail(email); err == nil || u != nil {
		return user.ErrEmailAlreadyExists
	}

	hash, err := util.BcryptPassword(password)

	if err != nil {
		return err
	}

	u := &user.User{
		Email:    email,
		Password: hash,
	}

	return s.userRepo.Create(u)
}

func (s *Service) Update(u *user.User) error {
	//TODO: Implement
	return nil
}

func (s *Service) Delete(id uuid.UUID) error {
	//TODO: Implement
	return nil
}