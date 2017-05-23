package main

import (
	"github.com/astaxie/beego/context"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"

	"fmt"
)

//定义一个自己的中间件，这里将beego的context注入
func myContext() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		ctx := context.Context{Request: req, ResponseWriter: res}
		ctx.Input = context.NewInput(req)
		ctx.Output = context.NewOutput()
		c.Map(ctx)
	}
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer()) //注入中间件（渲染JSON和HTML模板的处理器中间件）
	m.Use(myContext())       //注入自己写的中间件

	m.Use(func(c martini.Context) {
		fmt.Println("before a request")
		c.Next() //Next方法之后最后处理
		fmt.Println("after a request")
	})

	//普通的GET方式路由
	m.Get("/", func() string {
		return "hello world!"
	})

	//路由分组
	m.Group("/books", func(r martini.Router) {
		r.Get("/list", getBooks)
		r.Post("/add", getBooks)
		r.Delete("/delete", getBooks)
	})

	//我们以中间件的方式来注入一个Handler
	m.Use(MyHeader(m))

	m.RunOnAddr(":8080") //运行程序监听端口
}

func getBooks() string {
	return "books"
}

//中间件Handler
func MyHeader(m *martini.ClassicMartini) martini.Handler {
	return func() {
		m.Group("/app", func(r martini.Router) {
			my := new(App)
			r.Get("/index", my.Index)
			r.Get("/test/:aa", my.Test)
		})
	}
}

//应用的处理
type App struct{}

func (this *App) Index(r render.Render, ctx context.Context) {
	fmt.Println(ctx.Input.Query("action"))
	ctx.WriteString("你好世界")
}

func (this *App) Test(r render.Render, params martini.Params, req *http.Request) {
	fmt.Println(params)

	parm := make(map[string]interface{})
	if t, ok := params["aa"]; ok {
		parm["aa"] = t
	}
	req.ParseForm()
	fmt.Println(parm, req.Form)

	r.Text(200, "----")
}
