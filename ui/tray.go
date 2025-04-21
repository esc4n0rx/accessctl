package ui

import "github.com/getlantern/systray"

func RunTray(onOpen, onQuit func()) {
	systray.Run(func() {
		systray.SetTitle("AccessCtl")
		systray.SetTooltip("AccessCtl daemon")
		mOpen := systray.AddMenuItem("Abrir CLI", "Reabrir janela")
		mQuit := systray.AddMenuItem("Sair", "Encerra daemon")

		go func() {
			for {
				select {
				case <-mOpen.ClickedCh:
					onOpen()
				case <-mQuit.ClickedCh:
					onQuit()
					systray.Quit()
					return
				}
			}
		}()
	}, func() {
		// cleanup se necessário
	})
}
