package stats

import (
    "fmt"
    "net"
    "time"
)

type PacketStats struct {
    Seq  int
    Addr net.Addr
    RTT  time.Duration
    Size int
    TTL  int
}

type Statistics struct {
    PacketsSent     int
    PacketsReceived int
    PacketsLost     int
    MinRTT          time.Duration
    MaxRTT          time.Duration
    AvgRTT          time.Duration
    TotalRTT        time.Duration
}

func NewStatistics() *Statistics {
    return &Statistics{
        MinRTT: time.Duration(1<<63 - 1), // Inicializar con el valor máximo posible
    }
}

func (s *Statistics) AddStat(stat *PacketStats) {
    s.PacketsSent++
    s.PacketsReceived++
    s.TotalRTT += stat.RTT

    if stat.RTT < s.MinRTT {
        s.MinRTT = stat.RTT
    }
    if stat.RTT > s.MaxRTT {
        s.MaxRTT = stat.RTT
    }
}

func (s *Statistics) PrintSummary() {
    s.PacketsLost = s.PacketsSent - s.PacketsReceived
    lossPercentage := float64(s.PacketsLost) / float64(s.PacketsSent) * 100

    if s.PacketsReceived > 0 {
        s.AvgRTT = s.TotalRTT / time.Duration(s.PacketsReceived)
    }

    fmt.Printf("\n--- Estadísticas de ping ---\n")
    fmt.Printf("%d paquetes transmitidos, %d recibidos, %.2f%% de pérdida de paquetes\n",
        s.PacketsSent, s.PacketsReceived, lossPercentage)
    fmt.Printf("RTT min/avg/max = %v/%v/%v\n", s.MinRTT, s.AvgRTT, s.MaxRTT)
}
