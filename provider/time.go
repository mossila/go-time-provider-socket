package provider

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func (p *provider) clientHandler(c net.Conn) {
	defer c.Close()
	for {
		msg, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Printf("client logout with error = %s\n", err.Error())
			p.remove(c)
			break
		}
		fmt.Printf("client msg=%s", msg)
	}
}
func tprovider(d time.Duration, c chan string) {
	for x := range time.Tick(d) {
		fmt.Println(x.String())
		c <- x.String()
	}
}

func (p *provider) broadcast() {
	for {
		msg := <-p.msgChan
		if len(p.Clients) == 0 {
			continue
		}
		fmt.Printf("Broadcast to %d clients\n", len(p.Clients))

		for _, conn := range p.Clients {
			conn.Write([]byte(fmt.Sprintf("%s\n", msg)))
		}
	}

}
func (p *provider) add(c net.Conn) {
	p.Clients = append(p.Clients, c)
	c.Write([]byte("hello client\n"))
	go p.clientHandler(c)
}
func (p *provider) remove(c net.Conn) {
	for i, conn := range p.Clients {
		if c == conn {
			p.Clients = append(p.Clients[:i], p.Clients[i+1:]...)
			return
		}
	}
}

type provider struct {
	Clients []net.Conn
	msgChan chan string
}

//TimeProvider provide time on socket with define port
func TimeProvider(port string) {
	l, err := net.Listen("tcp", ":1234")
	c := make(chan string)
	conns := []net.Conn{}
	p := &provider{conns, c}

	go tprovider(time.Millisecond*1, c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is ready.")
	fmt.Println("connect with this command")
	fmt.Printf("telnet localhost %s \n", port)
	defer l.Close()
	go p.broadcast()
	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
			fmt.Println("client logout")
		}
		p.add(conn)
	}

}
