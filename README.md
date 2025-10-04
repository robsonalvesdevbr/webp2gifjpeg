# WebP to GIF Converter

Aplicação em Go para converter arquivos WebP animados para GIF, processando recursivamente todos os arquivos em um diretório.

## Funcionalidades

- ✅ Conversão de WebP animado para GIF
- ✅ Processamento recursivo de diretórios
- ✅ Substituição automática dos arquivos WebP originais
- ✅ Logging de progresso e erros
- ✅ Testes unitários completos

## Requisitos

### Obrigatórios

1. **Go 1.21 ou superior**

   ```bash
   go version
   ```

2. **Python 3**

   ```bash
   python3 --version
   ```

3. **Pillow (biblioteca Python para processamento de imagens)**
   ```bash
   pip3 install Pillow
   # ou em sistemas que requerem flag
   pip3 install --break-system-packages Pillow
   ```

### Verificação Rápida

Execute para verificar se todos os requisitos estão instalados:

```bash
go version && python3 --version && python3 -c "import PIL; print('Pillow:', PIL.__version__)"
```

### Instalação dos Requisitos por Sistema Operacional

#### Ubuntu/Debian

```bash
# Instalar Go
sudo apt install golang-go

# Python3 geralmente já vem instalado
sudo apt install python3 python3-pip

# Instalar Pillow
pip3 install --break-system-packages Pillow
```

#### macOS

```bash
# Instalar Go
brew install go

# Python3 geralmente já vem instalado
brew install python3

# Instalar Pillow
pip3 install Pillow
```

#### Windows

- Baixar e instalar Go: https://go.dev/dl/
- Baixar e instalar Python: https://www.python.org/downloads/
- Instalar Pillow via CMD/PowerShell:
  ```cmd
  pip install Pillow
  ```

## Instalação da Aplicação

```bash
# Clone o repositório
git clone https://github.com/robson/webp2gifjpeg.git
cd webp2gifjpeg

# Compile a aplicação
go build -o webp2gifjpeg
```

**Importante:** Certifique-se de que o arquivo `webp_to_gif.py` está no mesmo diretório do binário `webp2gifjpeg`.

## Uso

### Processando o diretório atual

```bash
./webp2gifjpeg
```

### Processando um diretório específico

```bash
./webp2gifjpeg -dir /caminho/para/diretorio
```

### Exemplos

```bash
# Converter todos os WebP no diretório atual
./webp2gifjpeg

# Converter todos os WebP em um diretório específico
./webp2gifjpeg -dir ./imagens

# Converter todos os WebP incluindo subdiretórios
./webp2gifjpeg -dir /home/usuario/fotos
```

## Estrutura do Projeto

```
webp2gifjpeg/
├── main.go                    # Aplicação principal
├── webp2gifjpeg                   # Binário compilado
├── webp_to_gif.py            # Script Python para conversão
├── converter/
│   ├── converter.go          # Lógica de conversão
│   └── converter_test.go     # Testes unitários
├── go.mod                    # Dependências (vazio - sem deps externas)
└── README.md                 # Documentação
```

## Como Funciona

1. A aplicação Go percorre recursivamente o diretório especificado
2. Identifica todos os arquivos com extensão `.webp`
3. Para cada arquivo, chama o script Python `webp_to_gif.py`
4. O script Python usa Pillow para converter WebP → GIF (com suporte completo a animações)
5. Substitui o arquivo original `.webp` pelo novo `.gif`
6. Exibe um resumo com quantidade de arquivos convertidos e erros

### Arquitetura Híbrida

- **Go**: Gerenciamento de arquivos, busca recursiva, orquestração
- **Python/Pillow**: Conversão real (suporte completo para WebP animado com múltiplos frames)

## Testes

**Nota:** Os testes requerem `ffmpeg` instalado para criar arquivos WebP de teste.

Execute os testes com:

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes com verbose
go test -v ./...
```

## Testes Incluídos

- ✅ Conversão básica de WebP para GIF
- ✅ Processamento de diretórios recursivo
- ✅ Detecção de WebP animado
- ✅ Tratamento de erros (arquivo inexistente, diretório inválido)
- ✅ Verificação de substituição de arquivos

## Dependências

### Runtime

- **Python 3** com **Pillow** - Conversão de imagens

### Desenvolvimento (Go)

- Nenhuma dependência externa Go necessária

## Observações

- **Backup**: A aplicação substitui os arquivos originais. Faça backup antes de executar.
- **WebP Animado**: Suporte completo via Pillow - todos os frames e delays são preservados.
- **Performance**: O processamento é feito sequencialmente. Para grandes volumes, considere adicionar processamento paralelo.
- **Distribuição**: Para distribuir o binário, inclua tanto `webp2gifjpeg` quanto `webp_to_gif.py` no mesmo diretório.

## Melhorias Futuras

- [ ] Processamento paralelo de múltiplos arquivos
- [ ] Opção para preservar arquivos originais (flag `--keep-original`)
- [ ] Configuração de qualidade/compressão do GIF
- [ ] Progress bar para conversões longas
- [ ] Suporte a outras conversões (GIF→WebP, PNG→WebP, etc)
- [ ] Versão standalone sem dependência de Python (usando CGO + libwebp)

## Licença

MIT

## Contribuindo

Contribuições são bem-vindas! Por favor:

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Suporte

Para reportar bugs ou solicitar features, abra uma issue no GitHub.
