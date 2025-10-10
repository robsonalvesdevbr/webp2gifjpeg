# Security Policy

## Supported Versions

Usamos este projeto ativamente e fornecemos atualizações de segurança para as seguintes versões:

| Version | Supported          |
| ------- | ------------------ |
| 1.4.x   | :white_check_mark: |
| 1.3.x   | :white_check_mark: |
| < 1.3   | :x:                |

## Reporting a Vulnerability

A segurança do webpconvert é levada a sério. Se você descobrir uma vulnerabilidade de segurança, por favor siga as diretrizes abaixo.

### Como Reportar

**⚠️ NÃO abra issues públicas para vulnerabilidades de segurança.**

Para reportar vulnerabilidades de segurança:

1. **Email**: Envie um email para o mantenedor do projeto através do GitHub
2. **GitHub Security Advisory**: Use a funcionalidade [Security Advisories](https://github.com/robsonalvesdevbr/webpconvert/security/advisories/new) (recomendado)

### O que Incluir

Por favor inclua o máximo de informação possível:

- Tipo de vulnerabilidade (ex: buffer overflow, SQL injection, cross-site scripting, etc.)
- Caminho completo dos arquivos relacionados à manifestação da vulnerabilidade
- Localização do código vulnerável (tag/branch/commit ou URL direto)
- Configuração especial necessária para reproduzir o problema
- Passos detalhados para reproduzir a vulnerabilidade
- Proof-of-concept ou código de exploit (se possível)
- Impacto potencial da vulnerabilidade, incluindo como um atacante poderia explorá-la

### O que Esperar

Após reportar uma vulnerabilidade:

1. **Confirmação**: Você receberá uma confirmação dentro de 48 horas
2. **Avaliação**: Avaliaremos o problema e determinaremos sua severidade dentro de 7 dias
3. **Comunicação**: Manteremos você informado sobre o progresso da correção
4. **Resolução**: Trabalharemos para corrigir a vulnerabilidade o mais rápido possível
5. **Divulgação**: Coordenaremos a divulgação pública após a correção

### Severidade

Classificamos vulnerabilidades usando o [CVSS 3.1](https://www.first.org/cvss/v3.1/specification-document):

- **Critical** (9.0-10.0): Correção imediata (< 24 horas)
- **High** (7.0-8.9): Correção prioritária (< 48 horas)
- **Medium** (4.0-6.9): Correção na próxima release (< 7 dias)
- **Low** (0.1-3.9): Correção planejada (< 30 dias)

## Security Features

### Software Bill of Materials (SBOM)

O webpconvert publica SBOM completo para todas as releases, permitindo:

- **Rastreabilidade**: Inventário completo de todas as dependências
- **Detecção Rápida**: Identificação imediata de componentes vulneráveis
- **Automação**: Scanning automatizado de vulnerabilidades
- **Compliance**: Atendimento a frameworks de segurança

#### Acessar SBOM

```bash
# Baixar SBOM de uma release
gh release download v1.0.0 -p "*sbom.cyclonedx.json"

# Escanear vulnerabilidades
grype sbom:webpconvert_1.0.0_sbom.cyclonedx.json
```

**Documentação completa**: [docs/SBOM.md](docs/SBOM.md)

### Automação de Segurança

O projeto implementa as seguintes medidas automatizadas:

#### 1. Geração Automática de SBOM
- SBOM gerado para cada release em 3 formatos (CycloneDX, SPDX, Syft)
- Checksums SHA256 para verificação de integridade
- Disponível como release assets

#### 2. Dependency Tracking
- Submissão automática ao GitHub Dependency Graph
- Integração com Dependabot para alertas de segurança
- Atualização semanal de análise de dependências

#### 3. Vulnerability Scanning
- **Diário**: Scan completo com Grype e Trivy
- **Contínuo**: Análise em cada push para branch main
- **Release**: Validação de segurança antes de cada release

#### 4. Security Alerts
- Criação automática de issues para vulnerabilidades críticas/altas
- Relatórios SARIF no GitHub Security tab
- Notificações via GitHub

#### 5. GitHub Security Features
- Dependabot alerts habilitado
- Security advisories habilitado
- Code scanning com múltiplas ferramentas

### Dependency Management

Monitoramos e atualizamos dependências regularmente:

- **Patches de Segurança**: Dentro de 48 horas
- **Atualizações Menores**: Mensalmente
- **Atualizações Maiores**: Trimestralmente (com testes de compatibilidade)

Veja nossa [Política de Dependências](DEPENDENCY_POLICY.md) para detalhes.

### Build Security

#### Pipeline de CI/CD

- Builds executados em ambiente isolado (GitHub Actions)
- Verificação de checksums para dependências
- Logs públicos e auditáveis
- Assinatura de releases (planejado)

#### Binários

- Compilação com flags de segurança
- Bibliotecas nativas verificadas (libwebp, giflib, libjpeg)
- Checksums SHA256 para todos os binários
- Verificação de integridade disponível

### Runtime Security

#### CGO Safety

O projeto usa CGO para integração com bibliotecas C nativas:
- Bibliotecas auditadas: libwebp, giflib, libjpeg-turbo
- Binding seguro Go ↔ C
- Sem execução de código arbitrário
- Memory safety através de validações

#### Input Validation

- Validação de formatos de arquivo
- Sanitização de paths
- Proteção contra path traversal
- Tratamento robusto de erros

## Security Best Practices

### Para Usuários

1. **Sempre baixe de fontes oficiais**:
   - GitHub Releases: https://github.com/robsonalvesdevbr/webpconvert/releases
   - Verifique checksums SHA256

2. **Mantenha atualizado**:
   - Assine releases no GitHub para notificações
   - Use a versão mais recente para correções de segurança

3. **Verifique integridade**:
   ```bash
   # Baixar checksum
   wget https://github.com/robsonalvesdevbr/webpconvert/releases/download/v1.0.0/webpconvert-v1.0.0-linux-amd64.tar.gz.sha256

   # Verificar
   sha256sum -c webpconvert-v1.0.0-linux-amd64.tar.gz.sha256
   ```

4. **Escaneie vulnerabilidades**:
   ```bash
   # Baixar SBOM
   gh release download v1.0.0 -p "*sbom.cyclonedx.json"

   # Escanear
   grype sbom:webpconvert_1.0.0_sbom.cyclonedx.json
   ```

### Para Desenvolvedores

1. **Atualize dependências regularmente**:
   ```bash
   go get -u ./...
   go mod tidy
   ```

2. **Execute testes de segurança**:
   ```bash
   # Race detector
   go test -race ./...

   # Static analysis
   go vet ./...
   ```

3. **Valide builds**:
   ```bash
   # Verificar dependências
   go list -m all

   # Audit de vulnerabilidades
   go list -json -m all | grype
   ```

## Compliance

Este projeto adere a:

- **NIST Secure Software Development Framework (SSDF)**
- **Executive Order 14028** (Improving the Nation's Cybersecurity)
- **CISA Minimum Elements for SBOM**
- **NTIA Minimum Elements for SBOM**

## Security Contact

Para questões de segurança não relacionadas a vulnerabilidades:

- Abra uma issue com label `security`
- Entre em contato através do GitHub

## Acknowledgments

Agradecemos aos pesquisadores de segurança que ajudam a manter o webpconvert seguro através de divulgação responsável.

## Security Changelog

### Implemented
- ✅ SBOM generation and publication (v1.4.0)
- ✅ Automated vulnerability scanning (v1.4.0)
- ✅ GitHub Security integration (v1.4.0)
- ✅ Dependency tracking automation (v1.4.0)
- ✅ Security policy documentation (v1.4.0)

### Planned
- 🔄 Release signing with Sigstore/Cosign
- 🔄 SBOM signing for integrity verification
- 🔄 Supply chain level attestations (SLSA)
- 🔄 Automated dependency updates (Dependabot auto-merge)

## Resources

### Tools
- **Grype**: https://github.com/anchore/grype
- **Trivy**: https://github.com/aquasecurity/trivy
- **Syft**: https://github.com/anchore/syft

### Standards
- **SBOM**: https://www.cisa.gov/sbom
- **NIST SSDF**: https://csrc.nist.gov/Projects/ssdf
- **CVSS**: https://www.first.org/cvss/

### Community
- **OpenSSF**: https://openssf.org/
- **Sigstore**: https://www.sigstore.dev/
