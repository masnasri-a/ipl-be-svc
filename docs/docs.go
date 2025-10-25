package docs

import "github.com/swaggo/swag"

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "IPL Backend Service API",
	Description:      "RESTful API for IPL Backend Service with menu management",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

// docTemplate holds the base swagger template
var docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/health": {
            "get": {
                "description": "Check if the service is running",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "Service is running",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/menus/user/{user_id}": {
            "get": {
                "description": "Get list of menus accessible by a specific user ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "menus"
                ],
                "summary": "Get menus by user ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Menus retrieved successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/handler.MenuResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid user ID",
                        "schema": {
                            "$ref": "#/definitions/utils.APIResponse"
                        }
                    },
                    "404": {
                        "description": "No menus found",
                        "schema": {
                            "$ref": "#/definitions/utils.APIResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.APIResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.MenuResponse": {
            "type": "object",
            "properties": {
                "document_id": {
                    "type": "string",
                    "example": "mo5qqs8ezbruui07t91p6da8"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "is_active": {
                    "type": "boolean",
                    "example": true
                },
                "kode_menu": {
                    "type": "string",
                    "example": "master-data"
                },
                "nama_menu": {
                    "type": "string",
                    "example": "Master Data"
                },
                "published_at": {
                    "type": "string",
                    "example": "2025-10-23T15:16:28.206Z"
                },
                "urutan_menu": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "utils.APIResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        }
    }
}`