## 🔐 Variáveis de Ambiente

Crie um arquivo `.env` baseado no `env.example`:

```bash
cp env.example .env
```

### Arquivo `.env` completo

```env
# PostgreSQL (obrigatório)
DATABASE=postgres://user:pass@localhost:5432/uploader_db

# Upload configuration
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=10485760
ALLOWED_EXTENSIONS=.jpg,.jpeg,.png,.gif,.pdf,.txt,.doc,.docx

# CLI configuration
LOG_LEVEL=info
```

### Variáveis Detalhadas

#### `DATABASE` (obrigatório)

String de conexão PostgreSQL (suporta tanto `postgres://` quanto `postgresql://`):

```env
DATABASE=postgres://USER:PASSWORD@HOST:PORT/DATABASE
```

Exemplos:

```env
# Local
DATABASE=postgres://postgres:postgres@localhost:5432/uploader_db

# Docker Compose
DATABASE=postgres://uploader:uploader@db:5432/uploader_db

# Supabase
DATABASE=postgres://user:pass@db.xxx.supabase.co:5432/postgres

# Railway
DATABASE=postgres://user:pass@containers-us-west-1.railway.app:5432/railway

# Shard Cloud (com SSL)
DATABASE=postgres://user:pass@postgres.shardatabases.app:5432/database?ssl=true
```

#### `UPLOAD_DIR` (opcional, padrão: ./uploads)

Diretório onde os arquivos serão salvos:

```env
UPLOAD_DIR=./uploads
UPLOAD_DIR=/var/uploads
UPLOAD_DIR=/app/storage
```

#### `MAX_FILE_SIZE` (opcional, padrão: 10485760)

Tamanho máximo de arquivo em bytes (10MB por padrão):

```env
MAX_FILE_SIZE=10485760    # 10MB
MAX_FILE_SIZE=52428800    # 50MB
MAX_FILE_SIZE=104857600   # 100MB
```

#### `ALLOWED_EXTENSIONS` (opcional)

Extensões de arquivo permitidas separadas por vírgula:

```env
ALLOWED_EXTENSIONS=.jpg,.jpeg,.png,.gif,.pdf,.txt,.doc,.docx
ALLOWED_EXTENSIONS=.pdf,.doc,.docx,.txt
ALLOWED_EXTENSIONS=.jpg,.jpeg,.png,.gif,.webp
```

#### `LOG_LEVEL` (opcional, padrão: info)

Nível de log do aplicativo:

```env
LOG_LEVEL=info     # Opções: debug, info, warn, error
LOG_LEVEL=warn     # Apenas avisos e erros
LOG_LEVEL=error    # Apenas erros
```

## 🗄️ Banco de Dados

### Opção 1: Docker Compose (Recomendado)

```bash
docker-compose up -d db
```

Credenciais padrão:
- **User:** uploader
- **Password:** uploader
- **Database:** uploader_db
- **Port:** 5432

### Opção 2: PostgreSQL Local

```bash
# Criar usuário e banco
psql -U postgres
CREATE USER uploader WITH PASSWORD 'uploader';
CREATE DATABASE uploader_db OWNER uploader;
GRANT ALL PRIVILEGES ON DATABASE uploader_db TO uploader;
```

### Opção 3: PostgreSQL em Cloud

**Supabase (Grátis):**
1. Crie projeto em https://supabase.com
2. Vá em Settings > Database
3. Copie Connection String
4. Cole no `.env`

**Railway:**
1. Crie projeto em https://railway.app
2. Adicione PostgreSQL plugin
3. Copie `DATABASE_URL`

## 🔄 Migrations

O CLI usa **GORM AutoMigrate** para criar/atualizar tabelas automaticamente:

```go
// Auto migrate
if err := db.AutoMigrate(&models.File{}); err != nil {
    return nil, fmt.Errorf("failed to migrate database: %w", err)
}
```

**Tabela `files` criada automaticamente:**
- `id` (Primary Key)
- `name` (Nome do arquivo)
- `original_name` (Nome original)
- `path` (Caminho do arquivo)
- `size` (Tamanho em bytes)
- `mime_type` (Tipo MIME)
- `extension` (Extensão)
- `hash` (Hash MD5 único)
- `uploaded_at` (Data de upload)
- `updated_at` (Data de atualização)
- `deleted_at` (Soft delete)

## 🔐 Configuração de Segurança

### Validação de Arquivos

O CLI valida arquivos antes do upload:

1. **Tamanho:** Verifica se não excede `MAX_FILE_SIZE`
2. **Extensão:** Verifica se está em `ALLOWED_EXTENSIONS`
3. **Duplicação:** Verifica hash MD5 para evitar duplicatas
4. **Existência:** Verifica se o arquivo existe no sistema

### Hash MD5

Cada arquivo é identificado por seu hash MD5:

```go
func calculateFileHash(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }

    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
```

### Soft Delete

Arquivos são removidos usando soft delete (campo `deleted_at`):

```go
type File struct {
    // ... outros campos
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
```

## 🐳 Docker

### Build Customizado

```bash
# Build da imagem
docker build -t cli-uploader-go .

# Run com variáveis
docker run --rm -v $(pwd)/uploads:/app/uploads \
  -e DATABASE=postgres://user:pass@host:5432/db \
  -e UPLOAD_DIR=/app/uploads \
  cli-uploader-go upload --file /path/to/file
```

### Docker Compose Personalizado

```yaml
version: '3.8'
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-uploader}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-uploader}
      POSTGRES_DB: ${DB_NAME:-uploader_db}
    volumes:
      - postgres_data:/var/lib/postgresql/data

  cli:
    build: .
    environment:
      DATABASE: postgres://${DB_USER:-uploader}:${DB_PASSWORD:-uploader}@db:5432/${DB_NAME:-uploader_db}
      UPLOAD_DIR: /app/uploads
    volumes:
      - ./uploads:/app/uploads
    depends_on:
      - db
```

## 🔧 Configuração Avançada

### Logs Estruturados

O CLI usa **logrus** para logs estruturados:

```go
log := logger.New(cfg.LogLevel)
log.Info("File uploaded successfully:", file.OriginalName)
log.Error("Failed to upload file:", err)
```

### Connection Pool

GORM gerencia automaticamente o pool de conexões:

```go
db, err := gorm.Open(postgres.Open(convertedURL), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Error),
})
```

## 🎯 Próximos passos

Continue para [Rodando](./03-rodando.md) para executar e testar o CLI.