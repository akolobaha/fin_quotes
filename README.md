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
actor "Пользователь REST" as actor
actor "Пользователь Телеграм" as tg_user

component "Внутринние сервисы" as internal {

    rectangle rect_processing as "Бизнес логика" #aliceblue;line:blue;line.dotted;text:blue {
        artifact  {
            agent "Обработка данных" as processing 
            database database_processing [
                <b>MongoDB
                ====
                Отчетность
            ]
        }
        
        
        artifact {
            agent "Пользователи" as users_service 
            
            database database_user [
                <b>PostgreSQL
                ====
                Пользователи
                ----
                Токены
                ----
                Тикеры
            ]
        }

        

        artifact {
            agent "Цели" as targets_service 
            
            database database_targets [
                <b>PostgreSQL
                ====
                Цели
            ]
        }
    }
    
    
    artifact {
        agent "Рассылка уведомлений по email" as notification
    }
    
    artifact {
        agent "Рассылка уведомлений telegram" as notification_telegram
    }


    rectangle rect_data as "Данные" #aliceblue;line:blue;line.dotted;text:blue {
    artifact { 
        agent "Отчетность" as parser
    }

    artifact {
        agent "Котировки" as tracker
   
    }
    artifact {
        agent "Дивиденды" as dividends
   
    }
    }
    
    queue "RabbitMQ: отчетность" as rabbit_fundamentals
    queue "RabbitMQ: котировки" as rabbit_quotes 
    queue "RabbitMQ: дивиденды" as rabbit_dividends
    queue "RabbitMQ: задания на рассылку email" as rabbit_notifications
    
    queue "RabbitMQ: задания на рассылку telegram" as rabbit_notifications_telegram

}

cloud {
    agent "www.smartlab.ru" as externalData
}



cloud {
    agent "Мосбиржа" as market
}

cloud {
    agent "Пользователи" as telegramm
}

externalData --> parser: Парсинг

tracker-[dotted]->rabbit_quotes
dividends-[dotted]->rabbit_dividends
rabbit_quotes-[dotted]->processing
rabbit_dividends-[dotted]->processing

parser-[dotted]->rabbit_fundamentals
rabbit_fundamentals-[dotted]->processing

targets_service -[dotted]>rabbit_notifications

rabbit_notifications-[dotted]->notification

targets_service -[dotted]>rabbit_notifications_telegram
rabbit_notifications_telegram-[dotted]->notification_telegram

users_service <--> targets_service : gRPC
targets_service <--> processing : gRPC

market-->tracker : "API"
market-->dividends : "API"

notification-->telegramm
notification_telegram-->telegramm


 actor <--> users_service : REST
 tg_user <--> users_service : Telegram
@enduml
```

### Сборк
Перед сборкой контейнеров поднимаем общую сеть
```
docker network create fin-network
```