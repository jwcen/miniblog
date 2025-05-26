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

package post

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/jwcen/miniblog/internal/apiserver/model"
	"github.com/jwcen/miniblog/internal/apiserver/pkg/conversion"
	"github.com/jwcen/miniblog/internal/apiserver/store"
	"github.com/jwcen/miniblog/internal/pkg/contextx"
	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
	"github.com/onexstack/onexstack/pkg/store/where"
)

type PostBiz interface {
	Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error)
	Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error)
	List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error)

	PostExpansion
}

// PostExpansion 定义额外的帖子操作方法.
type PostExpansion interface{}

type postBiz struct {
	store store.IStore
}

var _ PostBiz = (*postBiz)(nil)

func New(store store.IStore) *postBiz {
	return &postBiz{store: store}
}

func (b *postBiz) Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	var postM model.PostM
	_ = copier.Copy(&postM, rq)
	postM.UserID = contextx.UserID(ctx)

	if err := b.store.Post().Create(ctx, &postM); err != nil {
		return nil, err
	}

	return &apiv1.CreatePostResponse{PostID: postM.PostID}, nil
}

func (b *postBiz) Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostID())
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	if rq.Title != nil {
		postM.Title = rq.GetTitle()
	}

	if rq.Content != nil {
		postM.Content = rq.GetContent()
	}

	if err := b.store.Post().Update(ctx, postM); err != nil {
		return nil, err
	}

	return &apiv1.UpdatePostResponse{}, nil
}

func (b *postBiz) Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostIDs())
	if err := b.store.Post().Delete(ctx, whr); err != nil {
		return nil, err
	}

	return &apiv1.DeletePostResponse{}, nil
}

func (b *postBiz) Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostID())
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	return &apiv1.GetPostResponse{Post: conversion.PostModelToPostV1(postM)}, nil
}

func (b *postBiz) List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	whr := where.T(ctx).P(int(rq.GetOffset()), int(rq.GetLimit()))
	count, postList, err := b.store.Post().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	posts := make([]*apiv1.Post, 0, len(postList))
	for _, post := range postList {
		converted := conversion.PostModelToPostV1(post)
		posts = append(posts, converted)
	}

	return &apiv1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}
