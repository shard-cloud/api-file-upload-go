# Rodando o projeto

## Desenvolvimento

### Linux/macOS
```bash
# Instalar dependências
go mod tidy

# Configurar variáveis de ambiente
cp env.example .env

# Build do binário
go build -o api-file-upload-go ./cmd/main.go

# Executar API
./api-file-upload-go
```

### Windows
```bash
# Instalar dependências
go mod tidy

# Configurar variáveis de ambiente
copy env.example .env

# Build do binário
go build -o api-file-upload-go.exe ./cmd/main.go

# Executar API
.\api-file-upload-go.exe
```

## Produção

### Linux/macOS
```bash
# Build do binário
go build -o api-file-upload-go ./cmd/main.go

# Executar API
./api-file-upload-go
```

### Windows
```bash
# Build do binário
go build -o api-file-upload-go.exe ./cmd/main.go

# Executar API
.\api-file-upload-go.exe
```

## Testando a API

### Health Check
```bash
curl http://localhost:80/health
```

### Upload de arquivo
```bash
curl -X POST -F "file=@document.pdf" http://localhost:80/api/v1/files/upload
```

### Listar arquivos
```bash
curl http://localhost:80/api/v1/files
```

### Download de arquivo
```bash
curl http://localhost:80/api/v1/files/1/download -o downloaded_file.pdf
```

### Obter detalhes do arquivo
```bash
curl http://localhost:80/api/v1/files/1
```

### Estatísticas
```bash
curl http://localhost:80/api/v1/stats
```

### Deletar arquivo
```bash
curl -X DELETE http://localhost:80/api/v1/files/1
```

## Endpoints da API

### File Operations
- `POST /api/v1/files/upload` – Upload de arquivo
- `GET /api/v1/files` – Listar arquivos (com paginação)
- `GET /api/v1/files/:id` – Obter detalhes do arquivo
- `GET /api/v1/files/:id/download` – Download do arquivo
- `DELETE /api/v1/files/:id` – Deletar arquivo

### Statistics
- `GET /api/v1/stats` – Estatísticas de upload

### System
- `GET /health` – Health check
- `GET /` – Informações da API

## Parâmetros de Query

### Listar arquivos
- `limit` (opcional): Número máximo de arquivos (padrão: 10, máximo: 100)
- `offset` (opcional): Número de arquivos para pular (padrão: 0)

Exemplo:
```bash
curl "http://localhost:80/api/v1/files?limit=5&offset=10"
```

## Respostas da API

### Upload bem-sucedido
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "file": {
    "id": 1,
    "name": "document.pdf",
    "size": 1024,
    "mime_type": "application/pdf",
    "extension": ".pdf",
    "hash": "abc123...",
    "uploaded_at": "2025-10-26T14:55:13Z"
  }
}
```

### Lista de arquivos
```json
{
  "success": true,
  "data": {
    "files": [
      {
        "id": 1,
        "name": "document.pdf",
        "size": 1024,
        "mime_type": "application/pdf",
        "extension": ".pdf",
        "hash": "abc123...",
        "uploaded_at": "2025-10-26T14:55:13Z",
        "download_url": "/api/v1/files/1/download"
      }
    ],
    "total": 1,
    "limit": 10,
    "offset": 0
  }
}
```

### Estatísticas
```json
{
  "success": true,
  "data": {
    "total_files": 10,
    "total_size": 1048576,
    "recent_uploads": 3,
    "largest_file": {
      "name": "large_file.pdf",
      "size": 524288
    },
    "extension_stats": [
      {
        "extension": ".pdf",
        "count": 5,
        "size": 524288
      }
    ]
  }
}
```

## Logs

A API gera logs estruturados com informações sobre:
- Requisições HTTP
- Uploads de arquivos
- Erros de validação
- Operações de banco de dados
- Health checks

## Testes

### Testes unitários
```bash
# Executar testes
go test ./...

# Executar com cobertura
go test -cover ./...
```

### Testes da API
```bash
# Iniciar a API primeiro
./api-file-upload-go  # Linux/macOS
# ou
.\api-file-upload-go.exe  # Windows

# Executar testes da API (em outro terminal)
go test ./tests/
```

## Troubleshooting

### Erro de conexão com banco
Verifique se a variável `DATABASE` está configurada corretamente:
```bash
echo $DATABASE  # Linux/macOS
echo %DATABASE%  # Windows
```

### Erro de porta em uso
Verifique se a porta está livre:
```bash
netstat -ano | findstr :80  # Windows
netstat -tulpn | grep :80   # Linux
lsof -i :80                 # macOS
```

### Erro de permissão de arquivo
Verifique se o diretório `UPLOAD_DIR` existe e tem permissões de escrita:
```bash
mkdir -p uploads
chmod 755 uploads
```