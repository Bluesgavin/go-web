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

func anyBook(c *server.Context) {
	c.Respond("any route!")
}

func main() {
	s := server.NewServer(server.MetricFilterBuilder)
	// 主页
	s.Route("GET", "/", home)
	// 新增图书
	s.Route("GET", "/book/add", addBook)
	// 查找图书
	s.Route("GET", "/book/get", getBook)
	// 删除图书
	s.Route("GET", "/book/delete", deleteBook)
	// 更新图书
	s.Route("GET", "/book/update", updateBook)

	s.Route("GET", "/book/*", anyBook)

	staticRoute := server.NewStaticHandler("public", "/public", server.WithMoreExtension(map[string]string{"mp3": "audio/mp3"}), server.WithFileCache(1<<20, 100))
	s.Route("GET", "/public/*", staticRoute.Handle)

	if err := s.Start(":8081"); err != nil {
		panic(err)
	}
	fmt.Println("server started!")
}
