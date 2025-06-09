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

package biz

import (
	"github.com/google/wire"
	postV1 "github.com/jwcen/miniblog/internal/apiserver/biz/v1/post"
	userV1 "github.com/jwcen/miniblog/internal/apiserver/biz/v1/user"
	"github.com/jwcen/miniblog/internal/apiserver/store"
	auth "github.com/onexstack/onexstack/pkg/authz"
)

var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(IBiz), new(*biz)))

type IBiz interface {
	// 获取用户业务接口.
	UserV1() userV1.UserBiz
	// 获取帖子业务接口.
	PostV1() postV1.PostBiz
}

type biz struct {
	store store.IStore
	authz *auth.Authz
}

var _ IBiz = (*biz)(nil)

func NewBiz(store store.IStore, authz *auth.Authz) *biz {
	return &biz{
		store: store,
		authz: authz,
	}
}

func (b *biz) UserV1() userV1.UserBiz {
	return userV1.New(b.store, b.authz)
}

func (b *biz) PostV1() postV1.PostBiz {
	return postV1.New(b.store)
}
