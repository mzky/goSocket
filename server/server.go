package main

import (
	"flag"
	"goSocket/server/utils"
	"net"
)

func main() {
	p := flag.String("p", "9000", "listen socket port")
	//c := flag.Bool("w", true, "c  chang")
	t := flag.Int("t", 3, "timeout --second")

	flag.Parse()
	startServer(*p, *t)
}

func startServer(port string, timeout int) {
	//	setup a socket and listen the port

	netListen, err := net.Listen("tcp", "0.0.0.0:"+port)
	utils.CheckError(err)
	defer netListen.Close()
	utils.Log("Waiting for clients")

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		utils.Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn, timeout)
	}

	// you can run this part of code in Window System

	//netListen, err := net.Listen("tcp", "localhost:1024")
	//utils.CheckError(err)
	//defer netListen.Close()
	//utils.Log("Waiting for clients")
	//
	//for {
	//	conn, err := netListen.Accept()
	//	if err != nil {
	//		continue
	//	}
	//
	//	utils.Log(conn.RemoteAddr().String(), " tcp connect success")
	//	go handleConnection(conn, 3)
	//}
}

//handle the connection
func handleConnection(conn net.Conn, timeout int) {

	tmpBuffer := make([]byte, 0)

	buffer := make([]byte, 1024)
	messnager := make(chan byte)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			utils.Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = utils.Depack(append(tmpBuffer, buffer[:n]...))
		utils.Log("receive data string:", string(tmpBuffer))
		utils.TaskDeliver(tmpBuffer, conn)
		//start heartbeating
		go utils.HeartBeating(conn, messnager, timeout)
		//check if get message from client
		go utils.GravelChannel(tmpBuffer, messnager)

	}
	defer conn.Close()

}
