// Package generated Автоматически сгенерированная документация Swagger.
// Будет перезаписана при запуске swag init.
package generated

import "github.com/swaggo/swag"

func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
	})
}

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "description": "API для импорта и доступа к данным сегментации из SAP",
    "title": "SAP Segmentation API",
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
        "tags": ["system"],
        "summary": "Проверка работоспособности",
        "description": "Проверяет работоспособность API сервера",
        "produces": ["application/json"],
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
        "tags": ["segmentation"],
        "summary": "Получить все сегменты",
        "description": "Возвращает список всех сегментов из базы данных",
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/model.Segmentation"
              }
            }
          }
        }
      }
    },
    "/api/segmentation/{id}": {
      "get": {
        "tags": ["segmentation"],
        "summary": "Получить сегмент по ID",
        "description": "Возвращает сегмент с указанным SAP ID",
        "produces": ["application/json"],
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
          }
        }
      }
    },
    "/api/segmentation/import": {
      "post": {
        "tags": ["segmentation"],
        "summary": "Импортировать сегментацию",
        "description": "Запускает процесс импорта данных из SAP API в базу данных",
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "properties": {
                "count": {
                  "type": "integer"
                },
                "message": {
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
