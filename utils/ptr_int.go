// Copyright 2023 Stichting ThingsIX Foundation
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
//
// SPDX-License-Identifier: Apache-2.0

package utils

func IntPtrToUintPtr(i *int) *uint {
	if i == nil {
		return nil
	}

	return Ptr(uint(*i))
}

func UintPtrToIntPtr(i *uint) *int {
	if i == nil {
		return nil
	}

	return Ptr(int(*i))
}

func Int32PtrToIntPtr(i *int32) *int {
	if i == nil {
		return nil
	}

	return Ptr(int(*i))
}

func IntPtrToInt32Ptr(i *int) *int32 {
	if i == nil {
		return nil
	}

	return Ptr(int32(*i))
}
