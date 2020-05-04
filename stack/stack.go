/*
 * Copyright 2020 The SealEVM Authors
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package stack

import (
	"SealEVM/evmErrors"
	"SealEVM/evmInt256"
)

type Stack struct {
	data [] *evmInt256.Int
	max int
}

func New(max int) *Stack {
	if max <= 0 {
		max = int(^uint(0) >> 1)
	}

	return &Stack{
		max: max,
	}
}

func (s *Stack) CheckMinDepth(min int) bool {
	sLen := len(s.data)
	return sLen >= min
}

func (s *Stack) Len() int {
	return len(s.data)
}

func (s *Stack) Push(i *evmInt256.Int) error {
	sLen := len(s.data)
	if sLen + 1 > s.max {
		return evmErrors.StackOverFlow
	}

	s.data = append(s.data, i)
	return nil
}

func (s *Stack) PushN(i []*evmInt256.Int) error {
	sLen := len(s.data)
	iLen := len(i)
	if sLen + iLen > s.max {
		return evmErrors.StackOverFlow
	}

	s.data = append(s.data, i...)
	return nil
}

func (s *Stack) Pop() (*evmInt256.Int, error) {
	sLen := len(s.data)
	if sLen == 0 {
		return nil, evmErrors.StackUnderFlow
	}

	i := s.data[sLen - 1]
	s.data = s.data[:sLen - 1]
	return i, nil
}

func (s *Stack) PopN(n int) ([]*evmInt256.Int, error) {
	sLen := len(s.data)
	var el []*evmInt256.Int
	if sLen >= n {
		el = s.data[sLen - n:]
		s.data = s.data[:sLen - n]
	} else {
		return nil, evmErrors.StackUnderFlow
	}

	//reverse to make sure the order
	for i, j := 0, len(el) - 1; i < j; i, j = i+1, j-1 {
		el[i], el[j] = el[j], el[i]
	}
	return el, nil
}

func (s *Stack) Peek() *evmInt256.Int {
	sLen := len(s.data)
	if sLen == 0 {
		return nil
	}

	i := s.data[sLen - 1]
	return i
}

func (s *Stack) PeekN(n int) []*evmInt256.Int {
	sLen := len(s.data)
	var el []*evmInt256.Int = nil
	if sLen >= n {
		el = s.data[sLen - n:]
	}

	return el
}

func (s *Stack) Swap(n int) error {
	n += 1
	sLen := len(s.data)
	if sLen < n {
		return evmErrors.StackUnderFlow
	}

	s.data[sLen - n], s.data[sLen - 1] = s.data[sLen - 1], s.data[sLen - n]

	return nil
}

func (s *Stack) Dup(n int) error {
	sLen := len(s.data)
	if sLen < n {
		return evmErrors.StackUnderFlow
	}

	i := s.data[sLen - n]
	err := s.Push(i)

	return err
}
