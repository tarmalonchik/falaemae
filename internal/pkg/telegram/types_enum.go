// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package telegram

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// ParseModeHtml is a ParseMode of type Html.
	ParseModeHtml ParseMode = iota
	// ParseModeMarkdown is a ParseMode of type Markdown.
	ParseModeMarkdown
)

var ErrInvalidParseMode = errors.New("not a valid ParseMode")

const _ParseModeName = "htmlmarkdown"

var _ParseModeMap = map[ParseMode]string{
	ParseModeHtml:     _ParseModeName[0:4],
	ParseModeMarkdown: _ParseModeName[4:12],
}

// String implements the Stringer interface.
func (x ParseMode) String() string {
	if str, ok := _ParseModeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ParseMode(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ParseMode) IsValid() bool {
	_, ok := _ParseModeMap[x]
	return ok
}

var _ParseModeValue = map[string]ParseMode{
	_ParseModeName[0:4]:                   ParseModeHtml,
	strings.ToLower(_ParseModeName[0:4]):  ParseModeHtml,
	_ParseModeName[4:12]:                  ParseModeMarkdown,
	strings.ToLower(_ParseModeName[4:12]): ParseModeMarkdown,
}

// ParseParseMode attempts to convert a string to a ParseMode.
func ParseParseMode(name string) (ParseMode, error) {
	if x, ok := _ParseModeValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _ParseModeValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return ParseMode(0), fmt.Errorf("%s is %w", name, ErrInvalidParseMode)
}

var errParseModeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *ParseMode) Scan(value interface{}) (err error) {
	if value == nil {
		*x = ParseMode(0)
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case int64:
		*x = ParseMode(v)
	case string:
		*x, err = ParseParseMode(v)
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(v); verr == nil {
				*x, err = ParseMode(val), nil
			}
		}
	case []byte:
		*x, err = ParseParseMode(string(v))
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(string(v)); verr == nil {
				*x, err = ParseMode(val), nil
			}
		}
	case ParseMode:
		*x = v
	case int:
		*x = ParseMode(v)
	case *ParseMode:
		if v == nil {
			return errParseModeNilPtr
		}
		*x = *v
	case uint:
		*x = ParseMode(v)
	case uint64:
		*x = ParseMode(v)
	case *int:
		if v == nil {
			return errParseModeNilPtr
		}
		*x = ParseMode(*v)
	case *int64:
		if v == nil {
			return errParseModeNilPtr
		}
		*x = ParseMode(*v)
	case float64: // json marshals everything as a float64 if it's a number
		*x = ParseMode(v)
	case *float64: // json marshals everything as a float64 if it's a number
		if v == nil {
			return errParseModeNilPtr
		}
		*x = ParseMode(*v)
	case *uint:
		if v == nil {
			return errParseModeNilPtr
		}
		*x = ParseMode(*v)
	case *uint64:
		if v == nil {
			return errParseModeNilPtr
		}
		*x = ParseMode(*v)
	case *string:
		if v == nil {
			return errParseModeNilPtr
		}
		*x, err = ParseParseMode(*v)
		if err != nil {
			// try parsing the integer value as a string
			if val, verr := strconv.Atoi(*v); verr == nil {
				*x, err = ParseMode(val), nil
			}
		}
	}

	return
}

// Value implements the driver Valuer interface.
func (x ParseMode) Value() (driver.Value, error) {
	return int64(x), nil
}

type NullParseMode struct {
	ParseMode ParseMode
	Valid     bool
}

func NewNullParseMode(val interface{}) (x NullParseMode) {
	x.Scan(val) // yes, we ignore this error, it will just be an invalid value.
	return
}

// Scan implements the Scanner interface.
func (x *NullParseMode) Scan(value interface{}) (err error) {
	if value == nil {
		x.ParseMode, x.Valid = ParseMode(0), false
		return
	}

	err = x.ParseMode.Scan(value)
	x.Valid = (err == nil)
	return
}

// Value implements the driver Valuer interface.
func (x NullParseMode) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	// driver.Value accepts int64 for int values.
	return int64(x.ParseMode), nil
}
