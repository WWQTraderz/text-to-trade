package user

import (
	"context"
	"errors"

	userpb "github.com/tjons/text-to-trade/pkg/api/user"
	"github.com/tjons/text-to-trade/pkg/model"
	"gorm.io/gorm"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	db *gorm.DB
}

func NewUserService(db *gorm.DB) userpb.UserServiceServer {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, req *userpb.User) (*userpb.User, error) {
	user := &model.User{
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		FirebaseUID: req.FirebaseUid,
		Username:    req.Username,
	}

	experienceLevel, ok := userpb.ExperienceLevel_name[int32(req.ExperienceLevel)]
	if !ok {
		return nil, errors.New("invalid experience level")
	}
	user.Experience = model.ExperienceLevel(experienceLevel)

	allocation, ok := userpb.Allocation_name[int32(req.Allocation)]
	if !ok {
		return nil, errors.New("invalid allocation")
	}
	user.Allocation = model.Allocation(allocation)

	risk, ok := userpb.RiskLevel_name[int32(req.RiskLevel)]
	if !ok {
		return nil, errors.New("invalid risk level")
	}
	user.Risk = model.RiskLevel(risk)

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return &userpb.User{
		InternalId:      uint32(user.ID),
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		FirebaseUid:     user.FirebaseUID,
		Username:        user.Username,
		ExperienceLevel: userpb.ExperienceLevel(userpb.ExperienceLevel_value[string(user.Experience)]),
		Allocation:      userpb.Allocation(userpb.Allocation_value[string(user.Allocation)]),
		RiskLevel:       userpb.RiskLevel(userpb.RiskLevel_value[string(user.Risk)]),
	}, nil

}

func (s *UserService) GetUser(ctx context.Context, req *userpb.User) (*userpb.User, error) {
	user := &model.User{}
	if err := s.db.First(user, req.InternalId).Error; err != nil {
		return nil, err
	}

	return &userpb.User{
		InternalId:      uint32(user.ID),
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		FirebaseUid:     user.FirebaseUID,
		Username:        user.Username,
		ExperienceLevel: userpb.ExperienceLevel(userpb.ExperienceLevel_value[string(user.Experience)]),
		Allocation:      userpb.Allocation(userpb.Allocation_value[string(user.Allocation)]),
		RiskLevel:       userpb.RiskLevel(userpb.RiskLevel_value[string(user.Risk)]),
	}, nil
}

func (s *UserService) OnboardFlow(ctx context.Context, req *userpb.UserFlowRequest) (*userpb.UserFlowResponse, error) {
	user := &model.User{}
	response := &userpb.UserFlowResponse{}
	if err := s.db.First(user, req.InternalId).Error; err != nil {
		return nil, err
	}

	switch req.CurrentStep {
	case userpb.Step_EXPERIENCE:
		_, ok := userpb.ExperienceLevel_value[req.Response]
		if !ok {
			return nil, errors.New("Unknown experience level")
		}
		user.Experience = model.ExperienceLevel(req.Response)
		response.NextStep = userpb.Step_ALLOCATION
		response.Message = "Choose Allocation"
		response.Options = []string{userpb.Allocation_LONG_TERM.String(), userpb.Allocation_SHORT_TERM.String()}
	case userpb.Step_ALLOCATION:
		_, ok := userpb.Allocation_value[req.Response]
		if !ok {
			return nil, errors.New("invalid allocation")
		}
		user.Allocation = model.Allocation(req.Response)
		response.NextStep = userpb.Step_RISK
		response.Message = "Choose Risk Level"
		response.Options = []string{userpb.RiskLevel_LOW.String(), userpb.RiskLevel_HIGH.String()}
	case userpb.Step_RISK:
		risk, ok := userpb.RiskLevel_value[req.Response]
		if !ok {
			return nil, errors.New("invalid risk level")
		}
		user.Risk = model.RiskLevel(risk)
		return s.sendWelcomeMessage(user)
	default:
		if user.Onboarded {
			return s.sendWelcomeMessage(user)
		} else {
			return s.startOnboarding(user)
		}
	}

	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return response, nil
}

func (s *UserService) startOnboarding(user *model.User) (*userpb.UserFlowResponse, error) {
	user.Onboarded = false
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return &userpb.UserFlowResponse{
		NextStep: userpb.Step_EXPERIENCE,
		Message:  "Hi! How comfortable are you with trading?",
		Options:  []string{userpb.ExperienceLevel_BEGINNER.String(), userpb.ExperienceLevel_ADVANCED.String()},
	}, nil
}

func (s *UserService) sendWelcomeMessage(user *model.User) (*userpb.UserFlowResponse, error) {
	user.Onboarded = true
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return &userpb.UserFlowResponse{
		Message: "Welcome to Text to Trade! You're all set up and ready to go!",
	}, nil
}
