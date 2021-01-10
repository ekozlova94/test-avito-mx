# На основе какого образа мы создаём наш образ
FROM alpine:3.12.0

RUN mkdir /app
# Указываем рабочую директорию для контейнера
WORKDIR /app

# Копируем из текущей папки необходимые файлы в образ
COPY bin/test-avito-mx .
COPY bin/simulator .
COPY example-goods.xlsx .
COPY migrations ./migrations

# Делаем исполняемым наш бинарник
RUN chmod +x test-avito-mx
RUN chmod +x simulator