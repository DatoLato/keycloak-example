# Как Synapse может использовать refresh_token?

По умолчанию Synapse не сохраняет refresh_token и не использует его автоматически.
Но есть обходные пути:

### Вариант 1: Кастомный OIDC-провайдер (через модуль Python)

Можно написать модуль, который:
*  Перехватывает refresh_token из Keycloak.
*  Обновляет токены при истечении access_token.

```python
from synapse.module_api import ModuleApi
from authlib.oauth2.rfc6749 import RefreshTokenGrant

class CustomOIDCProvider:
    def __init__(self, config: dict, api: ModuleApi):
        self.api = api
        self.refresh_tokens = {}  # Храним refresh_token для каждого пользователя

    async def on_oidc_callback(self, auth_result: dict):
        user_id = auth_result.get("sub")
        refresh_token = auth_result.get("refresh_token")
        if refresh_token:
            self.refresh_tokens[user_id] = refresh_token

    async def refresh_access_token(self, user_id: str):
        refresh_token = self.refresh_tokens.get(user_id)
        if not refresh_token:
            return None
        
        # Используем authlib или requests для обновления токена
        token_url = "https://keycloak.example.com/realms/myrealm/protocol/openid-connect/token"
        data = {
            "grant_type": "refresh_token",
            "refresh_token": refresh_token,
            "client_id": "synapse",
            "client_secret": "your-client-secret",
        }
        response = await self.api.http_client.post(token_url, data=data)
        return response.json()
```

Регистрация модуля в homeserver.yaml

```yaml
modules:
  - module: "custom_oidc.CustomOIDCProvider"
    config: {}
```

### Вариант 2: Использование expires_in + принудительный релогин

Если не нужен refresh_token, можно:

* Установить короткий Access Token Lifespan в Keycloak (например, 5 минут).
* Включить backchannel_logout в Synapse:

```yaml
oidc_providers:
  - ...
    backchannel_logout_enabled: true
```
