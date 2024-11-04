package main

import "reflect"

// RemoveZeroFields removes zero-value fields from a struct
func RemoveZeroFields(data interface{}) {
	v := reflect.ValueOf(data).Elem()
	removeZeroFieldsRecursive(v)
}

func removeZeroFieldsRecursive(v reflect.Value) {
	// Handle pointers and interfaces
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if !v.IsNil() {
			removeZeroFieldsRecursive(v.Elem())
		}
		return
	}

	// Only process structs
	if v.Kind() != reflect.Struct {
		return
	}

	// Iterate through struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		// fieldType := v.Type().Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			removeZeroFieldsRecursive(field)
		case reflect.Ptr, reflect.Interface:
			if !field.IsNil() {
				removeZeroFieldsRecursive(field.Elem())
			} else {
				field.Set(reflect.Zero(field.Type()))
			}
		case reflect.Slice, reflect.Map:
			if field.Len() == 0 {
				field.Set(reflect.Zero(field.Type()))
			}
		default:
			// Set field to zero value if it's already the zero value
			if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				field.Set(reflect.Zero(field.Type()))
			}
		}
	}
}
