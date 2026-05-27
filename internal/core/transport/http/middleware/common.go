package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/odelshchwank/BigProjectLesson/internal/core/logger"
	core_http_response "github.com/odelshchwank/BigProjectLesson/internal/core/transport/http/response"
	"go.uber.org/zap"
)

/*
Middleware(цепочка вызовов) для дурачков:

1. Прилетает сырой HTTP-запрос
2. Работает заложенная логика, которая шлепает уникальный id
3. Передаем запрос дальше
4. Когда основной код закончит работу, управление вернется в
мидлварку и ответ уйдет клиенту
*/

/*
Какой же путь именно в нашем случае с мидлваркой?
1. Запрос прилетает и сразу попадает на мидлварю RequestID чтобы можно было понимать какой именно запрос если что насрет
2. Запрос с уникальным ID попадает в логгер и логируется (ого)
3. После этого запрос переходит дальше и начинает отлавливать паники
4. Перед отработкой запроса он попадает на Trace, где фиксируется время "прихода" и статус-код со временем "ухода"
И потом по такому же маршруту в обратную сторону, чтобы всегда можно было посмотреть логи и проверить полный жизненный цикл каждого запроса
*/

/*
Почему ловим панику начиная не с момента прилета запроса?
Вопрос дискуссионный, однако по моему мнению, если запросу не был дан ID и он не был закинут в логгер,
и при этом будет выполняться логика, то пусть лучше прилетит необработанная паника которая прервет весь процесс,
чем всё будет тихо-мирно работать до момента какой-то жопы и мы в логах увидим гордый чистый лист
Ну и к тому же если паника будет обрабатываться раньше, чем присвоится ID, как узнать какой запрос ее вызвал?
*/

const (
	requestIDHeader = "X-Request-ID"
)

func RequestID() Middleware { // Просто функция-генератор для удобства жизни

	return func(next http.Handler) http.Handler { // Промежуточный слой,
		// который берет next обработчик и возвращает новый
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Адаптер из анонимной функции в объект, дабы все работало
			// Тут уже стартует сама логика

			// Генерим(или принимаем в запросе) id и пишем в заголовки
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			// Передаем запрос дальше, иначе запрос остановится здесь и никуда не пойдет
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				rw := core_http_response.NewResponseWriter(w)

				before := time.Now()

				log.Debug(
					">>> incoming HTTP request",
					zap.String("http_method", r.Method),
					zap.Time("time", time.Now().UTC()),
				)

				next.ServeHTTP(rw, r)

				log.Debug(
					"<<< done HTTP request",
					zap.Int("status_code", rw.GetStatusCode()),
					zap.Duration("latency", time.Since(before)),
				)
			})
	}
}

/*
Чтобы не отлавливать панику в каждой отдельной ручке и не дублировать код
по всем ручкам, выносим это дело в мидлварку
*/
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/*
				Чего мы хотим от этой мидлварки?
				Логировать панику и её причину (id, url,
				отдать статус-код, дать ответ об ошибке в JSON)
			*/

			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"got unexpected panic during handling HTTP request",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
