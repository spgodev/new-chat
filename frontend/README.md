# New Chat Frontend v2

Готовый статический frontend под текущий Go backend.

## Login
`username + password`

## Register
`username + email + password`

## WebSocket
`ws://localhost:8080/api/ws?token=<access_token>`

## Запуск
```powershell
cd new_chat_frontend_v2
python -m http.server 5500
```

Открыть:
http://localhost:5500

Backend:
http://localhost:8080

Если браузер покажет CORS error, нужно разрешить `http://localhost:5500` в Gin backend.
