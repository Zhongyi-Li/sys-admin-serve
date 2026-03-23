package auth

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	authdto "sys-admin-serve/internal/dto/auth"
	"sys-admin-serve/internal/model"
	jwtutil "sys-admin-serve/internal/pkg/jwt"
	repositoryauth "sys-admin-serve/internal/repository/auth"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrUserNotFound       = errors.New("user not found")
)

type Service struct {
	repo       *repositoryauth.Repository
	jwtManager *jwtutil.Manager
	log        *zap.Logger
}

func NewService(repo *repositoryauth.Repository, jwtManager *jwtutil.Manager, log *zap.Logger) *Service {
	return &Service{repo: repo, jwtManager: jwtManager, log: log}
}

func (s *Service) Login(ctx context.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	if user == nil {
		s.log.Warn("login failed: user not found", zap.String("username", username))
		return nil, ErrInvalidCredentials
	}

	if user.Status != 1 {
		s.log.Warn("login failed: user disabled", zap.Uint64("user_id", user.ID), zap.String("username", user.Username))
		return nil, ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.log.Warn("login failed: invalid password", zap.Uint64("user_id", user.ID), zap.String("username", user.Username))
		return nil, ErrInvalidCredentials
	}

	roles, err := s.repo.ListRoleCodesByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("list user roles: %w", err)
	}

	accessToken, expiresAt, err := s.jwtManager.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	if err := s.repo.UpdateUserLastLogin(ctx, user.ID); err != nil {
		s.log.Warn("update last login failed", zap.Uint64("user_id", user.ID), zap.Error(err))
	}

	s.log.Info("user login succeeded", zap.Uint64("user_id", user.ID), zap.String("username", user.Username))

	return &authdto.LoginResponse{
		TokenType:   "Bearer",
		AccessToken: accessToken,
		ExpiresAt:   expiresAt.Unix(),
		User: authdto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Status:   user.Status,
			Roles:    roles,
		},
	}, nil
}

func (s *Service) GetCurrentUser(ctx context.Context, userID uint64) (*authdto.UserInfo, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	roles, err := s.repo.ListRoleCodesByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("list user roles: %w", err)
	}

	return &authdto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Status:   user.Status,
		Roles:    roles,
	}, nil
}

func (s *Service) GetCurrentUserMenus(ctx context.Context, userID uint64) ([]authdto.UserMenu, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	menus, err := s.repo.ListMenusByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list menus by user id: %w", err)
	}

	return buildMenuTree(menus), nil
}

func buildMenuTree(menus []model.Menu) []authdto.UserMenu {
	nodeMap := make(map[uint64]*authdto.UserMenu, len(menus))
	roots := make([]*authdto.UserMenu, 0)

	for i := range menus {
		menu := menus[i]
		nodeMap[menu.ID] = &authdto.UserMenu{
			ID:         menu.ID,
			ParentID:   menu.ParentID,
			Name:       menu.Name,
			Title:      menu.Title,
			Path:       menu.Path,
			Component:  menu.Component,
			Icon:       menu.Icon,
			MenuType:   menu.MenuType,
			Permission: menu.Permission,
			Sort:       menu.Sort,
			Hidden:     menu.Hidden,
			Children:   make([]authdto.UserMenu, 0),
		}
	}

	for i := range menus {
		node := nodeMap[menus[i].ID]
		if node.ParentID == 0 {
			roots = append(roots, node)
			continue
		}

		parent, ok := nodeMap[node.ParentID]
		if !ok {
			roots = append(roots, node)
			continue
		}

		parent.Children = append(parent.Children, *node)
	}

	for i := range roots {
		sortMenuChildren(&roots[i].Children)
	}

	sort.SliceStable(roots, func(i, j int) bool {
		if roots[i].Sort != roots[j].Sort {
			return roots[i].Sort < roots[j].Sort
		}
		return roots[i].ID < roots[j].ID
	})

	result := make([]authdto.UserMenu, 0, len(roots))
	for i := range roots {
		result = append(result, *roots[i])
	}

	return result
}

func sortMenuChildren(children *[]authdto.UserMenu) {
	sort.SliceStable(*children, func(i, j int) bool {
		if (*children)[i].Sort != (*children)[j].Sort {
			return (*children)[i].Sort < (*children)[j].Sort
		}
		return (*children)[i].ID < (*children)[j].ID
	})

	for i := range *children {
		sortMenuChildren(&(*children)[i].Children)
	}
}
