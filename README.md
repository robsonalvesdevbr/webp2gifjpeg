# WebP Convert

Conversor de imagens **standalone** em Go com foco em conversões bi-direcionais de/para WebP. Suporta processamento em lote com conversão paralela entre os formatos: WebP, JPEG, PNG, GIF e BMP.

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
- ✅ **Opção para preservar arquivos originais** (flag `--keep-original`)
- ✅ Logging de progresso e erros em tempo real
- ✅ **Implementação nativa em C** (CGO + libwebp + libjpeg + giflib)
- ✅ **Zero dependências runtime** (apenas bibliotecas do sistema)

### Recursos Avançados de Qualidade

- ✅ **Decodificação Avançada WebP**: Configuração otimizada com fancy upsampling e controle de dithering
- ✅ **Quantização de Cores Octree**: Algoritmo customizado para paletas GIF de alta qualidade (256 cores)
- ✅ **Quantização Median Cut**: Algoritmo alternativo para conteúdo fotográfico
- ✅ **Paletas Locais por Frame**: Cada frame GIF tem sua própria paleta otimizada
- ✅ **Distância de Cor Perceptual**: Correspondência de cores ponderada pela sensibilidade humana
- ✅ **JPEG 4:4:4 Chroma**: Sem subsampling de croma para máxima qualidade de cor
- ✅ **Progressive JPEG**: Encoding progressivo com DCT de alta qualidade (JDCT_ISLOW)
- ✅ **Composição Alpha com Precisão**: Aritmética de ponto flutuante para evitar perda de qualidade

## Requisitos

### Para Uso (Runtime)

**Nenhuma dependência adicional!** O binário é standalone e usa apenas bibliotecas do sistema que já estão instaladas:

- `libwebp7` / `libwebpdemux2` (geralmente já instalado)
- `libgif7` (geralmente já instalado)
- `libjpeg` / `libjpeg-turbo` (geralmente já instalado)

### Para Desenvolvimento (Build)

1. **Go 1.23 ou superior**

   ```bash
   go version
   ```

2. **Compilador C e pkg-config**

   ```bash
   # Ubuntu/Debian
   sudo apt install build-essential pkg-config

   # macOS (via Xcode Command Line Tools)
   xcode-select --install

   # Fedora/RHEL
   sudo dnf install gcc pkg-config
   ```

3. **Bibliotecas de desenvolvimento**

   ```bash
   # Ubuntu/Debian
   sudo apt install libwebp-dev giflib-dev libjpeg-dev

   # macOS
   brew install webp giflib jpeg pkg-config

   # Fedora/RHEL
   sudo dnf install libwebp-devel giflib-devel libjpeg-turbo-devel

   # Windows (MSYS2)
   pacman -S mingw-w64-x86_64-libwebp \
             mingw-w64-x86_64-giflib \
             mingw-w64-x86_64-libjpeg-turbo \
             mingw-w64-x86_64-pkg-config \
             mingw-w64-x86_64-gcc
   ```

4. **CGO habilitado** (geralmente já está por padrão)

   ```bash
   export CGO_ENABLED=1
   ```

## Instalação da Aplicação

### Via go install (Recomendado)

```bash
CGO_ENABLED=1 go install github.com/robsonalvesdevbr/webpconvert@latest
```

A ferramenta estará disponível em `~/go/bin/webpconvert` (ou `$GOPATH/bin/webpconvert`).

**Instalando versões específicas:**

```bash
# Versão mais recente (v1.3.0 - implementação nativa CGO)
CGO_ENABLED=1 go install github.com/robsonalvesdevbr/webpconvert@v1.3.0

# Versão legada (v1.0.0 - implementação Python, não recomendado)
go install github.com/robsonalvesdevbr/webp2gifjpeg@v1.0.0
```

> **Nota sobre Versionamento:**
> - **v1.3.0+**: Implementação nativa com CGO (recomendado) - módulo `webpconvert`
> - **v1.2.0**: Transição (renomeação do projeto)
> - **v1.0.0**: Versão legada com dependência Python - módulo `webp2gifjpeg`
>
> ⚠️ **Breaking Changes em v1.3.0**: Migração completa de Python para CGO nativo, mudança de nome do módulo e nova API.

### Via Clone do Repositório

```bash
# Clone o repositório
git clone https://github.com/robsonalvesdevbr/webpconvert.git
cd webpconvert

# Compile a aplicação
CGO_ENABLED=1 go build -o webpconvert
```

**Nota**: CGO_ENABLED=1 é necessário para compilar o código C nativo.

## Uso

### Processando o diretório atual

```bash
./webpconvert
```

### Processando um diretório específico

```bash
./webpconvert -dir /caminho/para/diretorio
```

### Configurando qualidade JPEG

```bash
./webpconvert -quality 95
```

### Processamento paralelo

```bash
# Usar 4 workers paralelos
./webpconvert -workers 4

# Usar todos os núcleos da CPU (padrão)
./webpconvert

# Processamento sequencial (1 worker)
./webpconvert -workers 1
```

### Preservar arquivos originais

```bash
# Manter arquivos WebP originais após conversão
./webpconvert --keep-original

# Os arquivos convertidos terão sufixo "_converted"
# Exemplo: image.webp → image_converted.jpg + image.webp (preservado)
```

### Exemplos

```bash
# Converter todos os WebP no diretório atual (qualidade JPEG padrão: 100, workers: CPU count)
./webpconvert

# Converter todos os WebP em um diretório específico
./webpconvert -dir ./imagens

# Alta qualidade JPEG para fotos profissionais
./webpconvert -dir ./fotos -quality 100

# Qualidade menor para web (arquivos menores)
./webpconvert -dir ./web-images -quality 75

# Processamento rápido com 8 workers paralelos
./webpconvert -dir ./fotos -workers 8

# Processamento sequencial para economia de recursos
./webpconvert -dir ./imagens -workers 1

# Preservar arquivos WebP originais
./webpconvert --keep-original

# Preservar originais com qualidade customizada
./webpconvert --keep-original -quality 90

# Combinar todas as opções
./webpconvert -dir /home/usuario/fotos -quality 95 -workers 4 --keep-original
```

## Estrutura do Projeto

```
webpconvert/
├── main.go                    # Aplicação principal (CLI)
├── webpconvert                # Binário compilado
├── converter/
│   ├── converter.go           # Lógica de conversão e processamento
│   └── converter_test.go      # Testes unitários
├── native/                    # Implementação nativa em C via CGO
│   ├── webp_detector.go       # Detecção de tipo WebP (animado/estático)
│   ├── webp_decoder.go        # Decodificador WebP avançado com RGBA/BGRA
│   ├── webp_to_jpeg.go        # Conversão WebP → JPEG com 4:4:4 chroma
│   ├── webp_to_gif.go         # Conversão WebP → GIF com paletas locais
│   ├── octree_quantizer.go    # Algoritmo Octree para quantização de cores
│   └── median_cut.go          # Algoritmo Median Cut para conteúdo fotográfico
├── go.mod                     # Dependências
├── .gitignore                 # Arquivos ignorados pelo Git
└── README.md                  # Documentação
```

## Como Funciona

### Arquitetura Nativa

1. **Detecção de Tipo**: Usa `libwebp` (WebPDemux) para detectar se o WebP é animado ou estático

2. **Conversão WebP → JPEG** (Alta Qualidade):
   - Decode WebP usando decodificador avançado com configuração otimizada
   - Fancy upsampling para melhor qualidade
   - Composição alpha com aritmética de ponto flutuante (precisão)
   - Encode JPEG com:
     - Chroma subsampling 4:4:4 (sem perda de cor)
     - DCT método JDCT_ISLOW (máxima qualidade)
     - Progressive encoding
     - Huffman optimization
     - Qualidade configurável (default: 100)

3. **Conversão WebP Animado → GIF** (Paletas Otimizadas):
   - Demux WebP usando `WebPDemuxer` (libwebpdemux)
   - Decode de cada frame com `WebPDecodeRGBA` (composição automática)
   - Quantização de cores por frame usando **Octree** ou **Median Cut**:
     - Paleta local de 256 cores otimizada para cada frame
     - Distância de cor perceptual (ponderada: verde > vermelho > azul)
     - Cache de correspondência de cores para performance
   - Encode GIF usando `giflib`:
     - Suporte a looping infinito (Netscape 2.0 extension)
     - Preservação de timing entre frames
     - Disposal method configurável

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

## Troubleshooting

### Problemas Comuns de Build

#### "pkg-config not found"
Instale o pkg-config:
```bash
# Ubuntu/Debian
sudo apt install pkg-config

# macOS
brew install pkg-config

# Fedora/RHEL
sudo dnf install pkg-config
```

#### "Package libwebp was not found"
Instale as bibliotecas de desenvolvimento:
```bash
# Ubuntu/Debian
sudo apt install libwebp-dev

# macOS
brew install webp

# Fedora/RHEL
sudo dnf install libwebp-devel
```

Verifique a instalação:
```bash
pkg-config --modversion libwebp
pkg-config --cflags --libs libwebp
```

#### "undefined reference to `WebP*`"
Certifique-se que CGO está habilitado:
```bash
export CGO_ENABLED=1
go build -v
```

#### "cannot find -lgif"
Instale giflib:
```bash
# Ubuntu/Debian
sudo apt install giflib-dev

# macOS
brew install giflib

# Fedora/RHEL
sudo dnf install giflib-devel
```

### Problemas de Runtime

#### "error while loading shared libraries: libwebp.so.7"
Instale as bibliotecas runtime:
```bash
# Ubuntu/Debian
sudo apt install libwebp7 libwebpdemux2

# Fedora/RHEL
sudo dnf install libwebp

# macOS
brew install webp
```

#### Verificar bibliotecas disponíveis
```bash
# Linux
ldd ./webpconvert

# macOS
otool -L ./webpconvert
```

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

- **Backup**: Por padrão, a aplicação substitui os arquivos originais. Use `--keep-original` para preservá-los ou faça backup antes de executar.
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

### Implementado
- [x] ~~Processamento paralelo de múltiplos arquivos~~ ✅
- [x] ~~Versão standalone sem dependência de Python~~ ✅
- [x] ~~Decodificação avançada WebP~~ ✅
- [x] ~~Quantização de cores Octree/Median Cut~~ ✅
- [x] ~~Paletas locais por frame GIF~~ ✅
- [x] ~~JPEG 4:4:4 chroma subsampling~~ ✅
- [x] ~~Opção para preservar arquivos originais (flag `--keep-original`)~~ ✅

### Planejado
- [ ] Configuração de qualidade/compressão do GIF
- [ ] Progress bar para conversões longas
- [ ] Suporte a outras conversões (GIF→WebP, PNG→WebP, etc)
- [ ] Static linking opcional para binário completamente portável

### Otimizações de Performance (Identificadas)
- [ ] Corrigir memory leak do `defer` em loop (frames GIF)
- [ ] Buffer pooling com `sync.Pool` para reutilização de memória
- [ ] Reduzir alocações desnecessárias em RGBA→RGB
- [ ] Otimizar tamanho de channels para batch processing
- [ ] Remover código morto (`analyzeAllFramesForGlobalPalette`)
- [ ] Implementar paleta incremental para frames similares

## Build Avançado

### Static Linking (opcional)

Para criar um binário completamente estático (sem dependências de bibliotecas dinâmicas):

```bash
# Ubuntu/Debian - instalar versões estáticas
sudo apt install libwebp-dev:native libgif-dev:native libjpeg-dev:native

# Build estático
CGO_ENABLED=1 go build -ldflags="-linkmode external -extldflags '-static'" -o webpconvert
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
