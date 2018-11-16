package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Fprintf(w, `<table border="1"><tr><th>用户名</th><th>密码</th></tr>
		<tr> <td>`+ r.Form["username"][0]+ `</td>
		     <td>`+ r.Form["password"][0]+ `</td>
		</tr>
	  </table>`)
	}
}

func file(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("file.html")
	log.Println(t.Execute(w, nil))
}

func unknown(w http.ResponseWriter, r *http.Request){
	// http.Redirect(w,r,"/500",http.StatusBadGateway);
	log.Println("Bad Gateway!")
	w.WriteHeader(502)
}

func errorPage(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("500.html")
	t.Execute(w, nil)
}

func redirect(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("redirect.html")
	t.Execute(w, nil)
}


func getJS(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("redirect.js")
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", sayhelloName)       //默认路由
	http.HandleFunc("/login", login)         //表单处理
	http.HandleFunc("/file", file)         	//静态文件处理
	http.HandleFunc("/unknown", unknown)	//处理未知错误 - 502
	http.HandleFunc("/redirect", redirect) 	//通过JS处理重定向
	http.HandleFunc("/redirect.js",getJS)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}