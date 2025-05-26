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

package grpc

import (
	"context"

	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
)

// CreatePost 创建博客帖子.
func (h *Handler) CreatePost(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	return h.biz.PostV1().Create(ctx, rq)
}

// UpdatePost 更新博客帖子.
func (h *Handler) UpdatePost(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	return h.biz.PostV1().Update(ctx, rq)
}

// DeletePost 删除博客帖子.
func (h *Handler) DeletePost(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	return h.biz.PostV1().Delete(ctx, rq)
}

// GetPost 获取博客帖子.
func (h *Handler) GetPost(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	return h.biz.PostV1().Get(ctx, rq)
}

// ListPost 列出所有博客帖子.
func (h *Handler) ListPost(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	return h.biz.PostV1().List(ctx, rq)
}
