package utils

import "reflect"

func Ptr[T any](v T) *T {
	return &v
}

func Val[T any](v *T) T {
	return *v
}

func InterfacePointerIsNil(val any) bool {
	return reflect.ValueOf(val).IsNil()
}
