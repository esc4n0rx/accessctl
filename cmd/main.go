package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"accessctl/config"
	"accessctl/logger"
	"accessctl/monitor"
	"accessctl/ui"
)

func main() {
	// Cria um arquivo de log temporário na área de trabalho
	logFile, err := os.Create("accessctl_debug.log")
	if err == nil {
		defer logFile.Close()
		fmt.Fprintf(logFile, "Iniciando aplicação: %s\n", time.Now().String())
	}

	// 1. Carrega config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Println("Erro ao carregar config:", err)
		if logFile != nil {
			fmt.Fprintf(logFile, "Erro ao carregar config: %v\n", err)
		}
		os.Exit(1)
	}

	if logFile != nil {
		fmt.Fprintf(logFile, "Config carregada com sucesso: keywords=%d, domains=%d\n",
			len(cfg.Keywords), len(cfg.Domains))
	}

	// 2. Logger
	log, err := logger.New("accessctl.log")
	if err != nil {
		fmt.Println("Erro inicializando logger:", err)
		if logFile != nil {
			fmt.Fprintf(logFile, "Erro inicializando logger: %v\n", err)
		}
		os.Exit(1)
	}

	// Registra início nos logs
	log.Info.Println("Aplicação iniciada")
	if logFile != nil {
		fmt.Fprintf(logFile, "Logger iniciado com sucesso\n")
	}

	// 3. Canal de controle
	stop := make(chan struct{})
	var wg sync.WaitGroup

	actions := ui.Actions{
		EditConfig: func() {
			log.Info.Println("EditConfig acionado")
			if logFile != nil {
				fmt.Fprintf(logFile, "EditConfig acionado\n")
			}
			// reimplemente editor externo ou recarregue arquivo
		},
		Start: func() {
			log.Info.Println("Iniciando monitoramento")
			if logFile != nil {
				fmt.Fprintf(logFile, "Iniciando monitoramento\n")
			}
			log.Info.Println("Iniciando monitoramento automaticamente")
			wg.Add(2)
			go monitor.ProcessMonitor(cfg, log, stop, &wg)
			go monitor.NetworkMonitor(cfg, log, stop, &wg)
		},
		Stop: func() {
			log.Info.Println("Parando monitoramento")
			if logFile != nil {
				fmt.Fprintf(logFile, "Parando monitoramento\n")
			}
			close(stop)
		},
		Status: func() {
			log.Info.Println("Status: rodando")
			if logFile != nil {
				fmt.Fprintf(logFile, "Status: rodando\n")
			}
		},
	}

	// 4. Inicia tray + CLI
	if logFile != nil {
		fmt.Fprintf(logFile, "Iniciando interface\n")
	}

	// Inicia CLI em uma goroutine separada
	go func() {
		if logFile != nil {
			fmt.Fprintf(logFile, "Tentando iniciar CLI\n")
		}
		ui.RunCLI(actions)
		if logFile != nil {
			fmt.Fprintf(logFile, "CLI finalizado\n")
		}
	}()

	// Inicia interface de bandeja do sistema
	if logFile != nil {
		fmt.Fprintf(logFile, "Tentando iniciar Tray\n")
	}

	go ui.RunTray(func() {
		if logFile != nil {
			fmt.Fprintf(logFile, "Tray: ação Abrir CLI\n")
		}
		go ui.RunCLI(actions)
	}, func() {
		if logFile != nil {
			fmt.Fprintf(logFile, "Tray: ação Sair\n")
		}
		close(stop)
	})

	// 5. Aguarda sinal de término
	if logFile != nil {
		fmt.Fprintf(logFile, "Aguardando sinal de término\n")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs

	if logFile != nil {
		fmt.Fprintf(logFile, "Sinal de término recebido\n")
	}

	close(stop)
	wg.Wait()
	log.Info.Println("Daemon finalizado")

	if logFile != nil {
		fmt.Fprintf(logFile, "Aplicação finalizada: %s\n", time.Now().String())
	}
}
