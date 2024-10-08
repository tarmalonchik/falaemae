// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package entities

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// DirectionTypeInvalid is a DirectionType of type Invalid.
	DirectionTypeInvalid DirectionType = iota
	// DirectionTypeToVladikavkaz is a DirectionType of type To_vladikavkaz.
	DirectionTypeToVladikavkaz
	// DirectionTypeToTskhinvali is a DirectionType of type To_tskhinvali.
	DirectionTypeToTskhinvali
)

var ErrInvalidDirectionType = errors.New("not a valid DirectionType")

const _DirectionTypeName = "invalidto_vladikavkazto_tskhinvali"

var _DirectionTypeMap = map[DirectionType]string{
	DirectionTypeInvalid:       _DirectionTypeName[0:7],
	DirectionTypeToVladikavkaz: _DirectionTypeName[7:21],
	DirectionTypeToTskhinvali:  _DirectionTypeName[21:34],
}

// String implements the Stringer interface.
func (x DirectionType) String() string {
	if str, ok := _DirectionTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("DirectionType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x DirectionType) IsValid() bool {
	_, ok := _DirectionTypeMap[x]
	return ok
}

var _DirectionTypeValue = map[string]DirectionType{
	_DirectionTypeName[0:7]:                    DirectionTypeInvalid,
	strings.ToLower(_DirectionTypeName[0:7]):   DirectionTypeInvalid,
	_DirectionTypeName[7:21]:                   DirectionTypeToVladikavkaz,
	strings.ToLower(_DirectionTypeName[7:21]):  DirectionTypeToVladikavkaz,
	_DirectionTypeName[21:34]:                  DirectionTypeToTskhinvali,
	strings.ToLower(_DirectionTypeName[21:34]): DirectionTypeToTskhinvali,
}

// ParseDirectionType attempts to convert a string to a DirectionType.
func ParseDirectionType(name string) (DirectionType, error) {
	if x, ok := _DirectionTypeValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _DirectionTypeValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return DirectionType(0), fmt.Errorf("%s is %w", name, ErrInvalidDirectionType)
}

var errDirectionTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *DirectionType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = DirectionType(0)
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case int64:
		*x = DirectionType(v)
	case string:
		*x, err = ParseDirectionType(v)
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(v); verr == nil {
				*x, err = DirectionType(val), nil
			}
		}
	case []byte:
		*x, err = ParseDirectionType(string(v))
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(string(v)); verr == nil {
				*x, err = DirectionType(val), nil
			}
		}
	case DirectionType:
		*x = v
	case int:
		*x = DirectionType(v)
	case *DirectionType:
		if v == nil {
			return errDirectionTypeNilPtr
		}
		*x = *v
	case uint:
		*x = DirectionType(v)
	case uint64:
		*x = DirectionType(v)
	case *int:
		if v == nil {
			return errDirectionTypeNilPtr
		}
		*x = DirectionType(*v)
	case *int64:
		if v == nil {
			return errDirectionTypeNilPtr
		}
		*x = DirectionType(*v)
	case float64: // json marshals everything as a float64 if it's a number
		*x = DirectionType(v)
	case *float64: // json marshals everything as a float64 if it's a number
		if v == nil {
			return errDirectionTypeNilPtr
		}
		*x = DirectionType(*v)
	case *uint:
		if v == nil {
			return errDirectionTypeNilPtr
		}
		*x = DirectionType(*v)
	case *uint64:
		if v == nil {
			return errDirectionTypeNilPtr
		}
		*x = DirectionType(*v)
	case *string:
		if v == nil {
			return errDirectionTypeNilPtr
		}
		*x, err = ParseDirectionType(*v)
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(*v); verr == nil {
				*x, err = DirectionType(val), nil
			}
		}
	}

	return
}

// Value implements the driver Valuer interface.
func (x DirectionType) Value() (driver.Value, error) {
	return int64(x), nil
}

type NullDirectionType struct {
	DirectionType DirectionType
	Valid         bool
}

func NewNullDirectionType(val interface{}) (x NullDirectionType) {
	x.Scan(val) // yes, we ignore this error, it will just be an invalid value.
	return
}

// Scan implements the Scanner interface.
func (x *NullDirectionType) Scan(value interface{}) (err error) {
	if value == nil {
		x.DirectionType, x.Valid = DirectionType(0), false
		return
	}

	err = x.DirectionType.Scan(value)
	x.Valid = (err == nil)
	return
}

// Value implements the driver Valuer interface.
func (x NullDirectionType) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	// driver.Value accepts int64 for int values.
	return int64(x.DirectionType), nil
}
