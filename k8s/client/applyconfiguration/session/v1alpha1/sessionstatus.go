/*
Some boilerplate
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SessionStatusApplyConfiguration represents an declarative configuration of the SessionStatus type for use
// with apply.
type SessionStatusApplyConfiguration struct {
	TTL *v1.Time `json:"ttl,omitempty"`
}

// SessionStatusApplyConfiguration constructs an declarative configuration of the SessionStatus type for use with
// apply.
func SessionStatus() *SessionStatusApplyConfiguration {
	return &SessionStatusApplyConfiguration{}
}

// WithTTL sets the TTL field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TTL field is set to the value of the last call.
func (b *SessionStatusApplyConfiguration) WithTTL(value v1.Time) *SessionStatusApplyConfiguration {
	b.TTL = &value
	return b
}
