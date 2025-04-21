package ui

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
)

func RunTray(onOpen, onQuit func()) {
	// Tentar criar um arquivo de log para o tray
	logFile, _ := os.OpenFile("tray_debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if logFile != nil {
		defer logFile.Close()
		fmt.Fprintln(logFile, "Iniciando systray")
	}

	systray.Run(func() {
		if logFile != nil {
			fmt.Fprintln(logFile, "Systray iniciado")
		}

		systray.SetTitle("AccessCtl")
		systray.SetTooltip("AccessCtl daemon")

		mOpen := systray.AddMenuItem("Abrir CLI", "Reabrir janela")
		mQuit := systray.AddMenuItem("Sair", "Encerra daemon")

		if logFile != nil {
			fmt.Fprintln(logFile, "Menus criados")
		}

		go func() {
			for {
				select {
				case <-mOpen.ClickedCh:
					if logFile != nil {
						fmt.Fprintln(logFile, "Menu Abrir clicado")
					}
					onOpen()
				case <-mQuit.ClickedCh:
					if logFile != nil {
						fmt.Fprintln(logFile, "Menu Sair clicado")
					}
					onQuit()
					systray.Quit()
					return
				}
			}
		}()

		if logFile != nil {
			fmt.Fprintln(logFile, "Goroutine de eventos iniciada")
		}
	}, func() {
		// cleanup se necessário
		if logFile != nil {
			fmt.Fprintln(logFile, "Systray finalizado")
		}
	})
}
