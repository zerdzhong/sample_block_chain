package networking

import (
	"log"
	"net"
	"os"
	. "samplechain/blockchain"
)

var bcServer chan BlockChain

func StartServer() {
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))

	if nil != err {
		log.Fatal(err)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
}
