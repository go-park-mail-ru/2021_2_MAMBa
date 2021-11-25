# КиноПоиск
 Kino4U

 ## Деплой
 
 https://film4u.club/

## API

https://app.swaggerhub.com/apis/MAMBa/Film4U-API
 
 
 ## Авторы

 * [**Алексей Самойлов**](https://github.com/tr0llex) - *Фуллстек*
 * [**Марина Новикова**](https://github.com/fillinmar) - *Фронтенд*
 * [**Максим Дудник**](https://github.com/maksongold) - *Фронтенд*
 * [**Борис Кожуро**](https://github.com/BorisKoz) - *Бэкенд*

 ## Менторы
 * [**Елизавета Добрянская**](https://github.com/Betchika99) - *Фронтенд*
 * [**Екатерина Кириллова**](https://github.com/K1ola) - *Бэкенд*

 ## Ссылка на frontend

 https://github.com/frontend-park-mail-ru/2021_2_MAMBa


### tests
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
cat profile.cov | grep -v ".pb.go:\|mock" > profile1.cov
go tool cover -func profile1.cov