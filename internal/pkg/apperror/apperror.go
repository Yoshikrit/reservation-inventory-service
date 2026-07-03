package apperror

import "fmt"

type Category string

const (
	CategoryBadRequest    Category = "BadRequest"
	CategoryUnauthorized  Category = "Unauthorized"
	CategoryNotFound      Category = "NotFound"
	CategoryConflict      Category = "Conflict"
	CategoryUnprocessable Category = "UnprocessableEntity"
	CategoryInternal      Category = "Internal"
)

type AppError struct {
	Category    Category `json:"category"`
	CodeNumber  int      `json:"code_number"`
	Description string   `json:"description"`
	Cause       error    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Description
}

var defaultAppErrorDefinitions = []AppError{
	{
		Category:    CategoryBadRequest,
		CodeNumber:  40000000,
		Description: "Invalid request: %v",
	},
	{
		Category:    CategoryBadRequest,
		CodeNumber:  40000001,
		Description: "Validation failed: %v",
	},
	{
		Category:    CategoryNotFound,
		CodeNumber:  40400000,
		Description: "Cannot find %s with id %v",
	},
	{
		Category:    CategoryConflict,
		CodeNumber:  40900000,
		Description: "%s with id %v already exists",
	},
	{
		Category:    CategoryUnprocessable,
		CodeNumber:  42200000,
		Description: "Unprocessable entity: %v",
	},
	{
		Category:    CategoryInternal,
		CodeNumber:  50000000,
		Description: "Unexpected internal error: %v",
	},
}

var definitionMap map[int]AppError

func init() {
	definitionMap = make(map[int]AppError, len(defaultAppErrorDefinitions))
	for _, def := range defaultAppErrorDefinitions {
		definitionMap[def.CodeNumber] = def
	}
}

// NewErrorWithDescription builds an AppError using the description verbatim — used when
// propagating an upstream error whose message is already well-formed (e.g. from a gRPC status).
func NewErrorWithDescription(category Category, codeNumber int, description string, cause error) *AppError {
	return &AppError{
		Category:    category,
		CodeNumber:  codeNumber,
		Description: description,
		Cause:       cause,
	}
}

func NewError(codeNumber int, cause error, args ...any) *AppError {
	def, ok := definitionMap[codeNumber]
	if !ok {
		def = definitionMap[50000000]
	}

	description := def.Description
	if len(args) > 0 {
		description = fmt.Sprintf(def.Description, args...)
	} else if cause != nil {
		description = fmt.Sprintf(def.Description, cause)
	}

	return &AppError{
		Category:    def.Category,
		CodeNumber:  def.CodeNumber,
		Description: description,
		Cause:       cause,
	}
}
