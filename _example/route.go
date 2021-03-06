package main

import (
	gp "./gopeer"
	"fmt"
	"time"
	"crypto/rsa"
)

const (
	TITLE_MESSAGE = "MESSAGE"
	NODE1_ADDRESS = ":8080"
	NODE2_ADDRESS = ":9090"
)

/*
	client1 -> client2
	client1 <-> node1 <-> client2 <-> node2
*/

func main() {
	client1 := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	client2 := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	client1.SetHandle(handleFunc)
	client2.SetHandle(handleFunc)

	fmt.Println(gp.HashPublic(client1.Public()))

	node1 := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	node2 := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	node1.SetHandle(handleFunc)
	node2.SetHandle(handleFunc)

	go gp.NewNode(NODE1_ADDRESS, node1).Run()
	go gp.NewNode(NODE2_ADDRESS, node2).Run()
	time.Sleep(500 * time.Millisecond)

	client1.Connect(NODE1_ADDRESS)
	client2.Connect(NODE1_ADDRESS)

	client2.Connect(NODE2_ADDRESS)

	route := []*rsa.PublicKey{
		node1.Public(),
		node2.Public(),
	}

	pseudoSender := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	for i := 0; i < 10; i++ {
		res, err := client1.Send(
			client2.Public(), 
			gp.NewPackage(TITLE_MESSAGE, fmt.Sprintf("hello, world! [%d]", i)),
			route,
			pseudoSender,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(res)
	}
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	gp.Handle(TITLE_MESSAGE, client, pack, getMessage)
}

func getMessage(client *gp.Client, pack *gp.Package) (set string) {
	public := gp.ParsePublic(pack.Head.Sender)
	fmt.Printf("[%s] => '%s'\n", gp.HashPublic(public), pack.Body.Data)
	return "ok"
}
