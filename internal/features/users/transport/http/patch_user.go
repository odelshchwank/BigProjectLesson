package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/odelshchwank/BigProjectLesson/internal/core/domain"
	core_logger "github.com/odelshchwank/BigProjectLesson/internal/core/logger"
	core_http_request "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/request"
	core_http_response "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/response"
	core_http_types "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/types"
)

/*
При написании метода PatchUser возникла проблема
У нас есть 3 сценария с разным телом запроса в виде JSON
Разберем на примере "phone_number"
1. {} -- НИЧЕГО НЕ ДЕЛАТЬ С "phone_number"
2. {"phone_number"}: "+375111111111" -- установить новый номер телефона
3. {"phone_number"}: null -- удалить номер телефона
Беда в том, что при использовании в структуре PatchUserRequest что string, что *string, мы
не сможем никак отличить случай 1 и случай 3, т.к. если использовать просто строку, то там
и там придет просто "", а в случае с указателем -- <nil>
Для решения этой проблемы будем использовать internal/core/domain/nullable.go
*/
type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf(
				"`FullName` must be between 3 and 100 symbols",
			)
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf(
					"`PhoneNumber` must be between 10 and 15 symbols",
				)
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf(
					"`PhoneNumber` must start with `+` symbol",
				)
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser     godoc
// @Summary      Изменение пользователя
// @Description  Изменение информации об уже существующем в системе пользователе
// @Description  ### Логика обновления полей (Three-state logic):
// @Description  1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description  2. **Явно передано значение**: `phone_number`: "+375111111111" - устанавливает новый номер телефона в БД
// @Description  3. **Передан null**: `phone_number`: null - очищает поле в БД (set to NULL)
// @Description  Ограничения: `full_name` не может быть выставлен как null
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int                true           "ID изменяемого пользователя"
// @Param        request  body	    PatchUserRequest   true           "PatchUser тело запроса"
// @Success      201      {object}	PatchUserResponse                 "Успешно измененный пользователь"
// @Failure      400      {object}	core_http_response.ErrorResponse  "Bad request"
// @Failure      409      {object}	core_http_response.ErrorResponse  "Conflict"
// @Failure      500      {object}	core_http_response.ErrorResponse  "Internal server error"
// @Router       /users/{id} [PATCH]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
