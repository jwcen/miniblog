// User API 定义，包含用户信息、登录请求和响应等相关消息

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

// Code generated by protoc-gen-defaults. DO NOT EDIT.

package v1

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	_ *timestamppb.Timestamp
	_ *durationpb.Duration
	_ *wrapperspb.BoolValue
)

func (x *User) Default() {
}

func (x *LoginRequest) Default() {
}

func (x *LoginResponse) Default() {
}

func (x *RefreshTokenRequest) Default() {
}

func (x *RefreshTokenResponse) Default() {
}

func (x *ChangePasswordRequest) Default() {
}

func (x *ChangePasswordResponse) Default() {
}

func (x *CreateUserRequest) Default() {
	if x.Nickname == nil {
		v := string("你好世界")
		x.Nickname = &v
	}
}

func (x *CreateUserResponse) Default() {
}

func (x *UpdateUserRequest) Default() {
}

func (x *UpdateUserResponse) Default() {
}

func (x *DeleteUserRequest) Default() {
}

func (x *DeleteUserResponse) Default() {
}

func (x *GetUserRequest) Default() {
}

func (x *GetUserResponse) Default() {
}

func (x *ListUserRequest) Default() {
}

func (x *ListUserResponse) Default() {
}
