# Introdução

A File Upload API é uma API REST para upload de arquivos com rastreamento de metadados em banco de dados. Ela permite gerenciar arquivos de forma eficiente via endpoints HTTP, com controle de tamanho, tipos permitidos e verificação de duplicatas.

## Funcionalidades

- Upload de arquivos via API REST
- Controle de tamanho máximo
- Verificação de tipos de arquivo permitidos
- Detecção de duplicatas por hash MD5
- Listagem de arquivos com paginação
- Download de arquivos via API
- Estatísticas de upload
- Exclusão segura de arquivos
- Logs estruturados
- Health check endpoint
- Suporte a CORS

## Casos de uso

- Backup de arquivos via API
- Gerenciamento de documentos em aplicações web
- Upload de imagens para galerias
- Armazenamento de configurações
- Versionamento de arquivos
- Auditoria de uploads
- Integração com frontends web/mobile

## Tecnologias

- **Linguagem**: Go 1.21+
- **Framework Web**: Gin
- **Banco de dados**: PostgreSQL com GORM
- **Configuração**: Environment variables
- **Logs**: Logrus
- **Testes**: Testify
- **Containerização**: Docker

## Arquitetura

A API segue uma arquitetura limpa com separação de responsabilidades:

- **Handlers**: Gerenciam requisições HTTP
- **Models**: Definem estruturas de dados
- **Database**: Gerencia conexão e migrações
- **Config**: Centraliza configurações
- **Utils**: Funções utilitárias
- **Logger**: Sistema de logs estruturados