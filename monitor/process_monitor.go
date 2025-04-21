package monitor

import (
	"fmt"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"accessctl/config"
	"accessctl/controller"
	"accessctl/logger"

	"github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows"
)

// Windows API direto para funções não disponíveis
var (
	user32                       = windows.NewLazyDLL("user32.dll")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowText            = user32.NewProc("GetWindowTextW")
	procGetWindowTextLength      = user32.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
)

// Typedef para a callback EnumWindows
type EnumWindowsProc func(hwnd windows.Handle, lParam uintptr) uintptr

// Função para enumerar janelas
func enumWindows(enumFunc EnumWindowsProc, lParam uintptr) (err error) {
	r1, _, e1 := syscall.SyscallN(procEnumWindows.Addr(),
		syscall.NewCallback(enumFunc), uintptr(lParam))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// Obtém o ID do processo de uma janela
func getWindowThreadProcessId(hwnd windows.Handle) uint32 {
	var processID uint32
	syscall.SyscallN(procGetWindowThreadProcessId.Addr(),
		uintptr(hwnd), uintptr(unsafe.Pointer(&processID)))
	return processID
}

// Verifica se uma janela é visível
func isWindowVisible(hwnd windows.Handle) bool {
	ret, _, _ := syscall.SyscallN(procIsWindowVisible.Addr(), uintptr(hwnd))
	return ret != 0
}

// Obtém o tamanho do texto da janela
func getWindowTextLength(hwnd windows.Handle) int {
	ret, _, _ := syscall.SyscallN(procGetWindowTextLength.Addr(), uintptr(hwnd))
	return int(ret)
}

// Obtém o texto da janela
func getWindowText(hwnd windows.Handle, str *uint16, maxCount int) int {
	ret, _, _ := syscall.SyscallN(procGetWindowText.Addr(),
		uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	return int(ret)
}

// Obtém o título da janela de um processo
func getWindowTitle(pid uint32) string {
	var title string

	_ = enumWindows(func(hwnd windows.Handle, lParam uintptr) uintptr {
		currentPid := getWindowThreadProcessId(hwnd)

		if currentPid == pid {
			if isWindowVisible(hwnd) {
				length := getWindowTextLength(hwnd)
				if length > 0 {
					buf := make([]uint16, length+1)
					getWindowText(hwnd, &buf[0], length+1)
					title = syscall.UTF16ToString(buf)
					return 0 // Parar enumeração
				}
			}
		}
		return 1 // Continuar enumeração
	}, 0)

	return title
}

func ProcessMonitor(cfg *config.Config, log *logger.Logger, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Info.Println("Iniciando monitoramento de processos")

	// Buffer de teclas por processo
	buffers := make(map[uint32]string)
	var mu sync.Mutex

	// Loop de verificação periódica de processos ativos
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			log.Info.Println("Parando monitoramento de processos")
			return
		case <-ticker.C:
			// Obter lista de processos
			processes, err := ps.Processes()
			if err != nil {
				log.Error.Println("Erro ao listar processos:", err)
				continue
			}

			// Verificar títulos de janelas de processos ativos
			for _, proc := range processes {
				if proc == nil {
					continue
				}

				pid := uint32(proc.Pid())
				name := strings.ToLower(proc.Executable())

				// Verifica se é um navegador ou aplicativo de interesse
				if isBrowser(name) || isInterestingApp(name) {
					// Pega o título da janela (pode conter texto sensível)
					title := getWindowTitle(pid)

					mu.Lock()
					// Adiciona o título ao buffer desse processo
					if title != "" {
						buffers[pid] += " " + title

						// Limita o tamanho do buffer
						if len(buffers[pid]) > 1000 {
							buffers[pid] = buffers[pid][len(buffers[pid])-1000:]
						}

						// Verifica palavras-chave
						for _, kw := range cfg.Keywords {
							if strings.Contains(strings.ToLower(buffers[pid]), strings.ToLower(kw)) {
								log.Info.Println(fmt.Sprintf("Keyword '%s' detectada no PID %d (%s)", kw, pid, name))
								if err := controller.TerminateProcess(int(pid)); err != nil {
									log.Error.Println("falha ao terminar processo:", err)
								} else {
									log.Info.Println(fmt.Sprintf("Processo %d encerrado", pid))
								}
								buffers[pid] = ""
								break
							}
						}
					}
					mu.Unlock()
				}
			}
		}
	}
}

// Verifica se o nome do processo é um navegador conhecido
func isBrowser(name string) bool {
	browsers := []string{"chrome.exe", "firefox.exe", "msedge.exe", "opera.exe", "brave.exe", "iexplore.exe"}
	for _, b := range browsers {
		if strings.EqualFold(name, b) {
			return true
		}
	}
	return false
}

// Verifica se é um aplicativo de interesse para monitorar
func isInterestingApp(name string) bool {
	apps := []string{"notepad.exe", "wordpad.exe", "word.exe", "excel.exe", "powerpnt.exe"}
	for _, a := range apps {
		if strings.EqualFold(name, a) {
			return true
		}
	}
	return false
}
