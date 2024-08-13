# Устанавливаем аргумент для версии Go с значением по умолчанию
ARG GO_VERSION=1.22

# Используем аргумент для установки версии Go
FROM golang:${GO_VERSION}-alpine

# Устанавливаем необходимые зависимости
RUN apk add git make gcc libc-dev && rm -rf /var/cache/apk/*

# Добавляем переменную окружения для GitLab credentials
ENV GIT_CREDENTIALS=""
ENV GIT_DOMAIN=""

# Указываем рабочую директорию
ENV WORKDIR=/src
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

# Создаем entrypoint скрипт
RUN printf '#!/bin/sh\n\
set -e\n\
\n\
if [ "$1" = "with_creds" ]; then\n\
    # Удаляем первый аргумент, чтобы основная команда была корректно выполнена\n\
    shift\n\
    echo "Set credentials for ${GIT_DOMAIN}"\n\
    git config --global url."https://${GIT_CREDENTIALS}@${GIT_DOMAIN}".insteadOf "https://${GIT_DOMAIN}"\n\
    echo "Executing command: $@"\n\
    "$@"\n\
    echo "Unset credentials"\n\
    git config --global --unset url."https://${GIT_CREDENTIALS}@${GIT_DOMAIN}".insteadOf\n\
else\n\
    exec "$@"\n\
fi\n' > /entrypoint.sh

# Делаем скрипт исполняемым
RUN chmod +x /entrypoint.sh

# Устанавливаем скрипт в качестве entrypoint
ENTRYPOINT ["/entrypoint.sh"]
