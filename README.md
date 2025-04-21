 
# GuardianEye 🛡️👁️

**Proteção digital para famílias e ambientes de trabalho**

![GuardianEye Logo](https://via.placeholder.com/200x200.png?text=GuardianEye)

## Sobre o GuardianEye

GuardianEye é uma ferramenta de controle e filtragem de conteúdo desenvolvida em Go, projetada para proteger ambientes digitais contra acesso a conteúdos indesejados. Ela monitora a atividade do sistema e bloqueia automaticamente o acesso a conteúdos impróprios baseados em palavras-chave configuráveis e domínios proibidos.

### Principais Recursos

- **Monitoramento de Processos:** Analisa janelas e processos em tempo real
- **Verificação de Domínios:** Bloqueia acesso a sites proibidos (requer Npcap)
- **Lista de Palavras-Chave Configurável:** Personalização completa do filtro
- **Interface Intuitiva:** Ícone na bandeja do sistema e CLI para gerenciamento
- **Operação em Segundo Plano:** Funciona discretamente sem interromper o fluxo de trabalho

## Instalação

### Pré-requisitos

- Sistema operacional Windows
- [Go 1.16+](https://golang.org/dl/) (apenas para compilação)
- [Npcap](https://npcap.com) (para funcionalidade de monitoramento de rede)

### Instalação Rápida (Binário Pré-compilado)

1. Baixe o arquivo ZIP mais recente da seção de [Releases](https://github.com/yourusername/guardianeye/releases)
2. Extraia o conteúdo para uma pasta de sua escolha
3. Execute `guardianeye.exe` como administrador

### Compilação a partir do código-fonte

```bash
# Clone o repositório
git clone https://github.com/yourusername/guardianeye.git
cd guardianeye

# Baixe as dependências
go mod tidy

# Compile o projeto
go build -o guardianeye.exe cmd/main.go
```

## Configuração

O arquivo `config.yaml` deve estar no mesmo diretório do executável e contém as listas de palavras-chave e domínios a serem bloqueados:

```yaml
# Palavras-chave proibidas
keywords:
  - adult
  - gambling
  - piracy
  - porn
  - spyware
  - virus
  - malware
  - adware

# Domínios bloqueados
domains:
  - pornhub.com
  - xvideos.com
  - adultfriendfinder.com
  - xhamster.com
```

### Personalização

Você pode personalizar as listas de palavras-chave e domínios de acordo com suas necessidades:

1. Abra o arquivo `config.yaml` em qualquer editor de texto
2. Adicione ou remova entradas nas seções `keywords` e `domains`
3. Salve o arquivo e reinicie o GuardianEye (ou use a opção "Recarregar Config" na interface)

## Uso

1. Execute o programa `guardianeye.exe` (preferencialmente como administrador)
2. Um ícone aparecerá na bandeja do sistema (área de notificações)
3. Clique com o botão direito no ícone para acessar as opções:
   - **Abrir CLI:** Abre a interface de linha de comando
   - **Sair:** Encerra o programa

### Interface de Linha de Comando (CLI)

A interface CLI permite:

- Carregar/editar lista de bloqueios
- Iniciar monitoramento
- Parar monitoramento 
- Exibir status
- Sair da aplicação

## Logs

Os logs da aplicação são gravados no arquivo `accessctl.log` no diretório da aplicação e contêm informações sobre:

- Inicialização e encerramento do programa
- Detecções de conteúdo bloqueado
- Processos encerrados
- Erros e avisos

## Arquitetura do Projeto

O GuardianEye é organizado nos seguintes pacotes:

- **cmd:** Contém o ponto de entrada do programa (`main.go`)
- **config:** Gerenciamento e carregamento de configurações
- **controller:** Ações de controle como término de processos
- **logger:** Sistema de logging
- **monitor:** Monitoramento de processos, rede e teclado
- **ui:** Interface com o usuário (CLI e ícone da bandeja)

## Limitações Conhecidas

- O monitoramento de rede requer a instalação do Npcap
- Compatível apenas com Windows
- Não possui criptografia para logs e configurações
- Não inclui autenticação para impedir a desativação

## Roadmap (Próximas Versões)

- [ ] Suporte para autenticação com senha
- [ ] Criptografia de logs e configurações
- [ ] Interface web administrativa
- [ ] Suporte para filtros baseados em IA para melhor detecção
- [ ] Suporte para outros sistemas operacionais
- [ ] Atualização automática de listas de bloqueio
- [ ] Análise de imagens para conteúdo inadequado

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para enviar pull requests ou abrir issues no GitHub.

1. Faça um fork do projeto
2. Crie sua branch de recurso (`git checkout -b feature/AmazingFeature`)
3. Commit suas alterações (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

Link do Projeto: [https://github.com/esc4n0rx/guardianeye](https://github.com/esc4n0rx/guardianeye)

---

Desenvolvido com ❤️ para proteger o que importa.