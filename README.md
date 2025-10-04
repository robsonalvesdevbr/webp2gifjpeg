# WebP to GIF/JPEG Converter

Aplicação em Go para converter arquivos WebP automaticamente para o formato apropriado: WebP animados → GIF e WebP estáticos → JPEG, processando recursivamente todos os arquivos em um diretório.

## Funcionalidades

- ✅ Detecção automática de WebP animado vs estático
- ✅ Conversão de WebP animado para GIF
- ✅ Conversão de WebP estático para JPEG
- ✅ Qualidade JPEG configurável (1-100)
- ✅ Preservação de metadados EXIF (JPEG)
- ✅ Tratamento de transparência (fundo branco em JPEG)
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

**Importante:** Certifique-se de que os arquivos Python (`detect_webp_type.py`, `webp_to_gif.py`, `webp_to_jpeg.py`) estão no mesmo diretório do binário `webp2gifjpeg`.

## Uso

### Processando o diretório atual

```bash
./webp2gifjpeg
```

### Processando um diretório específico

```bash
./webp2gifjpeg -dir /caminho/para/diretorio
```

### Configurando qualidade JPEG

```bash
./webp2gifjpeg -quality 95
```

### Exemplos

```bash
# Converter todos os WebP no diretório atual (qualidade JPEG padrão: 85)
./webp2gifjpeg

# Converter todos os WebP em um diretório específico
./webp2gifjpeg -dir ./imagens

# Alta qualidade JPEG para fotos profissionais
./webp2gifjpeg -dir ./fotos -quality 95

# Qualidade menor para web (arquivos menores)
./webp2gifjpeg -dir ./web-images -quality 75

# Converter todos os WebP incluindo subdiretórios
./webp2gifjpeg -dir /home/usuario/fotos
```

## Estrutura do Projeto

```
webp2gifjpeg/
├── main.go                    # Aplicação principal (CLI)
├── webp2gifjpeg               # Binário compilado
├── detect_webp_type.py        # Script Python para detecção de tipo
├── webp_to_gif.py             # Script Python para conversão GIF
├── webp_to_jpeg.py            # Script Python para conversão JPEG
├── converter/
│   ├── converter.go           # Lógica de conversão e detecção
│   └── converter_test.go      # Testes unitários
├── go.mod                     # Dependências (vazio - sem deps externas)
└── README.md                  # Documentação
```

## Como Funciona

1. A aplicação Go percorre recursivamente o diretório especificado
2. Identifica todos os arquivos com extensão `.webp`
3. Para cada arquivo WebP:
   - Detecta se é animado ou estático usando `detect_webp_type.py`
   - **Se animado**: converte para GIF usando `webp_to_gif.py`
   - **Se estático**: converte para JPEG usando `webp_to_jpeg.py`
4. Os scripts Python usam Pillow para conversão com alta qualidade
5. Substitui o arquivo original `.webp` pelo novo `.gif` ou `.jpg`
6. Exibe um resumo detalhado com estatísticas de conversão

### Arquitetura Híbrida

- **Go**: Gerenciamento de arquivos, busca recursiva, orquestração, CLI
- **Python/Pillow**: Detecção e conversão de imagens
  - Suporte completo para WebP animado (múltiplos frames, delays)
  - Conversão de WebP estático com preservação de EXIF
  - Tratamento de transparência (composite em fundo branco para JPEG)

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
- ✅ Conversão de WebP para JPEG
- ✅ Processamento de diretórios recursivo com detecção automática
- ✅ Detecção de tipo de WebP (animado vs estático)
- ✅ Tratamento de erros (arquivo inexistente, diretório inválido)
- ✅ Verificação de substituição de arquivos
- ✅ Validação de qualidade JPEG

## Dependências

### Runtime

- **Python 3** com **Pillow** - Conversão de imagens

### Desenvolvimento (Go)

- Nenhuma dependência externa Go necessária

## Observações

- **Backup**: A aplicação substitui os arquivos originais. Faça backup antes de executar.
- **WebP Animado**: Suporte completo via Pillow - todos os frames e delays são preservados no GIF.
- **WebP Estático**: Convertido para JPEG com qualidade configurável (padrão: 85).
- **Transparência**: WebP com canal alpha são convertidos para JPEG com fundo branco.
- **EXIF**: Metadados EXIF são preservados na conversão para JPEG.
- **Performance**: O processamento é feito sequencialmente. Para grandes volumes, considere adicionar processamento paralelo.
- **Distribuição**: Para distribuir o binário, inclua todos os arquivos Python (`detect_webp_type.py`, `webp_to_gif.py`, `webp_to_jpeg.py`) no mesmo diretório do executável.

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
