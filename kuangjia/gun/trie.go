package gun

import "strings"

//不太会实现动态路由（通过网上看文档人家使用了前缀树来进行的）
//树节点
type node struct {
	pattern  string  //带匹配路由
	part     string  //路由中的部分
	children []*node //子节点（就是个结构体嵌套，同node理解）
	isWild   bool    //是否匹配精确？如含:or*时为true
}
//前缀树即每个节点都含一个字符，从根节点开始寻找如相同，则寻找其子节点，
//如找到，找其子节点以此类推

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	//for循环遍历查找子节点
	for _, child := range n.children {
		//如果子节点匹配
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)//创建结构体map
	//for循环遍历添加map
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入
func (n *node) insert(pattern string, parts []string, height int) {
	//如果成熟相同，则在改成加入一个新路由
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)//调用匹配的路由

	//如果子路由为空，则
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)//插入
}

//寻找节点数
func (n *node) search(parts []string, height int) *node {
	//如果长度相同或者能匹配到*
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//如果路由为空
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)            //匹配成功的节点

	for _, child := range children {
		result := child.search(parts, height+1)  //寻找节点数
		if result != nil {
			return result
		}
	}
	return nil
}
