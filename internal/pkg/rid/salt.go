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


package rid

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"os"
)

// Salt 计算机器 ID 的哈希值并返回一个 uint64 类型的盐值.
func Salt() uint64 {
	// 使用 FNV-1a 哈希算法计算字符串的哈希值
	hasher := fnv.New64a()
	hasher.Write(ReadMachineID())

	// 将哈希值转换为 uint64 型的盐
	hashValue := hasher.Sum64()
	return hashValue
}

// ReadMachineID 获取机器 ID，如果无法获取，则生成随机 ID.
func ReadMachineID() []byte {
	id := make([]byte, 3)
	machineID, err := readPlatformMachineID()
	if err != nil || len(machineID) == 0 {
		machineID, err = os.Hostname()
	}

	if err == nil && len(machineID) != 0 {
		hasher := sha256.New()
		hasher.Write([]byte(machineID))
		copy(id, hasher.Sum(nil))
	} else {
		// 如果无法收集机器 ID，则回退到生成随机数
		if _, randErr := rand.Reader.Read(id); randErr != nil {
			panic(fmt.Errorf("id: cannot get hostname nor generate a random number: %w; %w", err, randErr))
		}
	}
	return id
}

// readPlatformMachineID 尝试读取平台特定的机器 ID.
func readPlatformMachineID() (string, error) {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil || len(data) == 0 {
		data, err = os.ReadFile("/sys/class/dmi/id/product_uuid")
	}
	return string(data), err
}
