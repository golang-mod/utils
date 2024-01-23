package tree

import (
	"reflect"
)

type Tree map[string]interface{}

// Node 其他的结构体想要生成菜单树，直接实现这个接口
type Node interface {
	// GetId 获取id
	GetId() string
	// GetParentId 获取父id
	GetParentId() string
	// IsRoot 判断当前节点是否是顶层根节点
	IsRoot() bool
}
type Nodes []Node

// GenerateTree 自定义的结构体实现 Node 接口后调用此方法生成树结构
// nodes 需要生成树的节点
// menuTrees 生成成功后的树结构对象
func GenerateTree(nodes []Node) (trees []Tree) {
	// 定义顶层根和子节点
	var roots, childs []Node
	for _, v := range nodes {
		if v.IsRoot() {
			// 判断顶层根节点
			roots = append(roots, v)
		} else {
			childs = append(childs, v)
		}
	}
	for _, item := range roots {
		childTree := makeMap(item, "json")
		// 递归
		recursiveTree(childTree, childs)
		trees = append(trees, childTree)
	}
	return
}

// recursiveTree 递归生成树结构
// tree 递归的树对象
// nodes 递归的节点
func recursiveTree(tree Tree, nodes []Node) {
	var treeL []Tree
	for _, item := range nodes {
		if tree["id"] == item.GetParentId() {
			childTree := makeMap(item, "json")
			treeL = append(treeL, childTree)
			// 递归
			recursiveTree(childTree, nodes)
		}
	}
	tree["children"] = treeL
}

// makeMap 将结构体转成map
// node 待转结构体
// tree 转后map
func makeMap(node Node, key string) Tree {
	childTree := Tree{
		"id":        node.GetId(),
		"parent_id": node.GetParentId(),
	}
	var keys = reflect.TypeOf(node)
	var vals = reflect.ValueOf(node)
	for i := 0; i < keys.NumField(); i++ {
		k := keys.Field(i).Tag.Get(key)
		if k == "" {
			k = keys.Field(i).Name
		}
		v := vals.Field(i).Interface()
		childTree[k] = v
	}
	return childTree
}
