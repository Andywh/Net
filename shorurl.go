// // server.go

// package main

// import (
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	//"github.com/go-martini/martini"
// )

// func transfer(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("method:", r.Method)
// 	if r.Method == "GET" {
// 		t, _ := template.ParseFiles("index.gtpl")
// 		log.Println(t.Execute(w, nil))
// 	} else {
// 		//fmt.Println("method:", r.Method)
// 		fmt.Println("username:", r.Form["longurl"])
// 		//fmt.Println("password:", r.Form["shorturl"])
// 	}
// }

// func main() {
// 	// m := martini.Classic()
// 	// m.Get("/", transfer)
// 	// m.Run()
// 	http.HandleFunc("/", transfer)
// 	err := http.ListenAndServe(":9000", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe:", err)
// 	}
// }

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/alphazero/Go-Redis"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func int64ToString(i int64) string {
	s := "0123456789ABCEDF"
	var str string
	for i > 0 {
		temp := s[i%16]
		str = string(temp) + str
		i = i / 16
	}
	return str
}

func route(w http.ResponseWriter, r *http.Request) {
	//isFound := false
	r.ParseForm()
	fmt.Println("route", r.URL.Path)
	fmt.Println("equal or not: ", r.URL.Path == "/")
	if r.URL.Path == "/" {
		fmt.Println("method:", r.Method) //获取请求的方法
		//r.ParseForm()
		if r.Method == "GET" {
			t, _ := template.ParseFiles("index.gtpl")
			log.Println(t.Execute(w, nil))
		} else {
			//请求的是登陆数据，那么执行登陆的逻辑判断
			//fmt.Println("longurl:", r.Form["longurl"])
			//fmt.Println("shorturl:", r.Form["shorturl"])
			//r.Form["url"] = append(r.Form["shorturl"], "123")

			//t, _ := template.ParseFiles("index.gtpl")

			//if _, exist := r.Form["longurl"]; !exist {
			//	r.Form["longurl"] := make(map[string][]string)
			//}

			spec := redis.DefaultSpec().Db(0).Password("")
			client, err := redis.NewSynchClientWithSpec(spec)
			if err != nil {
				fmt.Println("Connect redis server fail")
				return
			}
			//var getValue []byte
			fmt.Printf("test %s\n", r.Form["longurl"][0])
			getValue, err := client.Get(r.Form["longurl"][0])
			fmt.Println("getValue == nil", getValue == nil)
			fmt.Println("err != nil", err != nil)
			fmt.Println(err)
			if getValue == nil {
				last, _ := client.Dbsize()
				last++
				str := int64ToString(last)
				fmt.Println(str)
				fmt.Println("err")
				getValue = []byte(str)
				//key := []byte(r.Form["longurl"][0])
				fmt.Printf("key: %s", r.Form["longurl"][0])
				client.Set(r.Form["longurl"][0], getValue)

				sp := redis.DefaultSpec().Db(1).Password("")
				cli, err := redis.NewSynchClientWithSpec(sp)
				if err != nil {
					fmt.Println("Connect redis server fail")
					return
				}
				v := []byte(r.Form["longurl"][0])
				cli.Set(str, v)
			}
			fmt.Println("shorturl:", string(getValue))
			val, _ := client.Get("test1")
			fmt.Println(string(val))

			fmt.Fprintf(w, "<html>"+
				"<head>"+
				"<title></title>"+
				"</head>"+
				"<body>"+
				"<form action=\"/\" method=\"post\">"+
				"请输入长网址:</br>"+
				"<input type=\"text\" name=\"longurl\">"+
				"<input type=\"submit\" value=\"生成短地址\"></br>"+
				"新网址:</br>"+
				"http://116.62.60.145/%s<output type=\"text\" name=\"shorturl\"></br>"+
				"原网址:</br>"+
				"%s"+
				"<output type=\"text\" name=\"shorturl\">"+
				"</form>"+
				"</body>"+
				"</html>", string(getValue), r.Form["longurl"][0])
		}
	} else {
		spec := redis.DefaultSpec().Db(1).Password("")
		client, err := redis.NewSynchClientWithSpec(spec)
		if err != nil {
			fmt.Println("Connect redis server fail")
			return
		}
		fmt.Println("before trimprefix", r.URL.Path)
		shorturl := strings.TrimPrefix(r.URL.Path, "/")
		fmt.Println("after trimprefix", shorturl)
		getValue, _ := client.Get(shorturl)
		u := string(getValue)
		fmt.Println(u)
		if getValue != nil {
			url := string(getValue)

			fmt.Println("this url: ", url)
			//w.Header().Set("Location", url)
			http.Redirect(w, r, url, 302)
			return
			//http.Redirect(w, r, url, http.StatusFound)
		} else {
			fmt.Fprint(w, "404 Page Not Found!")
		}
	}
}

func transfer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.URL.Path)
	fmt.Fprintf(w, r.URL.Path)
}

func main() {
	//http.HandleFunc("/", sayhelloName)       //设置访问的路由
	http.HandleFunc("/", route) //设置访问的路由
	//http.HandleFunc("/.*", transfer)
	err := http.ListenAndServe(":9000", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
