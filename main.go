package main

import (
    "flag"
    "fmt"
    "os"
    "time"

    "pingadvanced/config"
    "pingadvanced/pinger"
    "pingadvanced/stats"
)

func main() {
    cfg := config.NewConfig()
    flag.Parse()

    if len(flag.Args()) < 1 {
        fmt.Println("Uso: go run main.go [opciones] <destino>")
        flag.PrintDefaults()
        os.Exit(1)
    }

    dest := flag.Args()[0]
    p, err := pinger.NewPinger(dest, cfg)
    if err != nil {
        fmt.Printf("Error al crear el pinger: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("PING %s (%s):\n", dest, p.IPAddr())

    statsChan := make(chan *stats.PacketStats, cfg.Count)
    go p.Run(statsChan)

    statistics := stats.NewStatistics()

    for i := 0; i < cfg.Count; i++ {
        select {
        case s := <-statsChan:
            statistics.AddStat(s)
            fmt.Printf("%d bytes desde %v: icmp_seq=%d ttl=%d tiempo=%v\n",
                s.Size, s.Addr, s.Seq, s.TTL, s.RTT)
        case <-time.After(cfg.Timeout):
            fmt.Printf("Tiempo de espera agotado para la secuencia %d\n", i)
        }
    }

    statistics.PrintSummary()
}
