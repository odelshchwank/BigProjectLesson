package users_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

type CreateUserResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	/*
		Пояснение нахуя тут вообще мидлварка

			Для работы CreateUser нужно знать id запроса, который мы или получаем из хеддера,
		или генерируем сами(что и описано в логике мидлварки)
			Резонный вопрос - а нахуя усложнять, если можно в начале CreateUser это описать?
			Проблема в том, что хэндлеров может быть не 1 и не 2, а там тоже нужен id
		и нам теперь везде писать один и тот же код? Пахнет калом данная затея. Помним
		про DRY -- don't repeat yourself. И это только для users...

			Тогда следующие предложение -- вынести в отдельную функцию. Казалось бы логично
		Но обращение к функции это строка кода и её нужно везде шлепать в каждой ручке,
		а это опять дублирование одного и того же кода и лишний объем в ручке, к тому же
		если забыть ее добавить, то прога не упадет, но лога не будет и это черевато бедой
	*/

	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("ПОЛУНДРА НАХУЙ")
	}
}
