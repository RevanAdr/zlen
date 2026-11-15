package auth

import "context"

type UserService interface {
	List(ctx context.Context) ([]UserResponse, error)
	Get(ctx context.Context, id int) (UserResponse, error)
}

type UserServiceImpl struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) List(ctx context.Context) ([]UserResponse, error) {
	responses, err := s.repo.List(ctx)
	if err != nil {
		return []UserResponse{}, err
	}
	return ToResponses(responses), nil
}

func (s *UserServiceImpl) Get(ctx context.Context, id int) (UserResponse, error) {
	response, err := s.repo.Get(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	return ToResponse(response), nil
}
