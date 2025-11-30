Микросервис для обработки заказов: получение из Kafka, хранение в PostgreSQL и быстрый доступ через LRU-кеш. В комплекте есть пример фронтенда на HTML+JS.
## Возможности
 - Чтение заказов из Kafka
 - Сохраняет структуру заказа с платежом, доставкой и товарами
 - Быстрый in-memory LRU-кеш
 - HTTP API для получения заказа по ID
 - Простой веб-интерфейс

## Технологии
 - Golang
 - Kafka (segmentio/kafka-go), Kafka UI
 - PostgreSQL (sqlx)
 - Docker, Docker Compose
 - HTML+JS (frontend)

## Быстрый старт
  ### Клонирование
    git clone https://github.com/MustafaevAlim/level0.git
    cd level0

  ### Настрой .env
Создать .env на основе .env_example и заполнить значения:
    
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=yourpassword
    POSTGRES_DB=ordersdb
    KAFKA_BROKERS=localhost:9092
    KAFKA_TOPIC=orders-topic
    KAFKA_GROUP=orders-group
    CACHE_SIZE=100
    HTTP_PORT=8082

   ### Запуск через Docker Compose
    docker compose up --build

Чтобы управлять миграциями локально нужно скачать migrate:

    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
    | tar xvz -C /usr/local/bin

Используется makefile:

    make migrate-up # применить миграции
    make migrate-down # отменить миграции
    make migrate-version # посмотреть версию миграции

Это поднимет контейнеры с сервисом, Kafka, Kafka UI и PostgreSQL.

Создать топик в кафке (настроить по желанию):

    docker exec -it kafka /bin/bash # войти в контейнер с кафкой 
    kafka-topics --create --topic orders-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1


## HTTP API
Получить заказ:

    GET /order/{order_uid}
Пример:

    curl http://localhost:8082/order/OOBmrfkDRphyYFiH

## Статичная HTML-страница
Открыть /info/ или (при нужной конфигурации) /order.html в браузере, чтобы воспользоваться веб-интерфейсом поиска заказа.
## Структура проекта
    ├── cmd/
    │   └── myapp/          # main.go — точка входа
    ├── internal/
    │   ├── config/         # Настройки приложения
    │   ├── app/            # Жизненный цикл приложения
    │   ├── api/            # HTTP API и маршруты
    │   ├── repository/     # Работа с БД, Kafka, Cache
    │   ├── model/          # модели данных
    ├── web/                # фронтенд
    ├── migrations/         # миграции БД
    ├── scripts/            # скрипты(имитация записи сообщений в кафку)
    ├── vendor/             # зависимости
    ├── Dockerfile
    ├── docker-compose.yml
    ├── .env.example
    ├── Makefile
    └── README.md


## Оптимизации
Сбор профиля показывает такую ситуацию
```
(pprof) top20
Showing nodes accounting for 60ms, 100% of 60ms total
Showing top 20 nodes out of 51
      flat  flat%   sum%        cum   cum%
      10ms 16.67% 16.67%       20ms 33.33%  github.com/jmoiron/sqlx.compileNamedQuery
      10ms 16.67% 33.33%       10ms 16.67%  internal/runtime/syscall.Syscall6
      10ms 16.67% 50.00%       10ms 16.67%  runtime.(*gcBitsArena).tryAlloc (inline)
      10ms 16.67% 66.67%       30ms 50.00%  runtime.findRunnable
      10ms 16.67% 83.33%       10ms 16.67%  runtime.futex
      10ms 16.67%   100%       10ms 16.67%  runtime.netpoll
         0     0%   100%       10ms 16.67%  bufio.(*Writer).Flush
         0     0%   100%       20ms 33.33%  github.com/jmoiron/sqlx.(*Tx).NamedExecContext
         0     0%   100%       20ms 33.33%  github.com/jmoiron/sqlx.NamedExecContext
         0     0%   100%       20ms 33.33%  github.com/jmoiron/sqlx.bindNamedMapper
         0     0%   100%       20ms 33.33%  github.com/jmoiron/sqlx.bindStruct
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*Conn).do
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*Conn).doRequest
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*Conn).heartbeat
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*Conn).heartbeat.func1
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*Conn).writeOperation (inline)
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*Conn).writeRequest
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*ConsumerGroup).nextGeneration.(*Generation).heartbeatLoop.func6
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*Generation).Start.func1
         0     0%   100%       10ms 16.67%  github.com/segmentio/kafka-go.(*timeoutCoordinator).heartbeat
(pprof) 
```

Видно что `sqlx` занимает 33% времени. Видимо потому что все время компилится маппинг структур, можно заменить на `prepared statements`

```
pprof) top20
Showing nodes accounting for 40ms, 100% of 40ms total
Showing top 20 nodes out of 65
      flat  flat%   sum%        cum   cum%
      10ms 25.00% 25.00%       10ms 25.00%  internal/runtime/syscall.Syscall6
      10ms 25.00% 50.00%       10ms 25.00%  io.ReadAtLeast
      10ms 25.00% 75.00%       10ms 25.00%  runtime.(*mspan).writeHeapBitsSmall
      10ms 25.00%   100%       10ms 25.00%  runtime.futex
         0     0%   100%       10ms 25.00%  database/sql.(*DB).retry
         0     0%   100%       10ms 25.00%  database/sql.(*Stmt).ExecContext
         0     0%   100%       10ms 25.00%  database/sql.(*Stmt).ExecContext.func1
         0     0%   100%       10ms 25.00%  database/sql.ctxDriverStmtExec
         0     0%   100%       10ms 25.00%  database/sql.resultFromStatement
         0     0%   100%       10ms 25.00%  github.com/jmoiron/sqlx.(*NamedStmt).ExecContext
         0     0%   100%       10ms 25.00%  github.com/lib/pq.(*conn).send
         0     0%   100%       10ms 25.00%  github.com/lib/pq.(*stmt).Exec
         0     0%   100%       10ms 25.00%  github.com/lib/pq.(*stmt).ExecContext
         0     0%   100%       10ms 25.00%  github.com/lib/pq.(*stmt).exec
         0     0%   100%       10ms 25.00%  github.com/segmentio/kafka-go.(*Conn).do
         0     0%   100%       10ms 25.00%  github.com/segmentio/kafka-go.(*Conn).offsetCommit
         0     0%   100%       10ms 25.00%  github.com/segmentio/kafka-go.(*Conn).offsetCommit.func2
         0     0%   100%       10ms 25.00%  github.com/segmentio/kafka-go.(*Conn).offsetCommit.func2.1 (inline)
         0     0%   100%       10ms 25.00%  github.com/segmentio/kafka-go.(*Conn).writeOperation (inline)
         0     0%   100%       10ms 25.00%  github.com/segmentio/kafka-go.(*Generation).CommitOffsets
(pprof) 
```
Теперь `sqlx` пропал из топ 20. Пропал `runtime.findRunnable`

Теперь снимем профиль кучи.
```
(pprof) top20 
Showing nodes accounting for 7237.89kB, 100% of 7237.89kB total
Showing top 20 nodes out of 90
      flat  flat%   sum%        cum   cum%
    1539kB 21.26% 21.26%     1539kB 21.26%  runtime.allocm
 1040.08kB 14.37% 35.63%  1040.08kB 14.37%  github.com/segmentio/kafka-go.(*writeBatch).add
  544.67kB  7.53% 43.16%   544.67kB  7.53%  github.com/segmentio/kafka-go/protocol.newPage
  528.17kB  7.30% 50.46%   528.17kB  7.30%  regexp.(*bitState).reset
  513.50kB  7.09% 57.55%   513.50kB  7.09%  regexp/syntax.(*compiler).inst
  512.25kB  7.08% 64.63%   512.25kB  7.08%  l4.5/internal/repository.(*LRUcache).Push
  512.08kB  7.07% 71.70%   512.08kB  7.07%  github.com/segmentio/kafka-go/protocol.structDecodeFuncOf.func1.1
  512.05kB  7.07% 78.78%   512.05kB  7.07%  context.(*cancelCtx).Done
  512.05kB  7.07% 85.85%   512.05kB  7.07%  net.(*netFD).connect.func2
  512.03kB  7.07% 92.93%   512.03kB  7.07%  github.com/go-playground/validator.(*Validate).extractStructCache
  512.01kB  7.07%   100%   512.01kB  7.07%  github.com/lib/pq.(*conn).readStatementDescribeResponse
         0     0%   100%   512.01kB  7.07%  database/sql.(*DB).QueryContext
         0     0%   100%   512.01kB  7.07%  database/sql.(*DB).QueryContext.func1
         0     0%   100%   512.05kB  7.07%  database/sql.(*DB).connectionOpener
         0     0%   100%   512.01kB  7.07%  database/sql.(*DB).query
         0     0%   100%   512.01kB  7.07%  database/sql.(*DB).queryDC
         0     0%   100%   512.01kB  7.07%  database/sql.(*DB).queryDC.func1
         0     0%   100%   512.01kB  7.07%  database/sql.(*DB).retry
         0     0%   100%   512.01kB  7.07%  database/sql.ctxDriverQuery
         0     0%   100%   512.01kB  7.07%  database/sql.withLock
```

Здесь все хорошо.

Посмотрим аллокации.
```
l4.5 git:(main) ✗ go tool pprof http://localhost:8082/debug/pprof/heap                                                                  
Fetching profile over HTTP from http://localhost:8082/debug/pprof/heap
Saved profile in /home/traktor/pprof/pprof.main.alloc_objects.alloc_space.inuse_objects.inuse_space.002.pb.gz
File: main
Build ID: 3746e69505901f5b404be73c79b6350f7d0cb2f6
Type: inuse_space
Time: 2025-11-30 23:33:24 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) sample_index = alloc_space 
(pprof) top
Showing nodes accounting for 9710.83kB, 75.97% of 12783.11kB total
Showing top 10 nodes out of 122
      flat  flat%   sum%        cum   cum%
 3120.25kB 24.41% 24.41%  3120.25kB 24.41%  github.com/segmentio/kafka-go.(*writeBatch).add
    2052kB 16.05% 40.46%     2052kB 16.05%  runtime.allocm
  902.59kB  7.06% 47.52%  1415.15kB 11.07%  compress/flate.NewWriter
  544.67kB  4.26% 51.78%   544.67kB  4.26%  github.com/segmentio/kafka-go/protocol.newPage
  528.17kB  4.13% 55.91%   528.17kB  4.13%  regexp.(*bitState).reset
  513.50kB  4.02% 59.93%   513.50kB  4.02%  regexp/syntax.(*compiler).inst
  512.62kB  4.01% 63.94%   512.62kB  4.01%  github.com/lib/pq.(*writeBuf).string
  512.56kB  4.01% 67.95%   512.56kB  4.01%  compress/flate.newHuffmanEncoder (inline)
  512.25kB  4.01% 71.96%   512.25kB  4.01%  l4.5/internal/repository.(*LRUcache).Push
  512.22kB  4.01% 75.97%   512.22kB  4.01%  database/sql.driverArgsConnLocked
(pprof) sample_index = alloc_objects
(pprof) top                         
Showing nodes accounting for 62149, 98.73% of 62950 total
Dropped 35 nodes (cum <= 314)
Showing top 10 nodes out of 87
      flat  flat%   sum%        cum   cum%
     32768 52.05% 52.05%      32768 52.05%  github.com/lib/pq.(*conn).readStatementDescribeResponse
      8192 13.01% 65.07%       8192 13.01%  github.com/go-playground/validator.(*Validate).extractStructCache
      5461  8.68% 73.74%       5461  8.68%  net.(*netFD).connect.func2
      4681  7.44% 81.18%       4681  7.44%  context.(*cancelCtx).Done
      4096  6.51% 87.69%       4096  6.51%  regexp.mergeRuneSets.func2
      3277  5.21% 92.89%       3277  5.21%  github.com/segmentio/kafka-go/protocol.structDecodeFuncOf.func1.1
      1170  1.86% 94.75%       1170  1.86%  database/sql.driverArgsConnLocked
      1025  1.63% 96.38%       1025  1.63%  runtime.allocm
      1024  1.63% 98.00%       1024  1.63%  l4.5/internal/repository.(*LRUcache).Push
       455  0.72% 98.73%        455  0.72%  compress/flate.newHuffmanEncoder
(pprof) 
```

Много аллокаций от драйвера PostgreSQL. Можно попробовать настроить количество соединений с бд.


```
Showing nodes accounting for 41448, 99.92% of 41480 total
Dropped 32 nodes (cum <= 207)
Showing top 10 nodes out of 48
      flat  flat%   sum%        cum   cum%
     21845 52.66% 52.66%      21845 52.66%  reflect.MakeSlice
      5461 13.17% 65.83%      27314 65.85%  github.com/segmentio/kafka-go.(*conn).run
      4681 11.28% 77.11%       4681 11.28%  context.(*cancelCtx).Done
      4681 11.28% 88.40%       4681 11.28%  regexp/syntax.(*parser).newRegexp (inline)
      2731  6.58% 94.98%       2731  6.58%  github.com/jmoiron/sqlx.scanAll
      1025  2.47% 97.45%       1025  2.47%  runtime.allocm
      1024  2.47% 99.92%       1024  2.47%  regexp.mergeRuneSets.func2 (inline)
         0     0% 99.92%       4681 11.28%  database/sql.(*DB).connectionOpener
         0     0% 99.92%       5705 13.75%  github.com/go-playground/validator.init
         0     0% 99.92%       2731  6.58%  github.com/jmoiron/sqlx.(*DB).SelectContext
(pprof) 
```

Аллокации от PostgreSQL пропали.