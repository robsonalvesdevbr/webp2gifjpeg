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

1. **Go 1.25 ou superior**

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
   sudo apt install libwebp-dev libgif-dev libjpeg-dev

   # macOS
   brew install webp giflib jpeg pkg-config

   # Fedora/RHEL
   sudo dnf install libwebp-devel giflib-devel libjpeg-turbo-devel
   ```

### Windows

Este projeto usa CGO (Go com código C) e requer bibliotecas nativas do sistema. Você tem três opções principais:

#### Instalação Rápida via go install (PowerShell/CMD)

<details>
<summary><b>Instalar diretamente sem compilar manualmente</b></summary>

Se você já tem o Go instalado no Windows e quer apenas instalar a ferramenta:

**Pré-requisitos:**
1. Go 1.25 ou superior instalado ([download aqui](https://go.dev/dl/))
2. MSYS2 com bibliotecas instaladas (veja Método 1 abaixo)

**Instalação:**

```powershell
# No PowerShell ou CMD
go install github.com/robsonalvesdevbr/webpconvert@latest
```

O executável será instalado em:
- `%USERPROFILE%\go\bin\webpconvert.exe` (normalmente `C:\Users\SeuUsuario\go\bin\`)

**Usar a ferramenta:**

```powershell
# Adicione ao PATH se ainda não estiver
$env:PATH += ";$env:USERPROFILE\go\bin"

# Execute
webpconvert.exe
```

**Importante:** Você ainda precisa ter as DLLs do MSYS2 disponíveis:
- Adicione `C:\msys64\mingw64\bin` ao PATH do Windows, ou
- Copie as DLLs necessárias para `%USERPROFILE%\go\bin\`

**Verificar instalação do Go:**
```powershell
# Verificar se Go está instalado
go version

# Se não estiver instalado, baixe em: https://go.dev/dl/
# Execute o instalador go1.25.x.windows-amd64.msi
# Reinicie o PowerShell e teste novamente
```

</details>

#### Método 1: MSYS2/MinGW (Build manual completo)

<details>
<summary><b>O que é MSYS2 e por que preciso dele?</b></summary>

MSYS2 é um ambiente de desenvolvimento que fornece:
- Ferramentas Unix-like (bash, pacman) para Windows
- MinGW-w64: compilador GCC nativo para Windows
- Gerenciador de pacotes (pacman) para bibliotecas C/C++
- Ambiente isolado que não interfere com o Windows

**Por que usar MSYS2?**
- Este projeto requer bibliotecas específicas (libwebp, giflib, libjpeg)
- MSYS2 fornece versões pré-compiladas dessas bibliotecas
- MinGW-w64 gera executáveis nativos do Windows (.exe)
- É o método mais direto para projetos CGO no Windows

</details>

##### Passo 1: Instalar MSYS2

1. Baixe o instalador do MSYS2:
   - Acesse: https://www.msys2.org/
   - Baixe o instalador: `msys2-x86_64-<data>.exe`
   - Execute o instalador

2. Durante a instalação:
   - Mantenha o caminho padrão: `C:\msys64`
   - Marque "Run MSYS2 now" ao finalizar

3. Atualize o sistema base:
   ```bash
   pacman -Syu
   ```
   - Feche a janela quando solicitado
   - Reabra "MSYS2 MSYS" do menu Iniciar
   - Execute novamente:
   ```bash
   pacman -Su
   ```

##### Passo 2: Instalar Dependências

1. Feche qualquer janela MSYS2 aberta

2. Abra **MSYS2 MINGW64** (não MSYS2 MSYS):
   - Busque no menu Iniciar: "MSYS2 MinGW 64-bit"
   - Ou execute: `C:\msys64\mingw64.exe`

3. Instale as bibliotecas de desenvolvimento:
   ```bash
   pacman -S mingw-w64-x86_64-libwebp \
             mingw-w64-x86_64-giflib \
             mingw-w64-x86_64-libjpeg-turbo \
             mingw-w64-x86_64-pkg-config \
             mingw-w64-x86_64-gcc
   ```

4. Instale o Go:

   **Opção A - Via MSYS2 (mais simples):**
   ```bash
   pacman -S mingw-w64-x86_64-go
   ```

   **Opção B - Go nativo do Windows:**

   Se preferir usar o instalador oficial do Go para Windows:

   1. Baixe o instalador: https://go.dev/dl/
   2. Execute `go1.25.x.windows-amd64.msi`
   3. O Go será instalado em `C:\Program Files\Go`
   4. O instalador adiciona automaticamente ao PATH do Windows
   5. No MSYS2 MINGW64, verifique: `which go`

   Ambas as opções funcionam perfeitamente no MSYS2 MINGW64

##### Passo 3: Verificar Ambiente

```bash
# Verificar compilador C
gcc --version

# Verificar pkg-config
pkg-config --version

# Verificar bibliotecas
pkg-config --libs libwebp
pkg-config --libs libgif
pkg-config --libs libjpeg

# Verificar Go
go version
```

##### Passo 4: Compilar o Projeto

```bash
# Clone o repositório (se ainda não fez)
cd ~
git clone https://github.com/robsonalvesdevbr/webpconvert.git
cd webpconvert

# Compile
CGO_ENABLED=1 go build -o webpconvert.exe
```

##### Passo 5: Usar o Executável

**Dentro do MSYS2 MINGW64:**
```bash
./webpconvert.exe
```

**No Windows (PowerShell/CMD):**

Para usar fora do MSYS2, você precisa das DLLs do MinGW:

Opção A - Adicionar ao PATH do Windows:
1. Adicione `C:\msys64\mingw64\bin` ao PATH do sistema
2. Reinicie o terminal

Opção B - Copiar DLLs para a pasta do executável:
```bash
# No MSYS2 MINGW64:
cp /mingw64/bin/libwebp-7.dll .
cp /mingw64/bin/libgif-7.dll .
cp /mingw64/bin/libjpeg-8.dll .
cp /mingw64/bin/libgcc_s_seh-1.dll .
cp /mingw64/bin/libwinpthread-1.dll .
cp /mingw64/bin/libstdc++-6.dll .
```

**Notas Importantes:**
- Sempre use **MSYS2 MINGW64** para compilar (não MSYS2 MSYS)
- CGO_ENABLED=1 é necessário
- Os executáveis são nativos do Windows (.exe)
- DLLs do MinGW são necessárias em runtime

#### Método 2: WSL (Windows Subsystem for Linux)

<details>
<summary><b>Quando usar WSL em vez de MSYS2?</b></summary>

**Vantagens do WSL:**
- Ambiente Linux completo no Windows
- Mais simples se você já conhece Linux
- Não precisa lidar com DLLs do MinGW
- Melhor integração com ferramentas Linux

**Desvantagens:**
- Executável gerado não é nativo do Windows
- Precisa do WSL para executar
- Requer Windows 10 (build 19041+) ou Windows 11

**Escolha WSL se:**
- Você já usa ou quer aprender Linux
- Não precisa de executável nativo do Windows
- Prefere simplicidade sobre portabilidade

</details>

##### Instalar WSL

1. Abra PowerShell como Administrador

2. Instale WSL com Ubuntu:
   ```powershell
   wsl --install
   ```

3. Reinicie o computador quando solicitado

4. Na primeira execução, crie usuário e senha

##### Instalar Dependências no Ubuntu

```bash
sudo apt update
sudo apt install -y libwebp-dev \
                    libgif-dev \
                    libjpeg-dev \
                    pkg-config \
                    build-essential \
                    git \
                    golang-go
```

##### Compilar e Usar

```bash
# Clonar repositório
git clone https://github.com/robsonalvesdevbr/webpconvert.git
cd webpconvert

# Compilar
go build -o webpconvert

# Usar
./webpconvert
```

**Acessar arquivos do Windows:**
```bash
# Seus arquivos estão em /mnt/c/
./webpconvert /mnt/c/Users/SeuUsuario/Downloads/image.webp
```

#### Solução de Problemas (Windows)

<details>
<summary><b>Erro: "gcc: command not found"</b></summary>

**Solução:**
1. Certifique-se de estar no **MSYS2 MINGW64** (não MSYS2 MSYS)
2. Reinstale o gcc:
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```
3. Reinicie o shell MINGW64

</details>

<details>
<summary><b>Erro: "package libwebp was not found"</b></summary>

**Solução:**
1. Verifique se as bibliotecas estão instaladas:
   ```bash
   pacman -Qs libwebp
   ```
2. Reinstale se necessário:
   ```bash
   pacman -S mingw-w64-x86_64-libwebp
   ```

</details>

<details>
<summary><b>Erro ao executar: "DLL não encontrada"</b></summary>

**Soluções:**

Opção 1 - Adicionar MinGW ao PATH:
- Adicione `C:\msys64\mingw64\bin` às variáveis de ambiente do Windows

Opção 2 - Copiar DLLs:
```bash
# No MSYS2 MINGW64:
ldd webpconvert.exe
# Copie cada DLL listada de /mingw64/bin para a pasta do executável
```

Opção 3 - Executar do MSYS2:
- Use sempre o terminal MINGW64 (DLLs já no PATH)

</details>

<details>
<summary><b>Build muito lento no Windows</b></summary>

**Solução:**
Adicione exceções no antivírus para:
- `C:\msys64`
- `C:\Users\SeuUsuario\go`
- Processos: `gcc.exe`, `go.exe`

</details>

<details>
<summary><b>Diferença entre MSYS2 MSYS e MINGW64?</b></summary>

- **MSYS2 MSYS**: Ambiente Unix-like puro - NÃO use para compilar
- **MSYS2 MINGW64**: Gera binários nativos Windows - USE ESTE

Como identificar: Prompt MINGW64 mostra `MINGW64` no início da linha

</details>

4. **CGO habilitado** (geralmente já está por padrão)

   ```bash
   export CGO_ENABLED=1
   ```

## Instalação da Aplicação

### Via Binários Pré-compilados (Recomendado)

Baixe a versão mais recente para seu sistema operacional na [página de releases](https://github.com/robsonalvesdevbr/webpconvert/releases/latest).

#### Linux

```bash
# Baixe a versão mais recente (substitua VERSION pela versão desejada, ex: v1.4.0)
wget https://github.com/robsonalvesdevbr/webpconvert/releases/download/VERSION/webpconvert-VERSION-linux-amd64.tar.gz

# Extraia o arquivo
tar xzf webpconvert-VERSION-linux-amd64.tar.gz

# Mova para um diretório no PATH
sudo mv webpconvert /usr/local/bin/

# Verifique a instalação
webpconvert --version
```

#### Windows

1. Baixe `webpconvert-VERSION-windows-amd64.zip` da [página de releases](https://github.com/robsonalvesdevbr/webpconvert/releases/latest)
2. Extraia o arquivo ZIP para um diretório de sua escolha (ex: `C:\Program Files\webpconvert\`)
3. Adicione o diretório ao PATH do sistema
4. Abra um novo terminal e execute: `webpconvert --version`

### Via go install

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
sudo apt install libgif-dev

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

## Segurança e Compliance

### SBOM (Software Bill of Materials)

O webpconvert gera e publica **Software Bill of Materials (SBOM)** completo para cada release, proporcionando transparência total sobre dependências e habilitando análise automatizada de segurança.

#### O que é SBOM?

SBOM é um inventário formal de todos os componentes de software, incluindo dependências diretas e transitivas, versões, licenças e hashes criptográficos. Isso permite:

- **Detecção rápida de vulnerabilidades** (CVEs) em dependências
- **Compliance** com frameworks como NIST SSDF e Executive Order 14028
- **Transparência** da cadeia de suprimentos de software
- **Rastreabilidade** de componentes afetados por vulnerabilidades

#### Acessar SBOM

Cada release inclui SBOM em três formatos:

```bash
# Baixar SBOM CycloneDX (recomendado)
gh release download v1.0.0 -p "*sbom.cyclonedx.json"

# Baixar SBOM SPDX (padrão ISO)
gh release download v1.0.0 -p "*sbom.spdx.json"

# Verificar integridade
wget https://github.com/robsonalvesdevbr/webpconvert/releases/download/v1.0.0/sbom_checksums.txt
sha256sum -c sbom_checksums.txt --ignore-missing
```

#### Escanear Vulnerabilidades

```bash
# Usando Grype (recomendado)
grype sbom:webpconvert_1.0.0_sbom.cyclonedx.json

# Usando Trivy
trivy sbom webpconvert_1.0.0_sbom.cyclonedx.json
```

#### Automação de Segurança

O projeto implementa:
- ✅ Geração automática de SBOM em cada release
- ✅ Submissão contínua ao GitHub Dependency Graph
- ✅ Scanning diário de vulnerabilidades com Grype e Trivy
- ✅ Criação automática de issues para vulnerabilidades críticas/altas
- ✅ Integração com GitHub Security (SARIF reports)

#### Documentação Completa

Para informações detalhadas sobre SBOM, formatos disponíveis, análise de licenças e integração com CI/CD, veja:

**[docs/SBOM.md](docs/SBOM.md)** - Guia completo de SBOM

### Reportar Vulnerabilidades

Para reportar vulnerabilidades de segurança, veja [SECURITY.md](SECURITY.md).

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
