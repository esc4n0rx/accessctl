package monitor

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	ps "github.com/mitchellh/go-ps"

	"accessctl/config"
	"accessctl/controller"
	"accessctl/logger"
)

func NetworkMonitor(cfg *config.Config, log *logger.Logger, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	iface := "any" // Windows: usa "any" ou especifique o nome da interface
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Error.Println("PCAP error:", err)
		return
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("udp port 53"); err != nil {
		log.Error.Println("BPF filter error:", err)
		return
	}

	src := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <-stop:
			return
		case packet := <-src.Packets():
			dnsLayer := packet.Layer(layers.LayerTypeDNS)
			if dnsLayer == nil {
				continue
			}
			dns := dnsLayer.(*layers.DNS)
			for _, q := range dns.Questions {
				name := string(q.Name)
				for _, domain := range cfg.Domains {
					if strings.Contains(strings.ToLower(name), strings.ToLower(domain)) {
						log.Info.Println(fmt.Sprintf("DNS bloqueado: %s", name))
						// encerra navegadores conhecidos
						killBrowsers(log)
					}
				}
			}
		}
	}
}

func killBrowsers(log *logger.Logger) {
	browsers := []string{"chrome.exe", "firefox.exe", "msedge.exe"}
	procs, _ := ps.Processes()
	for _, p := range procs {
		for _, b := range browsers {
			if strings.EqualFold(p.Executable(), b) {
				if err := controller.TerminateProcess(p.Pid()); err != nil {
					log.Error.Println("erro encerrando", b, p.Pid(), err)
				} else {
					log.Info.Println("encerrou processo", b, p.Pid())
				}
			}
		}
	}
}
