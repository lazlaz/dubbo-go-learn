package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/laz/dubbo-go/kownledge/fourteen/addressbook"
	"io/ioutil"
	"log"
)

func main() {
	// 自定义AddressBook内容
	book := &pb.AddressBook{
		People: []*pb.Person{
			&pb.Person{
				Id:    1,
				Name:  "zyq",
				Email: "77@qq.com",
				Phones: []*pb.Person_PhoneNumber{
					&pb.Person_PhoneNumber{
						Number: "11111",
						Type:   pb.Person_MOBILE,
					},
					&pb.Person_PhoneNumber{
						Number: "22222",
						Type:   pb.Person_HOME,
					},
				},
			},
		},
	}
	fmt.Println("book : ", book)

	fname := "address.dat"
	// 将book进行序列化
	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}
	// 将序列化的内容写入文件
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}

	// 读取写入的二进制数据
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	// 定义一个空的结构体
	book2 := &pb.AddressBook{}
	// 将从文件中读取的二进制进行反序列化
	if err := proto.Unmarshal(in, book2); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	fmt.Println("book2: ", book2)
}
