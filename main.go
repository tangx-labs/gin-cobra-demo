package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r1 := &Router{
		Path:    "/root",
		Method:  http.MethodGet,
		Handler: pingHandler,
	}

	// r1.Run()

	r2 := &Router{
		Path:    "/r2",
		Method:  http.MethodPost,
		Handler: pingHandler,
	}

	r3 := &Router{
		Path:    "/r3",
		Method:  http.MethodGet,
		Handler: pingHandler,
	}

	r1.AddRouters(r2, r3)
	r2.AddRouters(r3)

	r2.Run()

}

func pingHandler(c *gin.Context) {
	c.String(200, "pong")
}

type Router struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc

	parent      *Router
	children    []*Router
	engine      *gin.Engine
	routerGroup *gin.RouterGroup
}

func (r *Router) AddRouters(routers ...*Router) {
	for i, x := range routers {
		if r == x {
			panic("自己不能成为自己的子路由")
		}

		routers[i].parent = r

		r.children = append(r.children, x)
	}
}

func (r *Router) Run() {
	if r.parent != nil {
		r.parent.Run()
	}

	if r.engine == nil {
		r.engine = gin.Default()
	}

	r.register()

	fmt.Println("当前路径: ", r.Path)
	r.engine.Run(":8081")
}

// register。 注册子路由, 只能 root 发起调用。
func (r *Router) register() {

	// root 从 engine 开始注册
	if r.parent == nil {
		r.routerGroup = r.engine.Group(r.Path)
		r.routerGroup.Handle(r.Method, "/", r.Handler)
	}

	// 注册子路由
	for _, child := range r.children {

		// 子路由注册当前到节点下
		child.routerGroup = r.routerGroup.Group(child.Path)
		child.routerGroup.Handle(child.Method, "", child.Handler)

		// 递归, 进行子路由注册
		child.register()
	}
}
