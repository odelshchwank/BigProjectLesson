package users_transport_http

import (
	"net/http"

	core_logger "github.com/odelshchwank/BigProjectLesson/internal/core/logger"
	core_http_request "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/request"
	core_http_response "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/response"
)

// DeleteUser    godoc
// @Summary      Удалить пользователя
// @Description  Удалить существующего пользователя из системы
// @Tags         users
// @Accept       query
// @Produce
// @Param        request body	DeleteUserRequest	true "DeleteUser тело запроса"
// @Success      204 			CreateUserResponse "Успешно удаленный пользователь"
// @Failure      400 {object}	core_http_response.ErrorResponse "BadRequest"
// @Failure      404 {object}	core_http_response.ErrorResponse "NotFound"
// @Failure      500 {object}	core_http_response.ErrorResponse "Internal server error"
// @Router       /users/{id} [DELETE]
func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
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

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)

		return
	}

	responseHandler.NoContentResponse()
}
