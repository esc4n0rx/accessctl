package monitor

import (
	"fmt"
	"strings"
	"sync"

	"accessctl/config"
	"accessctl/controller"
	"accessctl/logger"

	"github.com/moutend/go-hook/pkg/hook"
	"github.com/moutend/go-hook/pkg/hook/winhook"
)

func ProcessMonitor(cfg *config.Config, log *logger.Logger, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	buffers := make(map[uint32]string)
	var mu sync.Mutex

	callback := hook.EventCallback(func(e hook.Event) {
		if e.Kind != hook.KindKeyboard {
			return
		}
		r := winhook.VKeycodeToString(e.VKeyCode)
		if r == "" {
			return
		}
		pid := e.ProcessID

		mu.Lock()
		buffers[pid] += r
		// manter no máximo 100 caracteres
		if len(buffers[pid]) > 100 {
			buffers[pid] = buffers[pid][len(buffers[pid])-100:]
		}
		// busca por palavras-chave
		for _, kw := range cfg.Keywords {
			if strings.Contains(strings.ToLower(buffers[pid]), strings.ToLower(kw)) {
				log.Info.Println(fmt.Sprintf("Keyword '%s' detectada no PID %d", kw, pid))
				if err := controller.TerminateProcess(int(pid)); err != nil {
					log.Error.Println("falha ao terminar processo:", err)
				} else {
					log.Info.Println(fmt.Sprintf("Processo %d encerrado", pid))
				}
				buffers[pid] = ""
				break
			}
		}
		mu.Unlock()
	})

	// registra hook global de teclado
	hook.Register(hook.HookTypeKeyboard, []hook.VKeyCode{}, callback)
	go hook.Run()

	// aguarda sinal de parada
	<-stop
	hook.Unregister(hook.HookTypeKeyboard, []hook.VKeyCode{}, callback)
}
