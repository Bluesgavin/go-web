package main

import (
	"fmt"
	"test-web/server"
)

func home(c *server.Context) {
	c.Respond("HOME PAGE")
}
func addBook(c *server.Context) {
	c.Respond("ADD A BOOK!")
}
func getBook(c *server.Context) {
	c.Respond("GET A BOOK!")
}
func updateBook(c *server.Context) {
	c.Respond("UPDATE A BOOK!")
}
func deleteBook(c *server.Context) {
	c.Respond("DELETE A BOOK!")
}

func  anyBook(c *server.Context){
	c.Respond("any route!")
}

func main() {
	server := server.NewServer(server.MetricFilterBuilder)
	// 主页
	server.Route("GET", "/", home)
	// 新增图书
	server.Route("GET", "/book/add", addBook)
	// 查找图书
	server.Route("GET", "/book/get", getBook)
	// 删除图书
	server.Route("GET", "/book/delete", deleteBook)
	// 更新图书
	server.Route("GET", "/book/update", updateBook)

	server.Route("GET", "/book/*", anyBook)

	if err := server.Start(":8081"); err != nil {
		panic(err)
	}
	fmt.Println("server started!")
}
