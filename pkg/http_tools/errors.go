package http_tools

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

type ErrorDto struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message any    `json:"message"`
} //	@name	ErrorDto

func LogError(err error) {
	GetLogger().Print(err)
}

func ErrorResponse(w http.ResponseWriter,
	_ *http.Request, status int, message any) {
	env := ErrorDto{
		Status:  status,
		Error:   http.StatusText(status),
		Message: message,
	}
	err := WriteJSON(w, status, env, nil)
	if err != nil {
		LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(w http.ResponseWriter,
	r *http.Request, err error) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}

	message := "the server encountered a problem and could not process your request"
	/*todo: if app.config.Env != config.ProdEnv {
		message = err.Error()
	}*/

	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func BadRequestResponse(w http.ResponseWriter,
	r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func RateLimitExceededResponse(w http.ResponseWriter,
	r *http.Request) {
	message := "rate limit exceeded"
	ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func UnauthorizedResponse(w http.ResponseWriter,
	r *http.Request, message string) {
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func ForbiddenResponse(w http.ResponseWriter,
	r *http.Request) {
	message := "user has no access to this resource"
	ErrorResponse(w, r, http.StatusForbidden, message)
}

func ConflictResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	resourceName string,
	identifier string,
	identifierValue string,
	jsonField string,
) {
	/*todo
	value := helpers.AnyToString(identifierValue)

	if err == nil || errors.Is(err, services.ErrRecordUniqueValue) {
		message := fmt.Sprintf(
			"%s with %s '%s' already exists",
			resourceName,
			identifier,
			value,
		)
		err := make(map[string]string)
		err[jsonField] = message
		app.errorResponse(w, r, http.StatusConflict, err)
	} else {
		app.serverErrorResponse(w, r, err)
	}
	*/
}

func NotFoundResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	resourceName string,
	identifier string, //nolint:unparam //should keep param
	identifierValue any,
	jsonField string,
) {
	/*todo
	value := helpers.AnyToString(identifierValue)

	if err == nil || errors.Is(err, services.ErrRecordNotFound) {
		message := fmt.Sprintf(
			"%s with %s '%s' doesn't exist",
			resourceName,
			identifier,
			value,
		)

		err := make(map[string]string)
		err[jsonField] = message

		app.errorResponse(w, r, http.StatusNotFound, err)
	} else {
		app.serverErrorResponse(w, r, err)
	}
	*/
}

func FailedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
