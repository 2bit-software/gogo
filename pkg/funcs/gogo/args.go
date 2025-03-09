package gogo

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/jessevdk/go-flags"
)

func ParseArgs(options any, args []string) ([]string, error) {
	return flags.ParseArgs(options, args)
}

// HydrateFromPositional fills struct fields with values from positional arguments
// based on the "order" struct tag. Fields that already have non-zero values are skipped.
//
// This function allows defining command argument handling in a centralized way,
// supporting both flag-based and positional arguments in a priority order.
func HydrateFromPositional(opts any, positional []string) error {
	val := reflect.ValueOf(opts)

	// Ensure we're working with a pointer to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct, got %T", opts)
	}

	val = val.Elem() // Get the struct value
	typ := val.Type()

	// Track fields by their position
	type fieldEntry struct {
		index    int  // Index in the struct
		position int  // Position from order tag
		isSet    bool // Whether field already has a value
	}

	// Collect fields with order tags
	fields := make([]fieldEntry, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if !fieldVal.CanSet() {
			continue
		}

		// Get and parse the order tag
		orderTag := field.Tag.Get("order")
		if orderTag == "" {
			continue
		}

		position, err := strconv.Atoi(orderTag)
		if err != nil {
			return fmt.Errorf("invalid order tag for field %s: %w", field.Name, err)
		}

		// Check if field already has a value
		isSet := !fieldVal.IsZero()

		fields = append(fields, fieldEntry{
			index:    i,
			position: position,
			isSet:    isSet,
		})
	}

	// Sort by position
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].position < fields[j].position
	})

	// Track which positional args have been used
	usedArgs := make([]bool, len(positional))

	// First pass: assign positional args to fields in order
	for _, field := range fields {
		if field.isSet {
			continue
		}

		// Skip if we've run out of positional args
		if field.position >= len(positional) {
			continue
		}

		// Skip if this arg has already been used
		if usedArgs[field.position] {
			continue
		}

		fieldVal := val.Field(field.index)
		fieldType := typ.Field(field.index)

		// Set the value based on field type
		if err := setFieldFromString(fieldVal, positional[field.position], fieldType.Name); err != nil {
			return err
		}

		usedArgs[field.position] = true
	}

	// Second pass: use any remaining args for unfilled fields
	argIndex := 0
	for _, field := range fields {
		if field.isSet || !val.Field(field.index).IsZero() {
			continue
		}

		// Find the next unused arg
		for argIndex < len(positional) && usedArgs[argIndex] {
			argIndex++
		}

		if argIndex >= len(positional) {
			break // No more args
		}

		fieldVal := val.Field(field.index)
		fieldType := typ.Field(field.index)

		// Set the value
		if err := setFieldFromString(fieldVal, positional[argIndex], fieldType.Name); err != nil {
			return err
		}

		usedArgs[argIndex] = true
		argIndex++
	}

	return nil
}

// setFieldFromString sets a field value from a string, with appropriate type conversion
func setFieldFromString(field reflect.Value, value string, fieldName string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
		return nil

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for %s: %w", fieldName, err)
		}
		field.SetBool(boolValue)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value for %s: %w", fieldName, err)
		}
		if field.OverflowInt(intValue) {
			return fmt.Errorf("value %d overflows %s", intValue, fieldName)
		}
		field.SetInt(intValue)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value for %s: %w", fieldName, err)
		}
		if field.OverflowUint(uintValue) {
			return fmt.Errorf("value %d overflows %s", uintValue, fieldName)
		}
		field.SetUint(uintValue)
		return nil

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value for %s: %w", fieldName, err)
		}
		if field.OverflowFloat(floatValue) {
			return fmt.Errorf("value %f overflows %s", floatValue, fieldName)
		}
		field.SetFloat(floatValue)
		return nil

	default:
		return fmt.Errorf("unsupported field type for %s: %s", fieldName, field.Kind())
	}
}
