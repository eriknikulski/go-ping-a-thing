package main

import (
	"fmt"
	"github.com/sparrc/go-ping"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Ping CLI"
	app.Usage = "Let's you ping a host!"

	myFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "google.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "ping",
			Usage: "Pings an IP address",
			Flags: myFlags,
			Action: func(c *cli.Context) error {
				pinger, err := ping.NewPinger(c.String("host"))

				if err != nil {
					panic(err)
				}

				pinger.SetPrivileged(true)

				pinger.OnRecv = func(pkt *ping.Packet) {
					fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
						pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
				}

				pinger.OnFinish = func(stats *ping.Statistics) {
					fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
					fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
						stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
					fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
						stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
				}

				fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
				pinger.Run()

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
