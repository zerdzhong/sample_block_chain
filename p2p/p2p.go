package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/go-libp2p"
	golog "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	gologging "github.com/whyrusleeping/go-logging"
	"io"
	"log"
	mathrand "math/rand"
	"os"
	"samplechain/blockchain"
	"samplechain/http"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

func MakeBasicHost(listenPort int, isSecureIO bool, randSeed int64) (host.Host, error) {
	var reader io.Reader
	if 0 == randSeed {
		reader = rand.Reader
	} else {
		reader = mathrand.New(mathrand.NewSource(randSeed))
	}

	golog.SetAllLoggers(gologging.INFO)

	privateKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, reader)

	if nil != err {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(privateKey)}

	if !isSecureIO {
		opts = append(opts, libp2p.NoEncryption())
	}

	basicHost, err := libp2p.New(context.Background(), opts...)

	if nil != err {
		return nil, err
	}

	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if isSecureIO {
		log.Printf("Now run \" -l %d -d %s --secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \" -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	basicHost.SetStreamHandler("/p2p/1.0.0", handleStream)

	return basicHost, nil

}

func DialToTarget(host host.Host, target string) {
	// The following code extracts target's peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(target)
	if err != nil {
		log.Fatalln(err)
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	host.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	log.Println("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol
	s, err := host.NewStream(context.Background(), peerid, "/p2p/1.0.0")
	if err != nil {
		log.Fatalln(err)
	}
	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	go writeData(rw)
	go readData(rw)

	select {} // hang forever
}

func handleStream(s net.Stream) {
	log.Println("Receive a new Stream")

	readWriter := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(readWriter)
	go writeData(readWriter)
}

func writeData(readWriter *bufio.ReadWriter) {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(blockchain.GetBlockChains())
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			mutex.Lock()
			readWriter.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			readWriter.Flush()
			mutex.Unlock()

		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)

		transaction := http.TransactionMessage{}
		json.Unmarshal([]byte(sendData), &transaction)

		fmt.Println(transaction)

		blockchain.Send(transaction.From, transaction.To, transaction.Value)

		if err != nil {
			log.Println(err)
		}

		mutex.Lock()
		bytes, err := json.Marshal(blockchain.GetBlockChains())
		readWriter.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		readWriter.Flush()
		mutex.Unlock()
	}
}

func readData(readWriter *bufio.ReadWriter) {
	for {
		str, err := readWriter.ReadString('\n')
		if nil != err {
			log.Fatal(err)
		}
		if "" == str || "\n" == str {
			return
		}

		chain := make([]blockchain.Block, 0)

		if err := json.Unmarshal([]byte(str), &chain); err != nil {
			log.Fatal(err)
		}

		mutex.Lock()

		if len(chain) > len(blockchain.GetBlockChains()) {

			blockchain.ReplaceBlockChain(chain)
			// Green console color: 	\x1b[32m
			// Reset console color: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

		mutex.Unlock()
	}
}
