# Dependency Management Policy

## Objetivo

Esta política define como o projeto webpconvert gerencia suas dependências de software para garantir segurança, estabilidade e manutenibilidade.

## Princípios

1. **Segurança em Primeiro Lugar**: Vulnerabilidades conhecidas são tratadas com máxima prioridade
2. **Estabilidade**: Preferência por dependências maduras e bem mantidas
3. **Minimalismo**: Manter o menor número possível de dependências
4. **Transparência**: SBOM completo publicado em cada release
5. **Automação**: Processos automatizados para detecção e atualização

## Categorias de Dependências

### 1. Dependências de Runtime (Bibliotecas Nativas)

Bibliotecas C/C++ linkadas dinamicamente ao binário:

| Biblioteca | Versão Mínima | Propósito | Criticidade |
|------------|---------------|-----------|-------------|
| libwebp | 0.6.0+ | Decodificação WebP | Alta |
| giflib | 5.1.0+ | Encoding GIF | Alta |
| libjpeg-turbo | 2.0.0+ | Encoding JPEG | Alta |

**Gerenciamento**:
- Mantidas pelo sistema operacional
- Verificadas em CI/CD
- Documentadas em SBOM

### 2. Dependências Go Diretas

Pacotes Go declarados em `go.mod`:

| Pacote | Propósito | Última Verificação |
|--------|-----------|-------------------|
| github.com/chai2010/webp | Bindings Go para libwebp | Ativo |
| github.com/spf13/cobra | Framework CLI | Ativo |
| github.com/spf13/viper | Configuração | Ativo |

**Critérios de Aceitação**:
- ✅ Mantido ativamente (commits nos últimos 6 meses)
- ✅ Sem vulnerabilidades conhecidas de severidade alta ou crítica
- ✅ Licença compatível (MIT, Apache 2.0, BSD)
- ✅ Comunidade ativa (> 100 stars ou mantedor reconhecido)
- ✅ Cobertura de testes razoável

### 3. Dependências Transitivas

Pacotes importados indiretamente pelas dependências diretas.

**Gerenciamento**:
- Rastreadas automaticamente via SBOM
- Monitoradas para vulnerabilidades
- Atualizadas junto com dependências diretas

## Processo de Atualização

### Patches de Segurança (CRÍTICO/ALTO)

**Timeline**: Dentro de 48 horas

1. **Detecção**:
   - Alertas automáticos via GitHub Dependabot
   - Scanning diário com Grype/Trivy
   - Monitoramento de CVEs

2. **Avaliação**:
   - Verificar severidade (CVSS score)
   - Confirmar se componente afetado está em uso
   - Avaliar disponibilidade de patch

3. **Aplicação**:
   ```bash
   # Atualizar dependência específica
   go get -u github.com/pacote/vulneravel@latest
   go mod tidy

   # Verificar quebras
   go test ./...

   # Gerar novo SBOM
   syft . -o cyclonedx-json=sbom.json
   ```

4. **Validação**:
   - Testes automatizados
   - Scan de vulnerabilidades
   - Teste manual se necessário

5. **Release**:
   - Patch release (ex: v1.4.0 → v1.4.1)
   - Changelog detalhado
   - Comunicação clara da correção

### Atualizações Menores (MEDIUM)

**Timeline**: Mensalmente

1. **Revisão**:
   - Executar `go list -u -m all`
   - Revisar changelogs
   - Priorizar atualizações de segurança

2. **Teste**:
   ```bash
   # Atualizar minor versions
   go get -u=patch ./...
   go mod tidy

   # Testes completos
   go test -race ./...
   ```

3. **Integração**:
   - Branch separado para testes
   - CI/CD completo
   - Merge após validação

### Atualizações Maiores (LOW)

**Timeline**: Trimestralmente

1. **Planejamento**:
   - Avaliar breaking changes
   - Estimar esforço de migração
   - Verificar compatibilidade

2. **Implementação**:
   - Branch de desenvolvimento dedicado
   - Refatoração se necessário
   - Atualização de testes

3. **Validação**:
   - Testes extensivos
   - Performance benchmarks
   - Beta testing (se aplicável)

## Critérios de Remoção

Uma dependência será considerada para remoção se:

1. **Não mais mantida**:
   - Sem commits por > 12 meses
   - Maintainer abandonou o projeto
   - Issues/PRs não respondidos

2. **Vulnerabilidades persistentes**:
   - Vulnerabilidades críticas não corrigidas por > 30 dias
   - Histórico de problemas de segurança recorrentes

3. **Licença incompatível**:
   - Mudança de licença para termos incompatíveis
   - Descoberta de licença previamente não detectada

4. **Redundância**:
   - Funcionalidade substituída por stdlib
   - Substituível por código próprio com pouco esforço

**Processo**:
1. Avaliar alternativas
2. Implementar substituição
3. Testar extensivamente
4. Documentar migração
5. Release com changelog detalhado

## Análise de Licenças

### Licenças Permitidas

- ✅ MIT
- ✅ Apache 2.0
- ✅ BSD (2-clause, 3-clause)
- ✅ ISC
- ✅ MPL 2.0

### Licenças Restritas

- ⚠️ GPL/LGPL: Requer revisão caso a caso
- ⚠️ AGPL: Evitar quando possível

### Licenças Proibidas

- ❌ Licenças proprietárias sem permissão
- ❌ Licenças com restrições de uso comercial

**Verificação**:
```bash
# Extrair licenças do SBOM
jq -r '.components[].licenses[]?.license?.id' sbom.cyclonedx.json | sort | uniq

# Ferramentas adicionais
go-licenses report ./... --template=licenses.tpl
```

## Monitoramento de Vulnerabilidades

### Ferramentas Automatizadas

1. **GitHub Dependabot**:
   - Alertas automáticos
   - Pull requests de atualização
   - Integrado ao GitHub Security

2. **Grype** (Diário):
   ```bash
   grype dir:. --fail-on critical
   ```

3. **Trivy** (Diário):
   ```bash
   trivy fs . --severity CRITICAL,HIGH
   ```

### Workflow

```
Detecção → Triagem → Priorização → Correção → Validação → Release → Comunicação
```

### Níveis de Resposta

| Severidade | SLA | Ação |
|------------|-----|------|
| Critical (9.0-10.0) | 24h | Patch release imediato |
| High (7.0-8.9) | 48h | Patch release prioritário |
| Medium (4.0-6.9) | 7 dias | Incluir em próxima release |
| Low (0.1-3.9) | 30 dias | Agendar correção |

## SBOM (Software Bill of Materials)

### Geração

SBOM é gerado automaticamente:

1. **Em releases**:
   - CycloneDX JSON (primário)
   - SPDX JSON (compliance)
   - Syft JSON (advanced)

2. **Continuamente**:
   - Submissão ao GitHub Dependency Graph
   - Atualização semanal

### Conteúdo

Cada SBOM inclui:
- ✅ Todas as dependências Go (diretas e transitivas)
- ✅ Versões exatas (com hash de commit quando aplicável)
- ✅ Licenças identificadas
- ✅ Informações de fornecedor
- ✅ Metadados de build

### Acesso

```bash
# Baixar SBOM de release
gh release download v1.0.0 -p "*sbom.cyclonedx.json"

# Gerar SBOM localmente
syft . -o cyclonedx-json=sbom.json
```

## Auditoria e Compliance

### Auditoria Trimestral

Revisão completa de dependências a cada trimestre:

1. **Inventário**:
   - Validar SBOM contra go.mod/go.sum
   - Identificar dependências não utilizadas
   - Verificar versões desatualizadas

2. **Segurança**:
   - Scan completo de vulnerabilidades
   - Revisar CVEs resolvidos
   - Atualizar matriz de risco

3. **Licenças**:
   - Verificar conformidade de licenças
   - Identificar mudanças de licenciamento
   - Atualizar documentação legal

4. **Performance**:
   - Avaliar tamanho de dependências
   - Identificar bloat
   - Considerar otimizações

### Relatório

Gerar relatório contendo:
- Total de dependências (diretas/transitivas)
- Vulnerabilidades encontradas e corrigidas
- Atualizações aplicadas
- Licenças em uso
- Recomendações

## Adição de Novas Dependências

### Processo de Aprovação

Antes de adicionar uma nova dependência, avaliar:

1. **Necessidade**:
   - [ ] Funcionalidade não disponível em stdlib?
   - [ ] Benefício justifica complexidade adicional?
   - [ ] Alternativas consideradas?

2. **Qualidade**:
   - [ ] Projeto ativo (commits recentes)?
   - [ ] Testes adequados (> 70% coverage)?
   - [ ] Documentação completa?
   - [ ] Comunidade ativa?

3. **Segurança**:
   - [ ] Sem vulnerabilidades conhecidas?
   - [ ] Histórico de segurança limpo?
   - [ ] Processo de disclosure responsável?

4. **Licença**:
   - [ ] Licença compatível?
   - [ ] Claramente identificada?
   - [ ] Sem restrições problemáticas?

5. **Manutenção**:
   - [ ] Maintainer responsivo?
   - [ ] Roadmap claro?
   - [ ] Breaking changes raros?

### Documentação

Ao adicionar dependência:

1. Atualizar `go.mod`:
   ```bash
   go get github.com/org/package@version
   ```

2. Documentar em commit:
   ```
   feat: add package X for Y functionality

   - Evaluated alternatives: A, B, C
   - Chosen for: performance/simplicity/features
   - License: MIT
   - Last security audit: 2024-01-01
   ```

3. Atualizar SBOM:
   ```bash
   syft . -o cyclonedx-json=sbom.json
   ```

## Ferramentas

### Gestão de Dependências

```bash
# Listar dependências
go list -m all

# Listar atualizações disponíveis
go list -u -m all

# Atualizar todas (minor/patch)
go get -u ./...

# Atualizar específica
go get -u github.com/org/package@latest

# Limpar não utilizadas
go mod tidy
```

### Análise de Vulnerabilidades

```bash
# Grype (recomendado)
grype dir:.

# Trivy
trivy fs .

# Go built-in (básico)
go list -json -m all | grype
```

### Análise de Licenças

```bash
# go-licenses (Google)
go install github.com/google/go-licenses@latest
go-licenses report ./... --template=licenses.tpl

# SBOM
jq -r '.components[].licenses[]?.license?.id' sbom.cyclonedx.json
```

## Exceções

Exceções a esta política requerem:

1. **Justificativa documentada**
2. **Aprovação de mantenedor**
3. **Plano de mitigação de riscos**
4. **Revisão periódica**

Documentar em `EXCEPTIONS.md`:
```markdown
## Exceção: Dependência X

**Razão**: Y
**Risco**: Z
**Mitigação**: W
**Revisão**: Trimestral
**Aprovado por**: Mantenedor (@github)
**Data**: 2024-01-01
```

## Contato

Para questões sobre dependências:
- **Vulnerabilidades**: Veja [SECURITY.md](SECURITY.md)
- **Outras questões**: Abra issue com label `dependencies`

## Recursos

### Ferramentas
- **Syft**: https://github.com/anchore/syft
- **Grype**: https://github.com/anchore/grype
- **Trivy**: https://github.com/aquasecurity/trivy
- **go-licenses**: https://github.com/google/go-licenses

### Padrões
- **CycloneDX**: https://cyclonedx.org/
- **SPDX**: https://spdx.dev/
- **CVSS**: https://www.first.org/cvss/

### Comunidade
- **OpenSSF**: https://openssf.org/
- **Go Security**: https://go.dev/security/

## Changelog

### 2024-01-09
- ✅ Política inicial criada
- ✅ SBOM automático implementado
- ✅ Scanning automatizado configurado
- ✅ Processo de resposta a vulnerabilidades definido

### Futuro
- 🔄 Dependabot auto-merge para patches
- 🔄 SBOM signing com Sigstore
- 🔄 Supply chain attestations (SLSA)
