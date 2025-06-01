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

package http

import (
	"github.com/gin-gonic/gin"
	"github.com/onexstack/onexstack/pkg/core"
)

// CreatePost 创建博客帖子.
func (h *Handler) CreatePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Create, h.val.ValidateCreatePostRequest)
}

// UpdatePost 更新博客帖子.
func (h *Handler) UpdatePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Update, h.val.ValidateUpdatePostRequest)
}

// DeletePost 删除博客帖子.
func (h *Handler) DeletePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Delete, h.val.ValidateDeletePostRequest)
}

// GetPost 获取博客帖子.
func (h *Handler) GetPost(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.PostV1().Get, h.val.ValidateGetPostRequest)
}

// ListPosts 列出用户的所有博客帖子.
func (h *Handler) ListPost(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.PostV1().List, h.val.ValidateListPostRequest)
}
