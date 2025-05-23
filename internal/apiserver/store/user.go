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

package store

import (
	"context"
	"errors"

	"github.com/jwcen/miniblog/internal/apiserver/model"
	"github.com/jwcen/miniblog/internal/pkg/errno"
	"github.com/jwcen/miniblog/internal/pkg/log"
	"github.com/onexstack/onexstack/pkg/store/where"
	"gorm.io/gorm"
)

// UserStore 定义了 user 模块在 store 层所实现的方法.
type UserStore interface {
	Create(ctx context.Context, obj *model.UserM) error
	Update(ctx context.Context, obj *model.UserM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.UserM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.UserM, error)

	UserExpansion
}

// UserExpansion 定义了用户操作的附加方法.
type UserExpansion interface{}

type userStore struct {
	store *datastore
}

var _ UserStore = (*userStore)(nil)

func newUserStore(store *datastore) *userStore {
	return &userStore{store}
}

func (s *userStore) Create(ctx context.Context, user *model.UserM) error {
	if err := s.store.DB(ctx).Create(&user).Error; err != nil {
		log.Errorw("Failed to insert user into database", "err", err, "user", user)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}

	return nil
}

func (s *userStore) Update(ctx context.Context, obj *model.UserM) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		log.Errorw("Failed to update user in database", "err", err, "user", obj)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}

	return nil
}

// Delete 根据条件删除用户记录.
func (s *userStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx, opts).Delete(new(model.UserM)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Failed to delete user from database", "err", err, "conditions", opts)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}

	return nil
}

// Get 根据条件查询用户记录.
func (s *userStore) Get(ctx context.Context, opts *where.Options) (*model.UserM, error) {
	var obj model.UserM
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		log.Errorw("Failed to retrieve user from database", "err", err, "conditions", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDBRead.WithMessage(err.Error())
	}

	return &obj, nil
}

// List 返回用户列表和总数.
// nolint: nonamedreturns
func (s *userStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.UserM, err error) {
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		log.Errorw("Failed to list users from database", "err", err, "conditions", opts)
		err = errno.ErrDBRead.WithMessage(err.Error())
	}
	return
}
