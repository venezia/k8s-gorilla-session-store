/*
Some boilerplate
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// SessionSpecApplyConfiguration represents an declarative configuration of the SessionSpec type for use
// with apply.
type SessionSpecApplyConfiguration struct {
	Data *string `json:"data,omitempty"`
}

// SessionSpecApplyConfiguration constructs an declarative configuration of the SessionSpec type for use with
// apply.
func SessionSpec() *SessionSpecApplyConfiguration {
	return &SessionSpecApplyConfiguration{}
}

// WithData sets the Data field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Data field is set to the value of the last call.
func (b *SessionSpecApplyConfiguration) WithData(value string) *SessionSpecApplyConfiguration {
	b.Data = &value
	return b
}
