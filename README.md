# WebP to GIF/JPEG Converter

Aplicação **standalone** em Go para converter arquivos WebP automaticamente para o formato apropriado: WebP animados → GIF e WebP estáticos → JPEG, processando recursivamente todos os arquivos em um diretório.

**✨ Versão Native**: Implementação 100% nativa usando CGO + libwebp, sem dependências de Python em runtime!

## Funcionalidades

- ✅ Detecção automática de WebP animado vs estático
- ✅ Conversão de WebP animado para GIF
- ✅ Conversão de WebP estático para JPEG
- ✅ Qualidade JPEG configurável (1-100, default: 100)
- ✅ **Processamento paralelo** com workers configuráveis
- ✅ Tratamento de transparência (fundo branco em JPEG)
- ✅ Processamento recursivo de diretórios
- ✅ Substituição automática dos arquivos WebP originais
- ✅ Logging de progresso e erros em tempo real
- ✅ **Implementação nativa em C** (CGO + libwebp + libjpeg + giflib)
- ✅ **Zero dependências runtime** (apenas bibliotecas do sistema)

## Requisitos

### Para Uso (Runtime)

**Nenhuma dependência adicional!** O binário é standalone e usa apenas bibliotecas do sistema que já estão instaladas:

- `libwebp7` (geralmente já instalado)
- `libgif7` (geralmente já instalado)
- `libjpeg` (geralmente já instalado)

### Para Desenvolvimento (Build)

1. **Go 1.21 ou superior**

   ```bash
   go version
   ```

2. **Bibliotecas de desenvolvimento**

   ```bash
   # Ubuntu/Debian
   sudo apt install libwebp-dev libgif-dev libjpeg-dev

   # macOS
   brew install webp giflib jpeg

   # Fedora/RHEL
   sudo dnf install libwebp-devel giflib-devel libjpeg-turbo-devel
   ```

3. **CGO habilitado** (geralmente já está por padrão)

   ```bash
   export CGO_ENABLED=1
   ```

## Instalação da Aplicação

### Via go install (Recomendado)

```bash
CGO_ENABLED=1 go install github.com/robsonalvesdevbr/webp2gifjpeg@latest
```

A ferramenta estará disponível em `~/go/bin/webp2gifjpeg` (ou `$GOPATH/bin/webp2gifjpeg`).

### Via Clone do Repositório

```bash
# Clone o repositório
git clone https://github.com/robsonalvesdevbr/webp2gifjpeg.git
cd webp2gifjpeg

# Compile a aplicação
CGO_ENABLED=1 go build -o webp2gifjpeg
```

**Nota**: CGO_ENABLED=1 é necessário para compilar o código C nativo.

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

### Processamento paralelo

```bash
# Usar 4 workers paralelos
./webp2gifjpeg -workers 4

# Usar todos os núcleos da CPU (padrão)
./webp2gifjpeg

# Processamento sequencial (1 worker)
./webp2gifjpeg -workers 1
```

### Exemplos

```bash
# Converter todos os WebP no diretório atual (qualidade JPEG padrão: 100, workers: CPU count)
./webp2gifjpeg

# Converter todos os WebP em um diretório específico
./webp2gifjpeg -dir ./imagens

# Alta qualidade JPEG para fotos profissionais
./webp2gifjpeg -dir ./fotos -quality 100

# Qualidade menor para web (arquivos menores)
./webp2gifjpeg -dir ./web-images -quality 75

# Processamento rápido com 8 workers paralelos
./webp2gifjpeg -dir ./fotos -workers 8

# Processamento sequencial para economia de recursos
./webp2gifjpeg -dir ./imagens -workers 1

# Combinar qualidade e workers
./webp2gifjpeg -dir /home/usuario/fotos -quality 95 -workers 4
```

## Estrutura do Projeto

```
webp2gifjpeg/
├── main.go                    # Aplicação principal (CLI)
├── webp2gifjpeg               # Binário compilado
├── converter/
│   ├── converter.go           # Lógica de conversão e processamento
│   └── converter_test.go      # Testes unitários
├── native/                    # Implementação nativa em C via CGO
│   ├── webp_detector.go       # Detecção de tipo WebP (animado/estático)
│   ├── webp_to_jpeg.go        # Conversão WebP → JPEG usando libjpeg
│   └── webp_to_gif.go         # Conversão WebP animado → GIF usando giflib
├── go.mod                     # Dependências
├── .gitignore                 # Arquivos ignorados pelo Git
└── README.md                  # Documentação
```

## Como Funciona

### Arquitetura Nativa

1. **Detecção de Tipo**: Usa `libwebp` (WebPDemux) para detectar se o WebP é animado ou estático
2. **Conversão WebP → JPEG**:
   - Decode WebP usando `WebPDecodeRGBA` (libwebp)
   - Tratamento de transparência (composite em fundo branco)
   - Encode JPEG usando `libjpeg` com qualidade configurável
3. **Conversão WebP Animado → GIF**:
   - Demux WebP animado usando `WebPAnimDecoder` (libwebpdemux)
   - Extração de todos os frames + delays
   - Quantização de cores (256 cores)
   - Encode GIF usando `giflib` com suporte a looping
4. **Processamento**:
   - Scan recursivo do diretório para encontrar arquivos `.webp`
   - Processamento paralelo usando goroutines (workers configuráveis)
   - Substituição automática dos arquivos originais
5. **Estatísticas**: Exibe resumo detalhado com contadores de conversão

### Performance

O processamento paralelo oferece ganhos significativos de performance:
- **1 worker** (sequencial): Baseline
- **4 workers**: ~3-4x mais rápido em CPUs quad-core
- **8 workers**: ~6-7x mais rápido em CPUs com 8+ cores
- **N workers**: Escalável até o número de núcleos disponíveis

**Performance Nativa vs Python**:
- **3-5x mais rápido** que a versão baseada em Python/Pillow
- Sem overhead de spawning processos externos
- Acesso direto às bibliotecas nativas via CGO

**Recomendações**:
- Para poucos arquivos (< 10): Use `-workers 1` (overhead mínimo)
- Para muitos arquivos: Use default (CPU count) ou ajuste conforme necessário
- Para sistemas com poucos recursos: Limite workers para evitar sobrecarga

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

# Executar com race detector
go test -race ./...
```

## Testes Incluídos

- ✅ Detecção de tipo WebP (animado vs estático)
- ✅ Conversão de WebP estático para JPEG com qualidade configurável
- ✅ Conversão de WebP animado para GIF
- ✅ Processamento de diretórios recursivo
- ✅ Processamento paralelo com múltiplos workers
- ✅ Tratamento de erros (arquivo inexistente, diretório inválido)
- ✅ Verificação de substituição de arquivos
- ✅ Validação de qualidade JPEG

## Dependências

### Runtime

**Nenhuma dependência Go!** Apenas bibliotecas do sistema:
- `libwebp7`, `libwebpdemux2`, `libwebpmux3` - Leitura e decode de WebP
- `libjpeg` / `libjpeg-turbo` - Encode JPEG
- `libgif7` (giflib) - Encode GIF

### Build (Desenvolvimento)

- `libwebp-dev` - Headers para desenvolvimento libwebp
- `libgif-dev` - Headers para desenvolvimento giflib
- `libjpeg-dev` - Headers para desenvolvimento libjpeg
- CGO habilitado (padrão no Go)

## Observações

- **Backup**: A aplicação substitui os arquivos originais. Faça backup antes de executar.
- **WebP Animado**: Suporte completo via libwebp - todos os frames e delays são preservados no GIF.
- **WebP Estático**: Convertido para JPEG com qualidade configurável (padrão: 100).
- **Transparência**: WebP com canal alpha são convertidos para JPEG com fundo branco.
- **Performance**: Implementação nativa em C oferece performance 3-5x superior à versão Python.
- **Thread Safety**: Código validado com `go test -race` - sem race conditions.
- **Cross-Platform**: Funciona em Linux, macOS e Windows (com mingw-w64).

## Vantagens da Versão Native

✅ **Zero dependências runtime** (exceto libs do sistema)
✅ **Binário standalone verdadeiro**
✅ **Performance 3-5x melhor** (sem overhead de processos Python)
✅ **Código 100% Go + C nativo**
✅ **Menor consumo de memória**
✅ **Startup instantâneo** (sem inicialização Python)

## Melhorias Futuras

- [x] ~~Processamento paralelo de múltiplos arquivos~~ ✅ Implementado!
- [x] ~~Versão standalone sem dependência de Python~~ ✅ Implementado!
- [ ] Opção para preservar arquivos originais (flag `--keep-original`)
- [ ] Configuração de qualidade/compressão do GIF
- [ ] Progress bar para conversões longas
- [ ] Suporte a outras conversões (GIF→WebP, PNG→WebP, etc)
- [ ] Static linking opcional para binário completamente portável

## Build Avançado

### Static Linking (opcional)

Para criar um binário completamente estático (sem dependências de bibliotecas dinâmicas):

```bash
# Ubuntu/Debian - instalar versões estáticas
sudo apt install libwebp-dev:native libgif-dev:native libjpeg-dev:native

# Build estático
CGO_ENABLED=1 go build -ldflags="-linkmode external -extldflags '-static'" -o webp2gifjpeg
```

**Nota**: Static linking pode não funcionar em todos os sistemas operacionais.

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
