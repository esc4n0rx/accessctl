package ui

import "github.com/rivo/tview"

type Actions struct {
	EditConfig func()
	Start      func()
	Stop       func()
	Status     func()
}

func RunCLI(actions Actions) {
	app := tview.NewApplication()
	menu := tview.NewList().
		AddItem("1) Carregar/Editar lista", "", '1', actions.EditConfig).
		AddItem("2) Iniciar Monitoramento", "", '2', actions.Start).
		AddItem("3) Parar Monitoramento", "", '3', actions.Stop).
		AddItem("4) Exibir Status", "", '4', actions.Status).
		AddItem("q) Sair", "", 'q', func() { app.Stop() })

	if err := app.SetRoot(menu, true).Run(); err != nil {
		panic(err)
	}
}
