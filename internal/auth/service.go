package auth

import "context"

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) List(ctx context.Context) ([]UserResponse, error) {
	responses, err := s.repo.List(ctx)
	if err != nil {
		return []UserResponse{}, err
	}
	return ToResponses(responses), nil
}

func (s *UserService) Get(ctx context.Context, id int) (UserResponse, error) {
	response, err := s.repo.Get(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	return ToResponse(response), nil
}

func (s *UserService) Create(ctx context.Context, user UserRequest) (UserCreateResponse, error) {
	response, err := s.repo.Create(ctx, user)
	if err != nil {
		return UserCreateResponse{}, err
	}
	return ToCreateRepsonse(response), nil
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	_, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	rErr := s.repo.Delete(ctx, id)
	if rErr != nil {
		return rErr
	}
	return nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (LoginResponse, error) {
	response, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return LoginResponse{}, err
	}
	return response, nil
}
