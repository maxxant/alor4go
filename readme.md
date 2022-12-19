# Unofficial Golang package for ALOR OpenAPI V2

## использование

документация разработчика https://alor.dev/docs

Golang >= 1.19 (возможно и ниже, но не проверялось)

Примеры для запуска в папке ./example/ можно запускать через **go run**, например:

 ```go run ./example/loctime.go```

 либо из терминала как бинарники

 ```./example/loctime.go```

для авторизованного запуска необходимо поместить refresh токен в переменную окружения ALOR_REFRESH

## обновление кода из спецификации OpenAPI

готовый сгенерированный код находится в файле **alor4go.gen.go**, в случае необходимости обновить вручную нужно выполнить следующие шаги

обновить .yaml спецификацию с сайта алора:

```make get```

сгенерировать бойлерплейт:

```make gen```

из-за ошибок генератора для успешной компиляции необходимо пропатчить код

```make patch```


NOTE: для генерации бойлерплейта при работе **make gen** используется <https://github.com/deepmap/oapi-codegen> -  генератор OpenAPI3 для golang, команда для его установки:

```go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest```  - последний релиз

```go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@master```  - последняя версия из мастера