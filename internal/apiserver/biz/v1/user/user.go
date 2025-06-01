// Copyright 2024 jayvee <jvvcen@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"context"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jwcen/miniblog/internal/apiserver/model"
	"github.com/jwcen/miniblog/internal/apiserver/pkg/conversion"
	"github.com/jwcen/miniblog/internal/apiserver/store"
	"github.com/jwcen/miniblog/internal/pkg/contextx"
	"github.com/jwcen/miniblog/internal/pkg/errno"
	"github.com/jwcen/miniblog/internal/pkg/known"
	"github.com/jwcen/miniblog/internal/pkg/log"
	"github.com/jwcen/miniblog/internal/pkg/rid"
	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
	"github.com/jwcen/miniblog/pkg/auth"
	"github.com/onexstack/onexstack/pkg/store/where"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error)
	Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error)
	List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)

	UserExpansion
}

// UserExpansion 定义用户操作的扩展方法.
type UserExpansion interface {
	Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error)
	RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error)
	ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)
}

type userBiz struct {
	store store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(store store.IStore) *userBiz {
	return &userBiz{
		store: store,
	}
}

func (u *userBiz) Login(ctx context.Context, req *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	whr := where.F("username", req.GetUsername())
	userM, err := u.store.User().Get(ctx, whr)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	if err := auth.Compare(userM.Password, req.GetPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}

	// TODO：实现 Token 签发逻辑

	return &apiv1.LoginResponse{
		Token:    "<placeholder>",
		ExpireAt: timestamppb.New(time.Now().Add(2 * time.Hour)),
	}, nil
}

// RefreshToken 用于刷新用户的身份验证令牌.
// 当用户的令牌即将过期时，可以调用此方法生成一个新的令牌.
func (u *userBiz) RefreshToken(ctx context.Context, req *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error) {
	// TODO：实现 Token 签发逻辑
	return &apiv1.RefreshTokenResponse{Token: "<placeholder>", ExpireAt: timestamppb.New(time.Now().Add(2 * time.Hour))}, nil
}

func (u *userBiz) ChangePassword(ctx context.Context, req *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error) {
	userM, err := u.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	if err := auth.Compare(userM.Password, req.GetOldPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}

	userM.Password, _ = auth.Encrypt(req.GetNewPassword())
	if err := u.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &apiv1.ChangePasswordResponse{}, nil
}

func (u *userBiz) ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	return nil, nil
}

func (u *userBiz) Create(ctx context.Context, req *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	var userM model.UserM

	if err := copier.Copy(&userM, req); err != nil {
		log.W(ctx).Errorw("copier.Copy(&userM, req) failed", "err", err)
		return nil, errno.ErrInvalidArgument
	}

	// 设置创建时间和更新时间
	now := time.Now()
	userM.CreatedAt = now
	userM.UpdatedAt = now

	// 生成临时 userID，后续会被 AfterCreate 钩子更新
	userM.UserID = rid.UserID.New(0)

	if err := u.store.User().Create(ctx, &userM); err != nil {
		return nil, err
	}

	return &apiv1.CreateUserResponse{
		UserID: userM.UserID,
	}, nil
}

func (u *userBiz) Update(ctx context.Context, req *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	userM, err := u.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	if req.Username != nil {
		userM.Username = *req.Username
	}
	if req.Email != nil {
		userM.Email = req.GetEmail()
	}
	if req.Nickname != nil {
		userM.Nickname = req.GetNickname()
	}
	if req.Phone != nil {
		userM.Phone = req.GetPhone()
	}

	if err := u.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &apiv1.UpdateUserResponse{}, nil
}

func (u *userBiz) Delete(ctx context.Context, req *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	// 只有 `root` 用户可以删除用户，并且可以删除其他用户
	// 所以这里不用 where.T()，因为 where.T() 会查询 `root` 用户自己
	if err := u.store.User().Delete(ctx, where.F("userID", req.GetUserID())); err != nil {
		return nil, err
	}

	return &apiv1.DeleteUserResponse{}, nil
}

func (u *userBiz) Get(ctx context.Context, req *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	userM, err := u.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	return &apiv1.GetUserResponse{
		User: conversion.UserModelToUserV1(userM),
	}, nil
}

func (u *userBiz) List(ctx context.Context, req *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	whr := where.P(int(req.GetOffset()), int(req.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}

	count, userList, err := u.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	// 设置最大并发数量为常量 MaxConcurrency
	eg.SetLimit(known.MaxErrGroupConcurrency)

	// 使用 goroutine 提高接口性能
	for _, user := range userList {
		user := user // 创建新的变量以避免闭包问题
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				count, _, err := u.store.Post().List(ctx, whr)
				if err != nil {
					return err
				}

				converted := conversion.UserModelToUserV1(user)
				converted.PostCount = count
				m.Store(user.UserID, converted) // 使用 UserID 作为 key

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.W(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	users := make([]*apiv1.User, 0, len(userList))
	for _, item := range userList {
		if user, ok := m.Load(item.UserID); ok { // 使用 UserID 作为 key
			users = append(users, user.(*apiv1.User))
		}
	}

	log.W(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &apiv1.ListUserResponse{
		TotalCount: count,
		Users:      users,
	}, nil
}
