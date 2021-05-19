package marisfrolg_utils

import "sync"

type Stack struct {
	data     []interface{}
	length   int
	capacity int
	sync.Mutex
}

// 构建一个空栈
func InitStack() *Stack {
	return &Stack{data: make([]interface{}, 8), length: 0, capacity: 8}
}

// 压栈操作
func (s *Stack) Push(data interface{}) {
	s.Lock()
	defer s.Unlock()

	if s.length+1 >= s.capacity {
		s.capacity <<= 1
		t := s.data
		s.data = make([]interface{}, s.capacity)
		copy(s.data, t)
	}
	s.data[s.length] = data
	s.length++

}

// 出栈操作
func (s *Stack) Pop() interface{} {
	s.Lock()
	defer s.Unlock()
	if s.length <= 0 {
		panic("堆栈弹出:索引超出范围")
	}
	t := s.data[s.length-1]
	s.data[s.length-1] = nil
	s.length--
	return t
}

// 返回栈顶元素
func (s *Stack) Top() interface{} {
	s.Lock()
	defer s.Unlock()

	if s.length <= 0 {
		panic("空栈")
	}

	return s.data[s.length-1]
}

// 返回当前栈元素个数
func (s *Stack) Count() int {
	s.Lock()
	defer s.Unlock()

	t := s.length

	return t
}

// 清空栈
func (s *Stack) Clear() {
	s.Lock()
	defer s.Unlock()

	s.data = make([]interface{}, 8)
	s.length = 0
	s.capacity = 8
}

// 栈是否为空
func (s *Stack) IsEmpty() bool {
	s.Lock()
	defer s.Unlock()
	b := s.length == 0
	return b
}
