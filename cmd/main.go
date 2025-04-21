package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"accessctl/config"
	"accessctl/logger"
	"accessctl/monitor"
	"accessctl/ui"
)

func main() {
	// 1. Carrega config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Println("Erro ao carregar config:", err)
		os.Exit(1)
	}

	// 2. Logger
	log, err := logger.New("accessctl.log")
	if err != nil {
		fmt.Println("Erro inicializando logger:", err)
		os.Exit(1)
	}

	// 3. Canal de controle
	stop := make(chan struct{})
	var wg sync.WaitGroup

	actions := ui.Actions{
		EditConfig: func() {
			// reimplemente editor externo ou recarregue arquivo
			log.Info.Println("EditConfig acionado")
		},
		Start: func() {
			log.Info.Println("Iniciando monitoramento")
			wg.Add(2)
			go monitor.ProcessMonitor(cfg, log, stop, &wg)
			go monitor.NetworkMonitor(cfg, log, stop, &wg)
		},
		Stop: func() {
			log.Info.Println("Parando monitoramento")
			close(stop)
		},
		Status: func() {
			log.Info.Println("Status: rodando")
		},
	}

	// 4. Inicia tray + CLI
	go ui.RunTray(func() { go ui.RunCLI(actions) }, func() { close(stop) })
	ui.RunCLI(actions)

	// 5. Aguarda sinal de término
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	close(stop)
	wg.Wait()
	log.Info.Println("Daemon finalizado")
}
