// Package docs contém a documentação Swagger para a API.
package docs

import (
	"github.com/swaggo/swag"
)

// SwaggerInfo contém as informações básicas da API.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "API de Gerenciamento de Seguradoras",
	Description:      "API para gerenciamento de usuários, tipos de perfil, seguradoras, eventos, objetos de contabilização e sistemas contábeis com autenticação JWT, limite de tentativas de login, rotação de tokens, auditoria e diversas medidas de segurança",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

// SwaggerInfo_swagger contém as informações básicas da API.
var SwaggerInfo_swagger = SwaggerInfo

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

// O template Swagger básico
const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Autenticação usando JWT. Exemplo: \"Bearer {token}\""
        }
    },
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Autentica um usuário e retorna um token JWT. Após 5 tentativas falhas em 15 minutos, a conta será bloqueada por 30 minutos.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "autenticação"
                ],
                "summary": "Login de usuário",
                "parameters": [
                    {
                        "description": "Credenciais de login",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "login": {
                                    "type": "string"
                                },
                                "senha": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "access_token": {
                                    "type": "string"
                                },
                                "refresh_token": {
                                    "type": "string"
                                },
                                "usuario": {
                                    "$ref": "#/definitions/models.Usuario"
                                },
                                "expires_in": {
                                    "type": "integer"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "429": {
                        "description": "Too Many Requests",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "Renova um token JWT válido",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "autenticação"
                ],
                "summary": "Renovar token",
                "parameters": [
                    {
                        "description": "Refresh token",
                        "name": "refresh",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "refresh_token": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "access_token": {
                                    "type": "string"
                                },
                                "refresh_token": {
                                    "type": "string"
                                },
                                "expires_in": {
                                    "type": "integer"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/eventos": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retorna todos os eventos cadastrados",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "eventos"
                ],
                "summary": "Listar todos os eventos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Evento"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cria um novo evento",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "eventos"
                ],
                "summary": "Criar evento",
                "parameters": [
                    {
                        "description": "Dados do evento",
                        "name": "evento",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Evento"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Evento"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/eventos/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retorna um evento específico pelo ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "eventos"
                ],
                "summary": "Buscar evento por ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do evento",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Evento"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Atualiza um evento existente",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "eventos"
                ],
                "summary": "Atualizar evento",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do evento",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Dados do evento",
                        "name": "evento",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Evento"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Evento"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Remove (desativa) um evento",
                "tags": [
                    "eventos"
                ],
                "summary": "Remover evento",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do evento",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/objetos-contabilizacao": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retorna todos os objetos de contabilização cadastrados",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "objetos-contabilizacao"
                ],
                "summary": "Listar todos os objetos de contabilização",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ObjetoContabilizacao"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cria um novo objeto de contabilização",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "objetos-contabilizacao"
                ],
                "summary": "Criar objeto de contabilização",
                "parameters": [
                    {
                        "description": "Dados do objeto de contabilização",
                        "name": "objeto",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ObjetoContabilizacao"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.ObjetoContabilizacao"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/sistemas-contabeis": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retorna todos os sistemas contábeis cadastrados",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sistemas-contabeis"
                ],
                "summary": "Listar todos os sistemas contábeis",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.SistemaContabil"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cria um novo sistema contábil",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sistemas-contabeis"
                ],
                "summary": "Criar sistema contábil",
                "parameters": [
                    {
                        "description": "Dados do sistema contábil",
                        "name": "sistema",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SistemaContabil"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SistemaContabil"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/sistemas-contabeis-config": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retorna todas as configurações de sistema contábil cadastradas",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sistemas-contabeis-config"
                ],
                "summary": "Listar todas as configurações de sistema contábil",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.SistemaContabilConfig"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cria uma nova configuração de sistema contábil",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sistemas-contabeis-config"
                ],
                "summary": "Criar configuração de sistema contábil",
                "parameters": [
                    {
                        "description": "Dados da configuração",
                        "name": "config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SistemaContabilConfig"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SistemaContabilConfig"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/objetos-contabilizacao-eventos": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retorna todas as relações entre objetos de contabilização e eventos",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "objetos-contabilizacao-eventos"
                ],
                "summary": "Listar todas as relações entre objetos de contabilização e eventos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ObjetoContabilizacaoEvento"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Cria uma nova relação entre objeto de contabilização e evento",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "objetos-contabilizacao-eventos"
                ],
                "summary": "Criar relação entre objeto de contabilização e evento",
                "parameters": [
                    {
                        "description": "Dados da relação",
                        "name": "relacao",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ObjetoContabilizacaoEvento"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.ObjetoContabilizacaoEvento"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Evento": {
            "type": "object",
            "properties": {
                "ativo": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "descricao": {
                    "type": "string"
                },
                "evento": {
                    "type": "integer"
                },
                "idCodigoEvento": {
                    "type": "integer"
                },
                "idSeguradora": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        },
        "models.ObjetoContabilizacao": {
            "type": "object",
            "properties": {
                "ativo": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "descricao": {
                    "type": "string"
                },
                "idObjetoContabilizacao": {
                    "type": "integer"
                },
                "idSeguradora": {
                    "type": "integer"
                },
                "objetoContabilizacao": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        },
        "models.ObjetoContabilizacaoEvento": {
            "type": "object",
            "properties": {
                "ativo": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "eventoDescricao": {
                    "type": "string"
                },
                "eventoNumero": {
                    "type": "integer"
                },
                "idCodigoEvento": {
                    "type": "integer"
                },
                "idObjetoContabilizacao": {
                    "type": "integer"
                },
                "idObjetoContabilizacaoEvento": {
                    "type": "integer"
                },
                "idSeguradora": {
                    "type": "integer"
                },
                "objetoContabilizacaoNome": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        },
        "models.Seguradora": {
            "type": "object",
            "properties": {
                "ativo": {
                    "type": "boolean"
                },
                "codigo_susep": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "id": {
                    "type": "integer"
                },
                "nome": {
                    "type": "string"
                },
                "nome_abreviado": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        },
        "models.SistemaContabil": {
            "type": "object",
            "properties": {
                "ativo": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "idSeguradora": {
                    "type": "integer"
                },
                "idSistemaContabil": {
                    "type": "integer"
                },
                "sistemaContabil": {
                    "type": "  "integer"
                },
                "sistemaContabil": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        },
        "models.SistemaContabilConfig": {
            "type": "object",
            "properties": {
                "ativo": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "eventoDescricao": {
                    "type": "string"
                },
                "eventoNumero": {
                    "type": "integer"
                },
                "idCodigoEvento": {
                    "type": "integer"
                },
                "idObjetoContabilizacao": {
                    "type": "integer"
                },
                "idSeguradora": {
                    "type": "integer"
                },
                "idSistemaContabil": {
                    "type": "integer"
                },
                "idSistemaContabilConfig": {
                    "type": "integer"
                },
                "objetoContabilizacaoNome": {
                    "type": "string"
                },
                "sistemaContabilNome": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        },
        "models.TipoPerfil": {
            "type": "object",
            "properties": {
                "ativo": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "id": {
                    "type": "integer"
                },
                "perfil": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        },
        "models.Usuario": {
            "type": "object",
            "properties": {
                "adminERP": {
                    "type": "boolean"
                },
                "ativo": {
                    "type": "boolean"
                },
                "bloqueado": {
                    "type": "boolean"
                },
                "bloqueado_ate": {
                    "type": "string",
                    "format": "date-time"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "idSeguradora": {
                    "type": "integer"
                },
                "idTipoPerfil": {
                    "type": "integer"
                },
                "login": {
                    "type": "string"
                },
                "nome": {
                    "type": "string"
                },
                "senha": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        }
    }
}`
