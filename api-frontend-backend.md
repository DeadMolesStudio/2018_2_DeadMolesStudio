# API Frontend<->Backend (REST API)

1. Логин

    **Запрос:**

    ```http
    POST /login HTTP/1.1
    ```
    ```json
    {
        user: "username",
        password: "password"
    }
    ```

    **Ответ:**

    Если метод не ```POST```, то:

    ```http
    HTTP/1.1 405 Method Not Allowed
    ```

    Если неверный формат JSON, то:

    ```http
    HTTP/1.1 400 Bad Request
    ```

    Неверная пара пользователь/пароль:

    ```http
    HTTP/1.1 403 Forbidden
    ```

    Успешный вход:

    ```http
    HTTP/1.1 200 OK
    Set-Cookie: session_id=k_-5sqLMSj2oIO_EsBui180GQMCPVWnj1Wcdu-hMngw=
    ```

    Если пользователь уже залогинен:

    ```http
    HTTP/1.1 200 OK
    ```

2. Логаут

    **Запрос:**

    ```http
    GET /logout HTTP/1.1
    ```

    **Ответ:**

    Вчерашняя кука, если залогинен:

    ```http
    HTTP/1.1 200 OK
    Set-Cookie: session_id=k_-5sqLMSj2oIO_EsBui180GQMCPVWnj1Wcdu-hMngw=; Expires=Thu, 20 Sep 2018 09:40:57 GMT
    ```

    Если уже разлогинен:

    ```http
    HTTP/1.1 200 OK
    ```

3. ...
