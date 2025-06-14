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

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameCasbinRuleM = "casbin_rule"

// CasbinRuleM mapped from table <casbin_rule>
type CasbinRuleM struct {
	ID    int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PType *string `gorm:"column:ptype" json:"ptype"`
	V0    *string `gorm:"column:v0" json:"v0"`
	V1    *string `gorm:"column:v1" json:"v1"`
	V2    *string `gorm:"column:v2" json:"v2"`
	V3    *string `gorm:"column:v3" json:"v3"`
	V4    *string `gorm:"column:v4" json:"v4"`
	V5    *string `gorm:"column:v5" json:"v5"`
}

// TableName CasbinRuleM's table name
func (*CasbinRuleM) TableName() string {
	return TableNameCasbinRuleM
}
