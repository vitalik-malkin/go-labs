package echo_api

const (
	Swagger = `
  {
    "swagger": "2.0",
    "info": {
      "title": "echo-api.proto",
      "version": "version not set"
    },
    "consumes": [
      "application/json"
    ],
    "produces": [
      "application/json"
    ],
    "paths": {
      "/v1/echo/{text}": {
        "get": {
          "operationId": "EchoV1",
          "responses": {
            "200": {
              "description": "A successful response.",
              "schema": {
                "$ref": "#/definitions/echo_apiEchoV1Response"
              }
            },
            "default": {
              "description": "An unexpected error response",
              "schema": {
                "$ref": "#/definitions/runtimeError"
              }
            }
          },
          "parameters": [
            {
              "name": "text",
              "description": "Любой текст, который будет возвращен в ответе.\nДлина текста ограничивается 1024 символами.",
              "in": "path",
              "required": true,
              "type": "string"
            }
          ],
          "tags": [
            "EchoAPI"
          ]
        }
      }
    },
    "definitions": {
      "echo_apiEchoV1Response": {
        "type": "object",
        "properties": {
          "text": {
            "type": "string",
            "description": "Текст, переданный в эхо-запросе."
          }
        },
        "description": "Параметры ответа эхо-запроса."
      },
      "protobufAny": {
        "type": "object",
        "properties": {
          "type_url": {
            "type": "string"
          },
          "value": {
            "type": "string",
            "format": "byte"
          }
        }
      },
      "runtimeError": {
        "type": "object",
        "properties": {
          "error": {
            "type": "string"
          },
          "code": {
            "type": "integer",
            "format": "int32"
          },
          "message": {
            "type": "string"
          },
          "details": {
            "type": "array",
            "items": {
              "$ref": "#/definitions/protobufAny"
            }
          }
        }
      }
    }
  }
 `
)
