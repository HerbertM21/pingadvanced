package pinger

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"

	"pingadvanced/config"
	"pingadvanced/stats"
)

type Pinger struct {
	dest   string
	addr   *net.IPAddr
	conn   *icmp.PacketConn
	config *config.Config
}

func NewPinger(dest string, cfg *config.Config) (*Pinger, error) {
	addr, err := net.ResolveIPAddr("ip4", dest)
	if err != nil {
		return nil, fmt.Errorf("error resolviendo %s: %v", dest, err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("error al escuchar: %v", err)
	}

	return &Pinger{
		dest:   dest,
		addr:   addr,
		conn:   conn,
		config: cfg,
	}, nil
}

func (p *Pinger) Run(statsChan chan<- *stats.PacketStats) {
	defer p.conn.Close()

	for i := 0; i < p.config.Count; i++ {
		if i > 0 {
			time.Sleep(p.config.Interval)
		}

		if err := p.sendPing(i); err != nil {
			fmt.Printf("Error al enviar ping: %v\n", err)
			continue
		}

		if s, err := p.receivePong(); err != nil {
			fmt.Printf("Error al recibir pong: %v\n", err)
		} else {
			statsChan <- s
		}
	}
}

func (p *Pinger) sendPing(seq int) error {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: make([]byte, p.config.Size),
		},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return fmt.Errorf("error al crear el mensaje: %v", err)
	}

	if _, err := p.conn.WriteTo(msgBytes, p.addr); err != nil {
		return fmt.Errorf("error al enviar el paquete: %v", err)
	}

	return nil
}

func (p *Pinger) receivePong() (*stats.PacketStats, error) {
	start := time.Now()
	reply := make([]byte, 1500)

	err := p.conn.SetReadDeadline(time.Now().Add(p.config.Timeout))
	if err != nil {
		return nil, fmt.Errorf("error al establecer el tiempo de espera: %v", err)
	}

	n, peer, err := p.conn.ReadFrom(reply)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	duration := time.Since(start)

	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		return nil, fmt.Errorf("error al analizar la respuesta: %v", err)
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		echo := rm.Body.(*icmp.Echo)
		return &stats.PacketStats{
			Seq:  echo.Seq,
			Addr: peer,
			RTT:  duration, // RTT es
			Size: n,
			TTL:  p.config.TTL, // Esto es una simplificación, idealmente obtendrías el TTL real del paquete
		}, nil
	default:
		return nil, fmt.Errorf("recibido %+v desde %v", rm, peer)
	}
}

func (p *Pinger) IPAddr() net.Addr {
	return p.addr
}
