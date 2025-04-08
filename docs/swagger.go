package docs

import "github.com/swaggo/swag"

// @title SAP Segmentation API
// @version 1.0
// @description API для импорта и доступа к данным сегментации из SAP
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https
func SwaggerInfo() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
	})
}

// Template для Swagger документации.
// Будет заменен автоматически генерируемой документацией.
const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "description": "API для импорта и доступа к данным сегментации из SAP",
    "title": "SAP Segmentation API",
    "termsOfService": "http://swagger.io/terms/",
    "contact": {
      "name": "API Support",
      "email": "support@example.com"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "paths": {
    "/api/health": {
      "get": {
        "description": "Проверяет работоспособность API сервера",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "system"
        ],
        "summary": "Проверка работоспособности",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "/api/segmentation": {
      "get": {
        "description": "Возвращает список всех сегментов из базы данных",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "segmentation"
        ],
        "summary": "Получить все сегменты",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/model.Segmentation"
              }
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "/api/segmentation/import": {
      "post": {
        "description": "Запускает процесс импорта данных из SAP API в базу данных",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "segmentation"
        ],
        "summary": "Импортировать сегментацию",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "interface{}"
              }
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "/api/segmentation/{id}": {
      "get": {
        "description": "Возвращает сегмент с указанным SAP ID",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "segmentation"
        ],
        "summary": "Получить сегмент по ID",
        "parameters": [
          {
            "type": "string",
            "description": "SAP ID сегмента",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/model.Segmentation"
            }
          },
          "404": {
            "description": "Not Found",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "model.Segmentation": {
      "type": "object",
      "properties": {
        "address_sap_id": {
          "type": "string"
        },
        "adr_segment": {
          "type": "string"
        },
        "segment_id": {
          "type": "integer",
          "format": "int64"
        }
      }
    }
  }
}`
