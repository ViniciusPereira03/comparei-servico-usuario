# comparei-servico-usuario

Este é um serviço backend desenvolvido em Go, que se conecta ao banco de dados MySQL e Redis para gerenciar usuários.

## **Tecnologias Utilizadas**

- **Go 1.x** – Linguagem de programação principal.
- **MySQL 8.x** – Banco de dados relacional para armazenar informações.
- **Redis 7.x** – Sistema de armazenamento em cache, utilizado para melhorar a performance.
- **GORM** – ORM para interação com o banco de dados MySQL.
- **Docker** – Contêineres para facilitar o desenvolvimento e a implantação.

## **Estrutura do Projeto**

Este projeto segue a seguinte estrutura de pastas:
```
/comparei-servico-usuario
  ├── /cmd                   
  │     └── main.go          # Arquivo principal para inicialização
  │
  ├── /config                
  │     ├── config.go        # Carregamento de configurações (e.g., env)
  │     ├── config.yaml      # Configurações padrões (e.g., para docker-compose)
  │
  ├── /internal              
  │     ├── /db              # Inicialização e interações com o banco de dados
  │     │     ├── mysql.go   # Conexão com MySQL
  │     │     ├── redis.go   # Conexão com Redis
  │     │
  │     ├── /cache           # Serviços relacionados ao cache
  │     │     └── cache.go   # Funções de gerenciamento de cache
  │     │
  │     ├── /repository      # Repositórios para acesso a dados
  │     │     ├── user.go    # Operações no banco para usuários
  │     │
  │     ├── /service         # Lógica de negócios
  │     │     ├── user.go    # Serviço de usuários
  │     │
  │     └── /api             # Handlers e rotas de API
  │           ├── user.go    # Rotas da API para usuários
  │
  ├── /pkg                   
  │     ├── /logger          # Log customizado
  │     ├── /utils           # Funções utilitárias (e.g., validações, helpers)
  │
  ├── /test                  
  │     ├── integration_test.go
  │     └── e2e_test.go
  │
  ├── Dockerfile             
  ├── docker-compose.yml     # Configuração do Docker Compose (serviço, MySQL, Redis)
  ├── go.mod                 
  ├── go.sum                 
  ├── README.md              
  └── .env                   # Variáveis de ambiente
```

## **Pré-Requisitos**

Antes de rodar o serviço, você precisa ter os seguintes programas instalados:

- **Go** (1.x)
- **Docker** (para rodar os containers MySQL e Redis)
- **Docker Compose** (opcional, para facilitar a orquestração de containers)

Além disso, crie um arquivo `.env` com as configurações de ambiente (veja abaixo).

## **Configuração de Ambiente**

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis de ambiente:
 ```env
    APP_PORT=8080

    MYSQL_ROOT_PASSWORD=senha_root
    MYSQL_DATABASE=users
    DB_USER=servico_usuario
    DB_PASSWORD=senha_db
    DB_HOST=mysql_db
    DB_PORT=3306

    REDIS_HOST=redis_cache
    REDIS_PORT=6379
 ```

 Estas configurações são usadas para configurar as conexões com o MySQL e o Redis. Você pode alterar essas variáveis conforme necessário para seu ambiente de desenvolvimento ou produção.

## **Rodando o Projeto com Docker**

### 1. **Construa e inicie os containers com Docker Compose**

Se você estiver usando o Docker Compose, pode rodar o seguinte comando para subir o ambiente de desenvolvimento com MySQL e Redis:

```bash
docker-compose up --build
```

Isso criará e iniciará os containers necessários, e o serviço será acessível localmente.

### 2. Subir somente o MySQL e Redis manualmente
Se você preferir rodar os containers manualmente, execute os seguintes comandos:

**Iniciar o Redis:**
```bash
docker run -d --name redis_cache -p 6379:6379 redis:7.0
```

**Iniciar o MySQL:**
```bash
docker run -d --name mysql_db -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=usuario_service -p 3306:3306 mysql:8.0

```

### 3. Rodar o serviço Go
Após subir os containers, execute o seguinte comando para rodar o serviço Go:
```bash
go run cmd/main.go
```

## **Executando as Migrações**

O serviço realiza migrações automaticamente quando é iniciado, para garantir que o banco de dados esteja configurado corretamente. O script de migração de criação da tabela de usuários está localizado em `/migrations/create_users_table.sql`.

## **Endpoints Disponíveis**
- **GET** `/user/{id}` – Obter detalhes de um usuário pelo ID.