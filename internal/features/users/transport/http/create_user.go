package users_transport_http

import (
	"net/http"

	"github.com/odelshchwank/BigProjectLesson/internal/core/domain"
	core_logger "github.com/odelshchwank/BigProjectLesson/internal/core/logger"
	core_http_request "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/request"
	core_http_response "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name"    validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UserDTOResponse

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
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
