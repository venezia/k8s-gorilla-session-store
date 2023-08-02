package janitor

const (
	// GenericErrorType is the generic error for the demo package.
	GenericErrorType ErrorType = "generic"

	// TransientErrorType are errors when a failure happened, but we think it is a transient error
	//
	// Example causes include:
	//  - network connectivity disrupted
	//
	// It is reasonable to think that attempting the operation later may result in a happier outcome
	TransientErrorType ErrorType = "transient"

	// ExecutionErrorType is the error type that can return if something was attempted, but we cannot fulfill the request
	//
	// Repeated attempts are unlikely to result in a better outcome, as best as we can tell, the input was valid
	//
	// Examples would include:
	//  - misconfiguration of an upstream service (not disconnection, misconfiguration)
	ExecutionErrorType ErrorType = "execution"

	// InputErrorType is the error type that can return if something was attempted, but we cannot fulfill the request
	//
	// Repeated attempts are unlikely to result in a better outcome, as best as we can tell, the input was invalid
	//
	InputErrorType ErrorType = "input"

	// NotImplementedErrorType are errors when what asked simply is not implemented
	//
	// It is reasonable to expect that this change will never change
	NotImplementedErrorType ErrorType = "not-implemented"

	// GenericError is a coverall - hopefully never used
	GenericError ErrorMessage = "an error occurred"

	FeatureNotImplementedError ErrorMessage = "this feature of the store is not implemented for this provider"
)

type ErrorType string
type ErrorMessage string

type Error struct {
	Type    ErrorType
	Details interface{}
	Message ErrorMessage
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return string(e.Message) + ": " + e.Err.Error()
	}
	return string(e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}
