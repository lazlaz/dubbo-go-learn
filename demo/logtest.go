package main

import "log"

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	log.Println("飞雪无情的博客:", "http://www.flysnow.org")
	log.Printf("飞雪无情的微信公众号：%s\n", "flysnow_org")
	log.Fatal("飞雪无情的博客:", "http://www.flysnow.org")
}
