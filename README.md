##
Сервис производит получение данных с API мосбиржи в заданном временном интервале


#### Общая идея
Наша цель – создать окружение вокруг одного основного сервиса с минимальной логической нагрузкой.
Простенький* сервис c полноценной обвязкой.


> \* c т.з бизнес логики.
Для дипломного проекта вам потребуется написать несколько сервисов. Поэтому это задание – хорошая возможность начать разработку одного из них. Вы можете считать это первым шагом к дипломной работе. И, самое главное, вы сейчас думаете о инфре, а потом только о логике(именно поэтому в больших компаниях есть платформенные команды). Ибо одновременно обо всем думать - гоняться на 2мя зайцами


- Делаем упор на архитектуру приложения(модульность) и обвязку вокруг него. Логику сервиса оставляем на десерт(т.е API чата c 20ю ручками и паданием с паникой и без логов это ~~не сдал~~ фаталити)
- Описываем документацию для вашего API(swagger). Было б очень классно иметь схемки/диаграммы(PlantUML)
- Протоколы и интеграции: grpc/http/graphql/socket/etc. - на выбор
- Конфигурация env/yaml
- Рейтлимиты/троттлинг. Пока в самом сервисе, просто для понимания устройства
- HealthCheck
- Логирование запросов и ответов(а можно ли автоматизировать?)
- Обработка ошибок и паник(хендлинг разных типов)
- Авторизация и аутентификация
- Валидация входящих данных
- Подключаем систему сбора метрик, чтобы можно было мониторить наши сервисы. Думаем и описываем список необходимых технических/бизнесовых метрик
- Реализуем трейсинг, чтобы отслеживать запросы через все сервисы и вовремя выявлять проблемы
- smth else?

## Схемы сервиса
```plantuml
@startuml
actor "user" as actor

component "Внутринние сервисы" as internal {

    artifact {
        agent "Котировки" as tracker
        database database_quotes [
            <b>Redis
            ====
            Кэш
        ]
    }
    
    rectangle rect_processing as "Бизнес логика" #aliceblue;line:blue;line.dotted;text:blue {
        artifact  {
            agent "Обработка данных" as processing 
            database database_processing [
                <b>Postgres
                ====
                Котировки
                ----
                Эмитенты
                ----
                Отчетность
            ]
        }
        
        
        artifact {
            agent "Пользователи и цели" as users_service 
            
            database database_user [
                <b>Mongo
                ====
                Цели
                ----
                Пользователи
            ]
        }
    }
    
    
    artifact {
        agent "Рассылка уведомлений" as notification
        database database_notification [
                <b>Mongo
                ====
                Журнал рассылок
            ]
    }
    
    artifact { 
        agent "Отчетсность" as parser
        database database_parser [
            <b>Redis
            ====
            Кэш
        ]
    }
    
    queue "отчетность" as rabbit_fundamentals
    queue "котировки" as rabbit_quotes 
    queue "задания на рассылку" as rabbit_notifications

}

cloud {
    agent "www.smartlab.ru" as externalData
}

artifact { 
    agent "API gateway" as api_gateway
    database database_auth[
        <b>Mongo
        ====
        Токены
    ] 
}

cloud {
    agent "Мосбиржа" as market
}

cloud {
    agent "телеграмм пользователей" as telegramm
}

externalData --> parser: Парсинг

tracker--[dotted]->rabbit_quotes
rabbit_quotes-[dotted]->processing

parser--[dotted]->rabbit_fundamentals
rabbit_fundamentals-[dotted]->processing

processing -[dotted]-> rabbit_notifications
rabbit_notifications-[dotted]->notification

users_service --> processing : gRPC

market-->tracker : "API"

notification-->telegramm

api_gateway <---> users_service : gRPC
actor <--> api_gateway : REST
@enduml
```

### Сборк
Перед сборкой контейнеров поднимаем общую сеть
```
docker network create fin-network
```