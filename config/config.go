package config

import (
    "flag"
    "time"
)

type Config struct {
    Count       int
    Interval    time.Duration
    Timeout     time.Duration
    Size        int
    TTL         int
    AllowICMPID bool
}

func NewConfig() *Config {
    cfg := &Config{}
    flag.IntVar(&cfg.Count, "c", 4, "número de paquetes a enviar")
    flag.DurationVar(&cfg.Interval, "i", time.Second, "intervalo entre paquetes")
    flag.DurationVar(&cfg.Timeout, "W", 5*time.Second, "tiempo de espera para respuestas")
    flag.IntVar(&cfg.Size, "s", 56, "tamaño de los datos en bytes")
    flag.IntVar(&cfg.TTL, "t", 64, "tiempo de vida (TTL)")
    flag.BoolVar(&cfg.AllowICMPID, "allow-icmp-id", false, "permitir ID ICMP personalizado")
    return cfg
}
