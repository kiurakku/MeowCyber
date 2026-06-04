# MeowCyber

Форк платформи AI-тестування безпеки (колишній CyberStrikeAI): український інтерфейс, бренд **MeowCyber**.

## Локальний запуск (Windows)

```powershell
cd D:\RE\CyberStrikeAI-1.6.30
$env:Path = "D:\go\mingw\mingw64\bin;" + $env:Path
$env:CGO_ENABLED = "1"
$env:GOMODCACHE = "D:\go\pkg\mod"
$env:GOCACHE = "D:\go\cache"
$env:TEMP = "D:\go\tmp"
go build -o meowcyber.exe cmd/server/main.go
.\meowcyber.exe -config config.yaml --https
```

Відкрийте **https://127.0.0.1:8088/** (або порт із `config.yaml`).

## Онлайн безкоштовно (GitHub + Render)

### 1. Репозиторій

Код уже налаштований під [github.com/kiurakku/MeowCyber](https://github.com/kiurakku/MeowCyber).

### 2. Render.com (безкоштовний план)

1. Увійдіть на [render.com](https://render.com) через GitHub.
2. **New → Blueprint** → підключіть репозиторій `kiurakku/MeowCyber`.
3. Render прочитає `render.yaml` і створить сервіс `meowcyber`.
4. У **Environment** додайте:
   - `OPENAI_API_KEY` — ваш ключ API;
   - `AUTH_PASSWORD` — надійний пароль (або згенерується автоматично);
   - `MEOWCYBER_HTTPS=0` — TLS надає Render, не самодіяльний сертифікат.
5. Після деплою URL буде на кшталт `https://meowcyber.onrender.com`.

**Обмеження free tier:** сервіс «засинає» після ~15 хв без запитів; перший запуск після сну — 30–60 с.

### 3. Власний домен (безкоштовно)

- У Render: **Settings → Custom Domains** → додайте домен.
- У DNS (Cloudflare, Namecheap тощо): CNAME на `meowcyber.onrender.com`.
- Безкоштовний домен: [Freenom](https://www.freenom.com) (обмежено), або піддомен Cloudflare / DuckDNS для тестів.

### 4. Альтернатива: Cloudflare Tunnel (з домашнього ПК)

```powershell
cloudflared tunnel --url http://localhost:8088
```

Отримаєте тимчасовий `*.trycloudflare.com` URL без відкриття портів.

## Мова інтерфейсу

За замовчуванням **українська** (`uk-UA`). Перемикач 🌐 у шапці: Українська / 中文 / English.

## Секрети

Не комітьте `config.yaml` з реальним `api_key`. Для хмари використовуйте змінні середовища `OPENAI_API_KEY` та `AUTH_PASSWORD`.
