# Software Bill of Materials (SBOM)

## Visão Geral

O webpconvert gera e publica Software Bill of Materials (SBOM) completo para cada release, proporcionando transparência total sobre componentes e dependências, além de permitir análise de segurança automatizada.

## O que é SBOM?

SBOM (Software Bill of Materials) é um inventário formal e estruturado de todos os componentes de software, incluindo:

- **Dependências diretas e transitivas**: Todas as bibliotecas utilizadas pelo projeto
- **Versões específicas**: Informações exatas de versão de cada componente
- **Informações de licenciamento**: Licenças de cada dependência
- **Hashes criptográficos**: Checksums para verificação de integridade
- **Metadados de fornecedor**: Informações sobre origem dos componentes

### Por que SBOM é importante?

1. **Segurança**: Identificação rápida de vulnerabilidades (CVEs) em dependências
2. **Compliance**: Atendimento a frameworks como NIST SSDF, Executive Order 14028
3. **Transparência**: Visibilidade completa da cadeia de suprimentos de software
4. **Rastreabilidade**: Capacidade de rastrear componentes afetados por vulnerabilidades
5. **Governança**: Facilita auditorias e gestão de licenças

## Formatos Disponíveis

O webpconvert disponibiliza SBOM em três formatos:

### 1. CycloneDX JSON (Recomendado)

**Arquivo**: `webpconvert_<versão>_sbom.cyclonedx.json`

- Formato moderno e rico em dados
- Amplamente suportado por ferramentas de análise
- Inclui informações de vulnerabilidades e licenças
- **Use este formato** para a maioria dos casos

### 2. SPDX JSON (Padrão ISO)

**Arquivo**: `webpconvert_<versão>_sbom.spdx.json`

- Padrão ISO/IEC 5962:2021
- Requisito em contratos governamentais
- Formato mais antigo e estável
- Use para requisitos de compliance

### 3. Syft JSON (Nativo)

**Arquivo**: `webpconvert_<versão>_sbom.syft.json`

- Formato nativo da ferramenta Syft (Anchore)
- Rico em metadados técnicos
- Otimizado para análise com Grype
- Use para troubleshooting avançado

## Como Acessar SBOM

### Via GitHub Releases

Cada release do webpconvert inclui arquivos SBOM como assets:

```bash
# Listar SBOMs disponíveis
gh release view v1.0.0 --json assets -q '.assets[].name | select(contains("sbom"))'

# Baixar SBOM CycloneDX
gh release download v1.0.0 -p "*sbom.cyclonedx.json"

# Baixar SBOM SPDX
gh release download v1.0.0 -p "*sbom.spdx.json"

# Baixar todos os SBOMs
gh release download v1.0.0 -p "*sbom*.json"
```

### Via Navegador

1. Acesse a página de [Releases](https://github.com/robsonalvesdevbr/webpconvert/releases)
2. Selecione a versão desejada
3. Baixe os arquivos `*_sbom.*.json` na seção de Assets

### Via GitHub Dependency Graph

O projeto também submete SBOM automaticamente para o GitHub Dependency Graph:

1. Acesse o repositório
2. Vá em **Insights** → **Dependency graph**
3. Visualize todas as dependências e suas vulnerabilidades conhecidas

## Análise de Vulnerabilidades

### Usando Grype (Recomendado)

[Grype](https://github.com/anchore/grype) é um scanner de vulnerabilidades de código aberto:

```bash
# Instalar Grype
curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh | sh -s -- -b /usr/local/bin

# Escanear SBOM
grype sbom:webpconvert_1.0.0_sbom.cyclonedx.json

# Gerar relatório JSON
grype sbom:webpconvert_1.0.0_sbom.cyclonedx.json -o json > vulnerabilities.json

# Filtrar apenas vulnerabilidades críticas/altas
grype sbom:webpconvert_1.0.0_sbom.cyclonedx.json --fail-on high
```

### Usando Trivy

[Trivy](https://github.com/aquasecurity/trivy) é outro scanner popular:

```bash
# Instalar Trivy (macOS)
brew install aquasecurity/trivy/trivy

# Instalar Trivy (Linux)
wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | sudo apt-key add -
echo "deb https://aquasecurity.github.io/trivy-repo/deb $(lsb_release -sc) main" | sudo tee -a /etc/apt/sources.list.d/trivy.list
sudo apt-get update && sudo apt-get install trivy

# Escanear SBOM
trivy sbom webpconvert_1.0.0_sbom.cyclonedx.json

# Gerar relatório em tabela
trivy sbom webpconvert_1.0.0_sbom.cyclonedx.json --format table

# Apenas vulnerabilidades críticas
trivy sbom webpconvert_1.0.0_sbom.cyclonedx.json --severity CRITICAL
```

### Usando Syft + Grype

Workflow completo de geração e análise:

```bash
# Gerar SBOM a partir do código fonte
syft dir:. -o cyclonedx-json=sbom.json

# Escanear com Grype
grype sbom:sbom.json
```

## Análise de Licenças

Extrair informações de licenciamento do SBOM:

```bash
# Listar todas as licenças (CycloneDX)
jq '.components[] | {name: .name, version: .version, licenses: .licenses}' \
  webpconvert_1.0.0_sbom.cyclonedx.json

# Encontrar dependências com licenças específicas
jq '.components[] | select(.licenses[]?.license?.id == "MIT")' \
  webpconvert_1.0.0_sbom.cyclonedx.json

# Contar por tipo de licença
jq -r '.components[].licenses[]?.license?.id' \
  webpconvert_1.0.0_sbom.cyclonedx.json | sort | uniq -c
```

## Verificação de Integridade

Cada release inclui checksums SHA256 para todos os SBOMs:

```bash
# Baixar checksums
wget https://github.com/robsonalvesdevbr/webpconvert/releases/download/v1.0.0/sbom_checksums.txt

# Baixar SBOM
wget https://github.com/robsonalvesdevbr/webpconvert/releases/download/v1.0.0/webpconvert_1.0.0_sbom.cyclonedx.json

# Verificar integridade
sha256sum -c sbom_checksums.txt --ignore-missing
```

## Integração com CI/CD

### GitHub Actions

```yaml
- name: Download SBOM
  run: |
    gh release download v1.0.0 -p "*sbom.cyclonedx.json"

- name: Scan SBOM
  uses: anchore/scan-action@v4
  with:
    sbom: webpconvert_1.0.0_sbom.cyclonedx.json
    fail-build: true
    severity-cutoff: high
```

### GitLab CI

```yaml
sbom-scan:
  image: anchore/grype:latest
  script:
    - wget https://github.com/robsonalvesdevbr/webpconvert/releases/download/v1.0.0/webpconvert_1.0.0_sbom.cyclonedx.json
    - grype sbom:webpconvert_1.0.0_sbom.cyclonedx.json --fail-on critical
```

### Jenkins

```groovy
stage('SBOM Scan') {
    steps {
        sh '''
            wget https://github.com/robsonalvesdevbr/webpconvert/releases/download/v1.0.0/webpconvert_1.0.0_sbom.cyclonedx.json
            docker run --rm -v $(pwd):/scan anchore/grype:latest sbom:/scan/webpconvert_1.0.0_sbom.cyclonedx.json
        '''
    }
}
```

## Automação de Segurança

O projeto webpconvert implementa automação completa de segurança:

### 1. Geração Automática em Releases

- SBOM gerado automaticamente para cada release
- Três formatos disponíveis (CycloneDX, SPDX, Syft)
- Checksums SHA256 para verificação de integridade

### 2. Submissão Contínua ao GitHub

- SBOM enviado ao GitHub Dependency Graph em cada push para `main`
- Executado semanalmente para detectar novas vulnerabilidades
- Scanning automatizado com Grype
- Relatórios no GitHub Security

### 3. Scanning Diário

- Execução diária de scans de vulnerabilidades
- Múltiplas ferramentas (Trivy + Grype)
- Issues automáticas para vulnerabilidades críticas/altas
- Relatórios SARIF no GitHub Security

### 4. Alertas Proativos

- Criação automática de issues para vulnerabilidades
- Atualização de issues existentes
- Integração com GitHub Dependabot
- Notificações de segurança

## Ferramentas e Recursos

### Geradores de SBOM

- **Syft** (Anchore): https://github.com/anchore/syft
- **CycloneDX GoMod**: https://github.com/CycloneDX/cyclonedx-gomod
- **SPDX Go Tools**: https://github.com/spdx/tools-golang

### Scanners de Vulnerabilidades

- **Grype** (Anchore): https://github.com/anchore/grype
- **Trivy** (Aqua Security): https://github.com/aquasecurity/trivy
- **Snyk**: https://snyk.io/
- **OWASP Dependency-Check**: https://owasp.org/www-project-dependency-check/

### Plataformas de Gestão

- **Dependency-Track**: https://dependencytrack.org/
- **GUAC**: https://guac.sh/
- **GitHub Dependency Graph**: Nativo no GitHub

## Padrões e Compliance

O webpconvert adere aos seguintes padrões:

### Frameworks

- **NIST SSDF** (Secure Software Development Framework)
- **US Executive Order 14028** (Improving the Nation's Cybersecurity)
- **CISA Minimum Elements for SBOM**
- **NTIA Minimum Elements for SBOM**

### Padrões

- **ISO/IEC 5962:2021** (SPDX)
- **CycloneDX 1.6**
- **SPDX 2.3**

## Perguntas Frequentes

### Como sei se há vulnerabilidades no meu binário?

```bash
# Baixe o SBOM da versão instalada
gh release download v1.0.0 -p "*sbom.cyclonedx.json"

# Escaneie
grype sbom:webpconvert_1.0.0_sbom.cyclonedx.json
```

### O SBOM inclui dependências de desenvolvimento?

Não. O SBOM reflete apenas dependências incluídas no binário de produção (`go.mod` com dependências diretas e transitivas).

### Com que frequência o SBOM é atualizado?

- **Releases**: SBOM gerado para cada release
- **Branch main**: Submissão ao GitHub Dependency Graph em cada push e semanalmente

### Posso confiar na integridade do SBOM?

Sim. Cada SBOM inclui:
- Checksum SHA256 verificável
- Geração automatizada em ambiente controlado (GitHub Actions)
- Logs de CI/CD públicos e auditáveis

### O SBOM funciona offline?

Sim. Baixe o arquivo SBOM e escaneie localmente com Grype ou Trivy. Bases de dados de vulnerabilidades podem ser baixadas para uso offline.

## Suporte

### Reportar Problemas com SBOM

Se encontrar problemas com SBOM:
- Abra uma issue: https://github.com/robsonalvesdevbr/webpconvert/issues
- Use a label `sbom`

### Questões de Segurança

Para reportar vulnerabilidades:
- Veja [SECURITY.md](../SECURITY.md)
- **Não** abra issues públicas para vulnerabilidades

### Contribuições

Contribuições para melhorar SBOM são bem-vindas:
- Sugestões de novos formatos
- Melhorias na automação
- Documentação adicional

## Recursos Adicionais

### Documentação Oficial

- **NTIA SBOM**: https://www.ntia.gov/SBOM
- **CISA SBOM**: https://www.cisa.gov/sbom
- **CycloneDX**: https://cyclonedx.org/
- **SPDX**: https://spdx.dev/

### Comunidade

- **OpenSSF**: https://openssf.org/
- **Sigstore**: https://www.sigstore.dev/
- **GUAC Project**: https://guac.sh/

### Artigos e Guias

- [SBOM at a Glance](https://www.cisa.gov/resources-tools/resources/sbom-glance)
- [NTIA Minimum Elements](https://www.ntia.gov/files/ntia/publications/sbom_minimum_elements_report.pdf)
- [CycloneDX Use Cases](https://cyclonedx.org/use-cases/)
