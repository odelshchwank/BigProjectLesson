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
	core_http_utils "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/utils"
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
				"`FullName` must be beteween 3 and 100 symbols",
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
					"`PhoneNumeber` must starts with `+` symbol",
				)
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
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

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
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
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
