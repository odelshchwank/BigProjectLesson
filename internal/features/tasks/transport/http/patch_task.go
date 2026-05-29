package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/odelshchwank/BigProjectLesson/internal/core/domain"
	core_logger "github.com/odelshchwank/BigProjectLesson/internal/core/logger"
	core_http_request "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/request"
	core_http_response "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/response"
	core_http_types "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"        swaggertype:"string" example:"Планы на вечер"`
	Description core_http_types.Nullable[string] `json:"description"  swaggertype:"string" example:"null"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"    swaggertype:"boolean"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be null")
		}

		titleLength := len([]rune(*r.Title.Value))
		if titleLength < 1 || titleLength > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLength := len([]rune(*r.Description.Value))
			if descriptionLength < 1 || descriptionLength > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be null")
		}
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchTask     godoc
// @Summary      Обновить задачу
// @Description  Обновляет информацию об уже существующей в системе задаче
// @Description  1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description  2. **Явно передано значение**: `description`: "Попить пиво" - устанавливает новую описание для задачи
// @Description  3. **Передан null**: `description`: null - очищает поле в БД (set to NULL)
// @Description  Ограничения: `title` и `completed` не могут быть выставлены как null
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id       path      int                true           "ID изменяемой задачи"
// @Param        request  body	    PatchTaskRequest   true           "PatchTask тело запроса"
// @Success      201      {object}	PatchTaskResponse                 "Успешно измененная задача"
// @Failure      400      {object}	core_http_response.ErrorResponse  "Bad request"
// @Failure      409      {object}	core_http_response.ErrorResponse  "Conflict"
// @Failure      500      {object}	core_http_response.ErrorResponse  "Internal server error"
// @Router       /tasks/{id} [PATCH]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get TaskID path value",
		)

		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
