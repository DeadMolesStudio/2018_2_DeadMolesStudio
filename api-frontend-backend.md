# API Frontend<->Backend (REST API)

1. Сессия

    **Запрос: `/session`**

    1. **GET**: получить сессию пользователя

        Если есть сессия, то она в куке session_id

        ```http
        GET /session HTTP/1.1
        Cookie: session_id=k_-5sqLMSj2oIO_EsBui180GQMCPVWnj1Wcdu-hMngw=
        ```

        **Ответ:**

        1. Пользователь залогинен:

            ```http
            HTTP/1.1 200 OK
            ```
            ```json
            {
                "session_id": "k_-5sqLMSj2oIO_EsBui180GQMCPVWnj1Wcdu-hMngw="
            }
            ```

        2. Не залогинен:

            ```http
            HTTP/1.1 401 Unauthorized
            ```

        3. Если ошибка в бд:

            ```http
            HTTP/1.1 500 Internal Server Error
            ```

    2. **POST**: создать сессию (залогинить пользователя)

        ```http
        POST /login HTTP/1.1
        ```
        ```json
        {
            "email": "email@email.com",
            "password": "password"
        }
        ```

        **Ответ:**

        1. Если неверный формат JSON, то:

            ```http
            HTTP/1.1 400 Bad Request
            ```

        2. Неверная пара пользователь/пароль:

            ```http
            HTTP/1.1 422 Unprocessable Entity
            ```

        3. Успешный вход:

            ```http
            HTTP/1.1 200 OK
            Set-Cookie: session_id=k_-5sqLMSj2oIO_EsBui180GQMCPVWnj1Wcdu-hMngw=
            ```

        4. Если пользователь уже залогинен:

            ```http
            HTTP/1.1 200 OK
            ```

        5. Если ошибка в бд:

            ```http
            HTTP/1.1 500 Internal Server Error
            ```

    3. **DELETE**: разлогинить пользователя

        ```http
        DELETE /session HTTP/1.1
        ```

        **Ответ:**

        1. Вчерашняя кука, если залогинен:

            ```http
            HTTP/1.1 200 OK
            Set-Cookie: session_id=k_-5sqLMSj2oIO_EsBui180GQMCPVWnj1Wcdu-hMngw=; Expires=Thu, 20 Sep 2018 09:40:57 GMT
            ```

        2. Если уже разлогинен:

            ```http
            HTTP/1.1 200 OK
            ```

2. Профиль

    **Запрос:** `/profile`

    1. **GET**: получить профиль

        По ID

        ```http
        GET /profile?id=32 HTTP/1.1
        ```

        По никнейму

        ```http
        GET /profile?nickname=Nick HTTP/1.1
        ```

        Без ID (из сессии)

        ```http
        GET /profile HTTP/1.1
        ```

        **Ответ:**

        1. Пользователь с этим ID, этой сессией, никнеймом существует:

            ```http
            HTTP/1.1 200 OK
            ```
            ```json
            {
                "id": 32,
                "nickname": "Nick",
                "email": "email@email.com",
                "record": 100500,
                "win": 21,
                "draws": 2,
                "loss": 15
            }
            ```

        2. Не найдено:

            ```http
            HTTP/1.1 404 Not Found
            ```

        3. Неправильный запрос

            ```http
            HTTP/1.1 400 Bad Request
            ```

        4. Если ошибка в бд:

            ```http
            HTTP/1.1 500 Internal Server Error
            ```

    2. **POST**: создать пользователя

        Все параметры при регистрация:

        ```http
        POST /profile HTTP/1.1
        ```
        ```json
        {
            "nickname": "Nick",
            "email": "email@email.com",
            "password": "mysecretpassword"
        }
        ```

        **Ответ:**

        Если неверный формат JSON, то:

        ```http
        HTTP/1.1 400 Bad Request
        ```

        Успешная регистрация и логин:

        ```http
        HTTP/1.1 200 OK
        Set-Cookie: session_id=k_-5sqLMSj2oIO_EsBui180GQMCPVWnj1Wcdu-hMngw=
        ```

        Занята почта или ник, пароль не удовлетворяет правилам безопасности, другие ошибки:

        ```http
        HTTP/1.1 403 Forbidden
        ```
        ```json
        {
            "error": [
                {
                    "field": "nickname",
                    "text": "Этот никнейм уже занят",
                },
                {
                    "field": "email",
                    "text": "Данная почта уже занята",
                },
                {
                    "field": "password",
                    "text": "Пароль должен быть не менее 8 символов",
                },
            ]
        }
        ```

        Возможные ошибки:
        * Никнейм:
            * Никнейм занят: "Этот никнейм уже занят";
            * Никнейм короткий (<4): "Никнейм должен быть не менее 4 символов";
            * Никнейм длинный (>32): "Никнейм должен быть не более 32 символов";
        * Почта:
            * Почта занята: "Данная почта уже занята"
            * Формат почты неверный: "Неверный формат почты"
        * Пароль:
            * Пароль короткий (<8): "Пароль должен быть не менее 8 символов"
            * Пароль короткий (>32): "Пароль должен быть не более 32 символов"

        При регистрации не все параметры:

        ```http
        HTTP/1.1 422 Unprocessable Entity
        ```

        Если ошибка в бд:

        ```http
        HTTP/1.1 500 Internal Server Error
        ```

    3. **PUT**: изменить пользователя

        Часть параметров или все:

        ```http
        PUT /profile HTTP/1.1
        ```
        ```json
        {
            "nickname": "NickName",
            "email": "newemail@email.com",
            "password": "mynewsecretpassword"
        }
        ```

        **Ответ:**

        Если неверный формат JSON, то:

        ```http
        HTTP/1.1 400 Bad Request
        ```

        Успешное изменение:

        ```http
        HTTP/1.1 200 OK
        ```

        Не залогинен:

        ```http
        HTTP/1.1 401 Unauthorized
        ```

        Если ошибка в бд:

        ```http
        HTTP/1.1 500 Internal Server Error
        ```

        Занята почта или ник, пароль не удовлетворяет правилам безопасности, другие ошибки:

        См. **POST** одноименный пункт

3. Скорборд

    **Запрос `/scoreboard`**

    **GET**: получить табличку или ее часть (limit, offset), в примере, 10 игроков, начиная с 11 места (11-20)

    ```http
    GET /scoreboard?limit=10&page=1 HTTP/1.1
    ```

    **Ответ:**

    1. Успех:

        ```http
        HTTP/1.1 200 OK
        ```
        ```json
        {
            "players": [
                {
                    "id": 143,
                    "nickname": "Nick",
                    "record": 100500
                },
                // 8 more...
                {
                    "id": 34,
                    "nickname": "LuckyBoy",
                    "record": 100000,
                },
            ],
            "total": 42
        }
        ```

        или
        ```http
        HTTP/1.1 200 OK
        ```
        ```json
        {
            "players": [],
            "total": 42
        }
        ```
        если пользователей нет.

    2. Если ошибка в бд:

        ```http
        HTTP/1.1 500 Internal Server Error
        ```

P.S. Ни один из методов (**GET**, **POST**, **PUT**, **DELETE**)

    ```http
    HTTP/1.1 405 Method Not Allowed
    ```
