// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2018-11-05 23:36:53.624856 +0300 MSK m=+0.192680864

package docs

import (
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "This is a backend server for the game.",
        "title": "The Ketnipz Game API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Artyom Andreev",
            "email": "aandreev06.1998@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "basePath": "/api",
    "paths": {
        "/profile": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Получить профиль пользователя по ID, email или из сессии",
                "operationId": "get-profile",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Никнейм",
                        "name": "nickname",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь найден, успешно",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.Profile"
                        }
                    },
                    "400": {
                        "description": "Неправильный запрос"
                    },
                    "401": {
                        "description": "Не залогинен"
                    },
                    "404": {
                        "description": "Не найдено"
                    },
                    "500": {
                        "description": "Ошибка в бд"
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Изменить профиль, должен быть залогинен",
                "operationId": "put-profile",
                "parameters": [
                    {
                        "description": "Новые никнейм, и/или почта, и/или пароль",
                        "name": "Profile",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.RegisterProfile"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь найден, успешно изменены данные"
                    },
                    "400": {
                        "description": "Неверный формат JSON"
                    },
                    "401": {
                        "description": "Не залогинен"
                    },
                    "403": {
                        "description": "Занята почта или ник, пароль не удовлетворяет правилам безопасности, другие ошибки",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ProfileErrorList"
                        }
                    },
                    "500": {
                        "description": "Ошибка в бд"
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Зарегистрировать по никнейму, почте и паролю и автоматически залогинить",
                "operationId": "post-profile",
                "parameters": [
                    {
                        "description": "Никнейм, почта и пароль",
                        "name": "Profile",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.RegisterProfile"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь зарегистрирован и залогинен успешно"
                    },
                    "400": {
                        "description": "Неверный формат JSON"
                    },
                    "403": {
                        "description": "Занята почта или ник, пароль не удовлетворяет правилам безопасности, другие ошибки",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ProfileErrorList"
                        }
                    },
                    "422": {
                        "description": "При регистрации не все параметры"
                    },
                    "500": {
                        "description": "Ошибка в бд"
                    }
                }
            }
        },
        "/scoreboard": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Получить таблицу лидеров (пагинация присутствует)",
                "operationId": "get-scoreboard",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Пользователей на страницу",
                        "name": "Limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Страница номер",
                        "name": "Page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Таблицу лидеров или ее страница и общее количество",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.PositionList"
                        }
                    },
                    "500": {
                        "description": "Ошибка в бд"
                    }
                }
            }
        },
        "/session": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Получить сессию пользователя, если есть сессия, то она в куке session_id",
                "operationId": "get-session",
                "responses": {
                    "200": {
                        "description": "Пользователь залогинен, успешно",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.Session"
                        }
                    },
                    "401": {
                        "description": "Не залогинен"
                    },
                    "500": {
                        "description": "Ошибка в бд"
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Залогинить пользователя (создать сессию)",
                "operationId": "post-session",
                "parameters": [
                    {
                        "description": "Почта и пароль",
                        "name": "UserPassword",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.UserPassword"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный вход / пользователь уже залогинен",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.Session"
                        }
                    },
                    "400": {
                        "description": "Неверный формат JSON"
                    },
                    "422": {
                        "description": "Неверная пара пользователь/пароль"
                    },
                    "500": {
                        "description": "Внутренняя ошибка"
                    }
                }
            },
            "delete": {
                "summary": "Разлогинить пользователя (удалить сессию)",
                "operationId": "delete-session",
                "responses": {
                    "200": {
                        "description": "Успешный выход / пользователь уже разлогинен"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Position": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 42
                },
                "nickname": {
                    "type": "string",
                    "example": "Nick"
                },
                "record": {
                    "type": "integer",
                    "example": 100500
                }
            }
        },
        "models.PositionList": {
            "type": "object",
            "properties": {
                "players": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Position"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "models.Profile": {
            "type": "object",
            "properties": {
                "draws": {
                    "type": "integer"
                },
                "email": {
                    "type": "string",
                    "example": "email@email.com"
                },
                "id": {
                    "type": "integer"
                },
                "loss": {
                    "type": "integer"
                },
                "nickname": {
                    "type": "string",
                    "example": "Nick"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                },
                "record": {
                    "type": "integer"
                },
                "win": {
                    "type": "integer"
                }
            }
        },
        "models.ProfileError": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string",
                    "example": "nickname"
                },
                "text": {
                    "type": "string",
                    "example": "Этот никнейм уже занят"
                }
            }
        },
        "models.ProfileErrorList": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ProfileError"
                    }
                }
            }
        },
        "models.RegisterProfile": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "email@email.com"
                },
                "nickname": {
                    "type": "string",
                    "example": "Nick"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                }
            }
        },
        "models.Session": {
            "type": "object",
            "properties": {
                "session_id": {
                    "type": "string",
                    "example": "ef84d238-47ef-4452-9536-99380db79911"
                }
            }
        },
        "models.UserPassword": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "email@email.com"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                }
            }
        }
    }
}`

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}
func init() {
	swag.Register(swag.Name, &s{})
}
