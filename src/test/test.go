package main

import (
	"chord"
	"common"
	"fmt"
	"kademlia"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

func init1() {
	rand.Seed(1)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type keyval struct {
	key string
	val string
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	//s,_:=syscall.Socket(syscall.AF_INET,syscall.SOCK_STREAM,syscall.IPPROTO_TCP)
	//var reuse byte=1
	//syscall.Setsockopt(s,syscall.SOL_SOCKET,syscall.SO_REUSEADDR,&reuse,int32(unsafe.Sizeof(reuse)))
	init1()
	var nodes [160]*kademlia.Client

	localAddress := chord.GetLocalAddress()
	fmt.Println("local address: " + localAddress)
	port := 3000
	nodes[0] = common.NewNode(port)
	nodes[0].Run()
	nodes[0].Create()
	kvMap := make(map[string]string)
	var nodecnt = 1
	for i := 0; i < 5; i++ {
		//join 15 nodes
		for j := 0; j < 15; j++ {
			var index = i*15 + j + 1
			port++
			nodes[index] = common.NewNode(port)
			nodes[index].Run()
			if !nodes[index].Join(localAddress + ":" + strconv.Itoa(3000+5*i)) {
				log.Fatal("join failed at ", 3000+5*i)
			}
			time.Sleep(1 * time.Second)
			fmt.Println("port ", port, " joined at", 3000+5*i)
		}
		nodecnt += 15
		time.Sleep(4 * time.Second)
		for j := i * 5; j <= i*15+15; j++ {
			nodes[j].Dump()
		}
		//put 30 kv
		for j := 0; j < 30; j++ {
			k := RandStringRunes(30)
			v := RandStringRunes(30)
			kvMap[k] = v
			tmp := rand.Intn(nodecnt) + i*5
			fmt.Println("put ", k, " => ", v, " from node ", tmp)
			nodes[tmp].Put(k, v)
		}
		//get 20 kv and check correctness
		var keyList [20]string
		cnt := 0
		for k, v := range kvMap {
			if cnt == 20 {
				break
			}
			var tmp = rand.Intn(nodecnt) + i*5
			fetchedVal, success := nodes[tmp].Get(k)
			if !success {
				closest := nodes[tmp].Node_.RoutingTable.GetClosest(chord.HashString(k), kademlia.K)
				fmt.Println("closest: ", closest)
				for j := i * 5; j <= i*15+15; j++ {
					nodes[j].Dump()
				}
				fetchedVal, success = nodes[tmp].Get(k)
				log.Fatal("error:can't find key ", k, " from node ", tmp)
			}
			if !fetchedVal.Exist(v) {
				log.Fatal("actual: ", fetchedVal, " expected: ", v)
			}
			keyList[cnt] = k
			cnt++
		}
		//delete 150 kv
		//for j := 0; j < 150; j++ {
		//	delete(kvMap, keyList[j])
		//	nodes[rand.Intn(nodecnt)+i*5].Del(keyList[j])
		//}
		for j := i * 5; j <= i*15+15; j++ {
			nodes[j].Dump()
		}
		////force quit and join 5 nodes
		//for j := 0; j < 5; j++ {
		//	nodes[j+i*5+5].ForceQuit()
		//	time.Sleep(3 * time.Second)
		//	fmt.Println("force quit node ", j+i*5+5)
		//}
		//time.Sleep(10 * time.Second)
		//for j := i * 5; j <= i*30+30; j++ {
		//	if nodes[j].Node_.Listening {
		//		nodes[j].Dump()
		//	}
		//}
		//for j := 0; j < 5; j++ {
		//	nodes[j+i*5+5] = common.NewNode(j + i*5 + 5 + 3000)
		//	nodes[j+i*5+5].Run(&wg)
		//	if !nodes[j+i*5+5].Join(localAddress + ":" + strconv.Itoa(3000+i*5)) {
		//		log.Fatal("join failed")
		//	}
		//	fmt.Println("port ", j+i*5+5, " joined at ", 3000+i*5)
		//	time.Sleep(3 * time.Second)
		//}
		//quit 5 nodes
		for j := 0; j < 5; j++ {
			nodes[j+i*5].Quit()
			fmt.Println("quit ", j+i*5, " node")
			time.Sleep(3 * time.Second)
			//nodes[j+i*5+1].Dump()
		}
		nodecnt -= 5
		time.Sleep(10 * time.Second)
		for j := i*5 + 5; j <= i*15+15; j++ {
			nodes[j].Dump()
		}
		//put 30 kv
		for j := 0; j < 30; j++ {
			k := RandStringRunes(30)
			v := RandStringRunes(30)
			kvMap[k] = v
			tmp := rand.Intn(nodecnt) + i*5 + 5
			fmt.Println("put ", k, " => ", v, " from node ", tmp)
			nodes[tmp].Put(k, v)
			//nodes[tmp].Dump()
		}
		//get 20 kv and check correctness

		cnt = 0
		for k, v := range kvMap {
			if cnt == 20 {
				break
			}
			var tmp = rand.Intn(nodecnt) + i*5 + 5
			fetchedVal, success := nodes[tmp].Get(k)
			if !success {
				closest := nodes[tmp].Node_.RoutingTable.GetClosest(chord.HashString(k), kademlia.K)
				fmt.Println("closest: ", closest)
				for j := i*5 + 5; j <= i*15+15; j++ {
					nodes[j].Dump()
				}
				fetchedVal, success = nodes[tmp].Get(k)
				log.Fatal("error:can't find key ", k, " from node ", tmp)

			}
			if !fetchedVal.Exist(v) {
				log.Fatal("actual: ", fetchedVal, " expected: ", v)
			}
			keyList[cnt] = k
			cnt++
		}
		//delete 150 kv
		//for j := 0; j < 150; j++ {
		//	delete(kvMap, keyList[j])
		//	nodes[rand.Intn(nodecnt)+i*5+5].Del(keyList[j])
		//}
		time.Sleep(4 * time.Second)
	}
}
