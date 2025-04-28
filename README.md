# Aplicação Go com MySQL

Esta é uma aplicação em Go que se conecta a um banco de dados MySQL e fornece uma API REST para gerenciar usuários, tipos de perfil, seguradoras, eventos, objetos de contabilização e sistemas contábeis, com autenticação JWT, limite de tentativas de login, rotação de tokens, auditoria e diversas medidas de segurança.

## Estrutura do Projeto

\`\`\`
EstudoGo/
├── .env                    # Variáveis de ambiente
├── go.mod                  # Definição do módulo e dependências
├── main.go                 # Ponto de entrada da aplicação
└── internal/               # Código interno da aplicação
    ├── auth/               # Autenticação JWT
    │   └── jwt.go
    ├── config/             # Configurações da aplicação
    │   └── config.go
    ├── database/           # Conexão com o banco de dados
    │   ├── database.go
    │   └── seed.go
    ├── docs/               # Documentação Swagger
    │   └── swagger.go
    ├── handlers/           # Manipuladores HTTP
    │   ├── handlers.go
    │   ├── auth_handler.go
    │   ├── tipo_perfil_handler.go
    │   ├── seguradora_handler.go
    │   ├── evento_handler.go
    │   ├── objeto_contabilizacao_handler.go
    │   ├── objeto_contabilizacao_evento_handler.go
    │   ├── sistema_contabil_handler.go
    │   ├── sistema_contabil_config_handler.go
    │   └── swagger_handler.go
    ├── middleware/         # Middlewares
    │   └── auth_middleware.go
    ├── models/             # Modelos de dados
    │   ├── usuario.go
    │   ├── tipo_perfil.go
    │   ├── seguradora.go
    │   ├── evento.go
    │   ├── objeto_contabilizacao.go
    │   ├── objeto_contabilizacao_evento.go
    │   ├── sistema_contabil.go
    │   └── sistema_contabil_config.go
    ├── security/           # Componentes de segurança
    │   ├── csrf.go
    │   ├── rate_limiter.go
    │   ├── security_headers.go
    │   └── password_policy.go
    ├── services/           # Serviços da aplicação
    │   └── audit_service.go
    └── utils/              # Utilitários
        ├── validator.go
        └── sanitizer.go
\`\`\`

## Requisitos

- Go 1.16 ou superior
- Acesso ao banco de dados MySQL (Hostgator)

## Configuração

1. Clone o repositório
2. O arquivo `.env` já está configurado com os dados de conexão ao banco
3. Execute `go mod tidy` para baixar as dependências
4. Execute `go run main.go` para iniciar a aplicação

### Resolvendo Problemas com Dependências

Se você encontrar erros relacionados a dependências ausentes no arquivo `go.sum`, execute:

\`\`\`bash
go mod tidy
\`\`\`

Se o problema persistir, você pode instalar manualmente as dependências específicas:

\`\`\`bash
go get golang.org/x/crypto/bcrypt
go get golang.org/x/tools/go/loader
go get golang.org/x/net/webdav
\`\`\`

### Gerando Documentação Swagger

Para gerar ou atualizar a documentação Swagger, execute:

\`\`\`bash
swag init -g main.go
\`\`\`

## Recursos de Segurança

A API implementa diversos recursos de segurança para proteger os dados e os usuários:

### 1. Autenticação e Autorização

- **Autenticação JWT**: Utiliza tokens JWT para autenticar usuários
- **Rotação de Tokens**: Implementa renovação automática de tokens antes da expiração
- **Refresh Tokens**: Utiliza tokens de refresh para permitir renovação de sessões sem reautenticação
- **Controle de Acesso Baseado em Perfil**: Restringe acesso a recursos com base no perfil do usuário

### 2. Proteção Contra Ataques

- **Limite de Tentativas de Login**: Bloqueia temporariamente contas após múltiplas tentativas de login malsucedidas
- **Rate Limiting**: Limita o número de requisições por IP para prevenir ataques de força bruta e DoS
- **Proteção CSRF**: Implementa tokens CSRF para prevenir ataques Cross-Site Request Forgery
- **Headers de Segurança HTTP**: Configura cabeçalhos de segurança para prevenir diversos ataques
- **Sanitização de Entrada**: Valida e sanitiza todas as entradas para prevenir injeção SQL e XSS
- **Proteção Contra Enumeração de Usuários**: Evita vazamento de informações sobre existência de usuários

### 3. Política de Senhas

- **Requisitos Robustos**: Exige senhas fortes com letras maiúsculas, minúsculas, números e caracteres especiais
- **Verificação de Senhas Comuns**: Bloqueia o uso de senhas conhecidas e facilmente adivináveis
- **Expiração de Senhas**: Força a troca periódica de senhas para maior segurança
- **Histórico de Senhas**: Impede a reutilização de senhas anteriores

### 4. Auditoria e Monitoramento

- **Log de Auditoria**: Registra todas as operações críticas (login, alterações de dados sensíveis) em um log de auditoria
- **Rastreamento de IP**: Registra os endereços IP de todas as requisições para fins de auditoria
- **Monitoramento de Atividades Suspeitas**: Detecta e registra padrões de comportamento potencialmente maliciosos

## Autenticação

A API utiliza autenticação JWT (JSON Web Token) para proteger as rotas. Para acessar as rotas protegidas, é necessário:

1. Fazer login através do endpoint `/auth/login` para obter tokens de acesso e refresh
2. Incluir o token de acesso no cabeçalho `Authorization` das requisições no formato `Bearer {token}`
3. Quando o token estiver próximo de expirar, a API automaticamente fornecerá um novo token no cabeçalho de resposta
4. Se o token expirar, use o endpoint `/auth/refresh` com o refresh token para obter novos tokens

### Endpoints de Autenticação

- `POST /auth/login` - Realiza login e retorna tokens JWT
  - Corpo da requisição: `{ "login": "seu_login", "senha": "sua_senha" }`
  - Resposta: `{ "access_token": "jwt_token", "refresh_token": "refresh_token", "usuario": {...}, "expires_in": 86400 }`

- `POST /auth/refresh` - Renova um token JWT válido
  - Corpo da requisição: `{ "refresh_token": "seu_refresh_token" }`
  - Resposta: `{ "access_token": "novo_jwt_token", "refresh_token": "novo_refresh_token", "expires_in": 86400 }`

### Limite de Tentativas de Login

Para proteger contra ataques de força bruta, a API implementa um limite de tentativas de login:

- Após 5 tentativas falhas em 15 minutos, a conta será bloqueada por 30 minutos
- Durante o bloqueio, qualquer tentativa de login resultará em erro 429 (Too Many Requests)
- O tempo restante de bloqueio é informado na resposta

## Documentação da API (Swagger)

A API possui documentação interativa usando Swagger. Para acessar:

1. Inicie a aplicação com `go run main.go`
2. Acesse `http://localhost:8080/swagger/index.html` no navegador

A documentação Swagger permite:
- Visualizar todos os endpoints disponíveis
- Testar as requisições diretamente pela interface
- Verificar os modelos de dados e parâmetros necessários
- Entender as respostas esperadas para cada operação

## Modelos de Dados

### Usuário
- `id` - Identificador único (auto-incremento)
- `nome` - Nome do usuário
- `email` - Email do usuário (único)
- `login` - Login do usuário (único)
- `senha` - Senha do usuário (armazenada com hash bcrypt)
- `idTipoPerfil` - ID do tipo de perfil do usuário
- `idSeguradora` - ID da seguradora associada ao usuário
- `AdminERP` - Indica se o usuário é administrador do ERP
- `bloqueado` - Indica se o usuário está bloqueado
- `bloqueado_ate` - Data até quando o usuário permanecerá bloqueado
- `created_at` - Data de criação do registro
- `updated_at` - Data da última atualização do registro
- `ativo` - Indica se o usuário está ativo no sistema

### Tipo de Perfil
- `id_tipo_perfil` - Identificador único (auto-incremento)
- `perfil` - Nome do perfil
- `created_at` - Data de criação do registro
- `updated_at` - Data da última atualização do registro
- `ativo` - Indica se o perfil está ativo no sistema

### Seguradora
- `id_seguradora` - Identificador único (auto-incremento)
- `seguradora` - Nome da seguradora
- `nome_abreviado` - Nome abreviado da seguradora
- `codigo_susep` - Código SUSEP da seguradora
- `created_at` - Data de criação do registro
- `updated_at` - Data da última atualização do registro
- `ativo` - Indica se a seguradora está ativa no sistema

### Evento
- `idCodigoEvento` - Identificador único (auto-incremento)
- `Evento` - Código numérico do evento
- `Descricao` - Descrição do evento
- `idSeguradora` - ID da seguradora associada
- `created_at` - Data de criação do registro
- `updated_at` - Data da última atualização do registro
- `ativo` - Indica se o evento está ativo no sistema

### Objeto de Contabilização
- `idObjetoContabilizacao` - Identificador único (auto-incremento)
- `ObjetoContabilizacao` - Nome do objeto de contabilização
- `Descricao` - Descrição do objeto
- `idSeguradora` - ID da seguradora associada
- `created_at` - Data de criação do registro
- `updated_at` - Data da última atualização do registro
- `ativo` - Indica se o objeto está ativo no sistema

### Sistema Contábil
- `idSistemaContabil` - Identificador único (auto-incremento)
- `SistemaContabil` - Nome do sistema contábil
- `idSeguradora` - ID da seguradora associada
- `created_at` - Data de criação do registro
- `updated_at` - Data da última atualização do registro
- `ativo` - Indica se o sistema está ativo

## Endpoints da API

### Autenticação
- `POST /auth/login` - Realiza login e retorna tokens JWT
- `POST /auth/refresh` - Renova tokens JWT
- `GET /csrf/token` - Obtém um token CSRF

### Usuários (Requer Autenticação)
- `GET /usuarios` - Lista todos os usuários
- `GET /usuarios/{id}` - Busca um usuário pelo ID
- `POST /usuarios` - Cria um novo usuário
- `PUT /usuarios/{id}` - Atualiza um usuário existente
- `DELETE /usuarios/{id}` - Remove um usuário (desativa)

### Tipos de Perfil (Requer Autenticação)
- `GET /tipos-perfil` - Lista todos os tipos de perfil
- `GET /tipos-perfil/{id}` - Busca um tipo de perfil pelo ID
- `POST /tipos-perfil` - Cria um novo tipo de perfil
- `PUT /tipos-perfil/{id}` - Atualiza um tipo de perfil existente
- `DELETE /tipos-perfil/{id}` - Remove um tipo de perfil (desativa)

### Seguradoras (Requer Autenticação)
- `GET /seguradoras` - Lista todas as seguradoras
- `GET /seguradoras/{id}` - Busca uma seguradora pelo ID
- `POST /seguradoras` - Cria uma nova seguradora
- `PUT /seguradoras/{id}` - Atualiza uma seguradora existente
- `DELETE /seguradoras/{id}` - Remove uma seguradora (desativa)

### Eventos (Requer Autenticação)
- `GET /eventos` - Lista todos os eventos
- `GET /eventos/{id}` - Busca um evento pelo ID
- `GET /eventos/seguradora/{id}` - Lista eventos de uma seguradora
- `POST /eventos` - Cria um novo evento
- `PUT /eventos/{id}` - Atualiza um evento existente
- `DELETE /eventos/{id}` - Remove um evento (desativa)

### Objetos de Contabilização (Requer Autenticação)
- `GET /objetos-contabilizacao` - Lista todos os objetos
- `GET /objetos-contabilizacao/{id}` - Busca um objeto pelo ID
- `GET /objetos-contabilizacao/seguradora/{id}` - Lista objetos de uma seguradora
- `POST /objetos-contabilizacao` - Cria um novo objeto
- `PUT /objetos-contabilizacao/{id}` - Atualiza um objeto existente
- `DELETE /objetos-contabilizacao/{id}` - Remove um objeto (desativa)

### Sistemas Contábeis (Requer Autenticação)
- `GET /sistemas-contabeis` - Lista todos os sistemas
- `GET /sistemas-contabeis/{id}` - Busca um sistema pelo ID
- `GET /sistemas-contabeis/seguradora/{id}` - Lista sistemas de uma seguradora
- `POST /sistemas-contabeis` - Cria um novo sistema
- `PUT /sistemas-contabeis/{id}` - Atualiza um sistema existente
- `DELETE /sistemas-contabeis/{id}` - Remove um sistema (desativa)

### Configurações de Sistema Contábil (Requer Autenticação)
- `GET /sistemas-contabeis-config` - Lista todas as configurações
- `GET /sistemas-contabeis-config/{id}` - Busca uma configuração pelo ID
- `GET /sistemas-contabeis-config/seguradora/{id}` - Lista configurações de uma seguradora
- `GET /sistemas-contabeis-config/sistema/{id}` - Lista configurações de um sistema
- `POST /sistemas-contabeis-config` - Cria uma nova configuração
- `PUT /sistemas-contabeis-config/{id}` - Atualiza uma configuração existente
- `DELETE /sistemas-contabeis-config/{id}` - Remove uma configuração (desativa)

## Exemplos de Uso

### Login

\`\`\`bash
curl -X POST http://localhost:8080/auth/login \
 -H "Content-Type: application/json" \
 -d '{
   "login": "admin",
   "senha": "Admin@123"
 }'
\`\`\`

### Obter Token CSRF

\`\`\`bash
curl -X GET http://localhost:8080/csrf/token \
 -H "Authorization: Bearer seu_token_jwt"
\`\`\`

### Criar um evento (com autenticação e CSRF)

\`\`\`bash
curl -X POST http://localhost:8080/eventos \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer seu_token_jwt" \
 -H "X-CSRF-Token: seu_token_csrf" \
 -d '{
   "evento": 1001,
   "descricao": "Pagamento de Sinistro",
   "idSeguradora": 1,
   "ativo": true
 }'
\`\`\`

## Contribuição

Para contribuir com o projeto:

1. Faça um fork do repositório
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Faça commit das suas alterações (`git commit -m 'Adiciona nova feature'`)
4. Faça push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo LICENSE para detalhes.
