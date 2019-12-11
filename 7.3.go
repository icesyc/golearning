package main

import (
	"fmt"
	"bytes"
	"strconv"
)


func main() {
	arr := []int{1, 20, 5, 8, 3, 7}
	var t *Tree
	t = AddAll(arr)
	fmt.Printf("%s\n", t)
}

type Tree struct{
	value int
	left, right *Tree
}

func Sort(values []int){
	var t *Tree
	for _, value := range values {
		t = Add(t, value)
	}
	appendValues(values[:0], t)
}

func AddAll(values []int) (t *Tree) {
	for _, value := range values {
		t = Add(t, value)
	}
	return t 
}

//这里使用非递归算法，递归算法表达其实要简单很多
func (t *Tree) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	var stack []*Tree
	node := t
	first := true
	stackPoped := false
	for node != nil || len(stack) > 0 {
		//1.左结点不为空，入栈
		//节点如果是刚刚出栈的，不需要处理左节点
		if !stackPoped && node.left != nil {
			stack = append(stack, node)
			node = node.left
			continue
		}

		//2.左结点为空，输出当前节点value
		value := strconv.Itoa(node.value)
		//不是根节点
		if first {
			first = false
		}else{
			buf.WriteByte(',')
		}
		buf.WriteString(value)

		//处理右节点
		node = node.right
		//3.右节点不为空, 继续[1]
		if node != nil {
			//继续[1]时要处理左节点
			stackPoped = false
			continue
		}
		//4.右节点为空, 出栈父节点, 继续循环
		if len(stack) > 0 {
			stackPoped = true
			node = stack[len(stack)-1]
			stack = stack[0:len(stack)-1]
		}
	}

	buf.WriteByte('}')
	return buf.String()
}

func appendValues(values []int, t *Tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func Add(t *Tree, value int) *Tree {
	if t == nil {
		t = new(Tree)
		t.value = value	
		return t
	}
	if value < t.value {
		t.left = Add(t.left, value)
	}else {
		t.right = Add(t.right, value)
	}
	return t
}

