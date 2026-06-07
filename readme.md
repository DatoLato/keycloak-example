# Установка synapse и kk.

### Для установки synapse рекомендуется использовать редми для synapse [synapse_redme.md]


### Установка kk.

Есть вот такой замечательный репозиторий: https://github.com/s-rb/keycloak-dockerized-ssl-nginx

* Берем смело и клонируем.

* cd keycloak-dockerized-ssl-nginx

* ткройте файл .env и отредактируйте следующие переменные:

```
KEYCLOAK_ADMIN=admin
KEYCLOAK_ADMIN_PASSWORD=password
PROXY_ADDRESS_FORWARDING=true
KC_PROXY=edge

KC_DB=postgres
KC_DB_URL=jdbc:postgresql://keycloak-postgres:5432/keycloak
KC_DB_USERNAME=keycloak
KC_DB_PASSWORD=password
POSTGRES_DB=keycloak
POSTGRES_USER=keycloak
POSTGRES_PASSWORD=password
```

* Данный пункт может выполняться и до первого шага - он никак не зависит от него. Далее в инструкции полагаем что у вас будет зарегистрирован свой домен (напримерsurkoff.com) и мы хотим чтобы Keycloak был бы доступен по my-keycloak.surkoff.com


* В конфигах nginx - default.conf_with_ssl, default.conf_without_ssl указываем свой домен:

Пример: 

```
server {
    listen 80;
    server_name my-keycloak.surkoff.com;

    location /.well-known/acme-challenge/ {
        root /data/letsencrypt;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name my-keycloak.surkoff.com;

    ssl_certificate /etc/letsencrypt/live/my-keycloak.surkoff.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/my-keycloak.surkoff.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://keycloak:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_set_header X-Forwarded-For $host;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
    }
}

```

* cp nginx/conf.d/default.conf_without_ssl nginx/conf.d/default.conf docker-compose up -d
* docker exec certbot certbot certonly --webroot --webroot-path=/data/letsencrypt -d my-keycloak.surkoff.com --email your_email@gmail.com --agree-tos --no-eff-email


* docker-compose down / cp nginx/conf.d/default.conf_with_ssl nginx/conf.d/default.conf / docker-compose up -d
