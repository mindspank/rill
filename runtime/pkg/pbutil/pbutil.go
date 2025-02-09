package pbutil

import (
	"fmt"
	"math"
	"math/big"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/marcboeker/go-duckdb"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ToValue converts any value to a google.protobuf.Value. It's similar to
// structpb.NewValue, but adds support for a few extra primitive types.
func ToValue(v any, t *runtimev1.Type) (*structpb.Value, error) {
	switch v := v.(type) {
	// In addition to the extra supported types, we also override handling for
	// maps and lists since we need to use valToPB on nested fields.
	case map[string]any:
		var t2 *runtimev1.StructType
		if t != nil {
			t2 = t.StructType
		}
		v2, err := ToStruct(v, t2)
		if err != nil {
			return nil, err
		}
		return structpb.NewStructValue(v2), nil
	case []any:
		v2, err := ToListValue(v, t)
		if err != nil {
			return nil, err
		}
		return structpb.NewListValue(v2), nil
	// Handle types not handled by structpb.NewValue
	case int8:
		return structpb.NewNumberValue(float64(v)), nil
	case int16:
		return structpb.NewNumberValue(float64(v)), nil
	case uint8:
		return structpb.NewNumberValue(float64(v)), nil
	case uint16:
		return structpb.NewNumberValue(float64(v)), nil
	case time.Time:
		if t != nil && t.Code == runtimev1.Type_CODE_DATE {
			s := v.Format(time.DateOnly)
			return structpb.NewStringValue(s), nil
		}
		s := v.Format(time.RFC3339Nano)
		return structpb.NewStringValue(s), nil
	case float32:
		// Turning NaNs and Infs into nulls until frontend can deal with them as strings
		// (They don't have a native JSON representation)
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			return structpb.NewNullValue(), nil
		}
		return structpb.NewNumberValue(float64(v)), nil
	case float64:
		// Turning NaNs and Infs into nulls until frontend can deal with them as strings
		// (They don't have a native JSON representation)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return structpb.NewNullValue(), nil
		}
		return structpb.NewNumberValue(v), nil
	case *big.Int:
		// Evil cast to float until frontend can deal with bigs:
		v2, _ := new(big.Float).SetInt(v).Float64()
		return structpb.NewNumberValue(v2), nil
		// This is what we should do when frontend supports it:
		// s := v.String()
		// return structpb.NewStringValue(s), nil
	case duckdb.Decimal:
		// Evil cast to float until frontend can deal with bigs:
		denom := big.NewInt(10)
		denom.Exp(denom, big.NewInt(int64(v.Scale)), nil)
		v2, _ := new(big.Rat).SetFrac(v.Value, denom).Float64()
		return structpb.NewNumberValue(v2), nil
	case duckdb.Map:
		return ToValue(map[any]any(v), t)
	case map[any]any:
		var t2 *runtimev1.MapType
		if t != nil {
			t2 = t.MapType
		}
		v2, err := ToStructCoerceKeys(v, t2)
		if err != nil {
			return nil, err
		}
		return structpb.NewStructValue(v2), nil
	case duckdb.Interval:
		m := map[string]any{"months": v.Months, "days": v.Days, "micros": v.Micros}
		v2, err := ToStruct(m, nil)
		if err != nil {
			return nil, err
		}
		return structpb.NewStructValue(v2), nil
	case []byte:
		if t != nil && t.Code == runtimev1.Type_CODE_UUID {
			uid, err := uuid.FromBytes(v)
			if err == nil {
				return structpb.NewStringValue(uid.String()), nil
			}
		}
		return structpb.NewValue(v)
	default:
		// Default handling for basic types (ints, string, etc.)
		return structpb.NewValue(v)
	}
}

// ToStruct converts a map to a google.protobuf.Struct. It's similar to
// structpb.NewStruct(), but it recurses on valToPB instead of structpb.NewValue
// to add support for more types. Providing t as a type hint is optional.
func ToStruct(v map[string]any, t *runtimev1.StructType) (*structpb.Struct, error) {
	x := &structpb.Struct{Fields: make(map[string]*structpb.Value, len(v))}
	if t == nil {
		for k, v := range v {
			if !utf8.ValidString(k) {
				return nil, fmt.Errorf("invalid UTF-8 in string: %q", k)
			}
			var err error
			x.Fields[k], err = ToValue(v, nil)
			if err != nil {
				return nil, err
			}
		}
	} else {
		for _, f := range t.Fields {
			var err error
			x.Fields[f.Name], err = ToValue(v[f.Name], f.Type)
			if err != nil {
				return nil, err
			}
		}
	}
	return x, nil
}

// ToStructCoerceKeys converts a map with non-string keys to a google.protobuf.Struct.
// It attempts to coerce the keys to JSON strings. Providing t as a type hint is optional.
func ToStructCoerceKeys(v map[any]any, t *runtimev1.MapType) (*structpb.Struct, error) {
	var keyType, valType *runtimev1.Type
	if t != nil {
		keyType = t.KeyType
		valType = t.ValueType
	}

	x := &structpb.Struct{Fields: make(map[string]*structpb.Value, len(v))}
	for k1, v := range v {
		k2, ok := k1.(string)
		if !ok {
			// Encode k1 using ToValue (to correctly coerce time, big numbers, etc.) and then to JSON.
			// This yields more idiomatic/consistent strings than using fmt.Sprintf("%v", k1).
			val, err := ToValue(k1, keyType)
			if err != nil {
				return nil, err
			}

			data, err := val.MarshalJSON()
			if err != nil {
				return nil, err
			}

			// Remove surrounding quotes returned by MarshalJSON for strings
			k2 = trimQuotes(string(data))
		}

		var err error
		x.Fields[k2], err = ToValue(v, valType)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

// trimQuotes removes surrounding double quotes from a string, if present.
// Examples: `"10"` -> `10` and `10` -> `10`.
func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// ToListValue converts a map to a google.protobuf.List. It's similar to
// structpb.NewList(), but it recurses on valToPB instead of structpb.NewList
// to add support for more types. Providing t as a type hint is optional.
func ToListValue(v []interface{}, t *runtimev1.Type) (*structpb.ListValue, error) {
	var elemType *runtimev1.Type
	if t != nil {
		elemType = t.ArrayElementType
	}
	x := &structpb.ListValue{Values: make([]*structpb.Value, len(v))}
	for i, v := range v {
		var err error
		x.Values[i], err = ToValue(v, elemType)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

func FromValue(val *structpb.Value) (any, error) {
	switch v := val.GetKind().(type) {
	case *structpb.Value_StringValue:
		return v.StringValue, nil
	case *structpb.Value_BoolValue:
		return v.BoolValue, nil
	case *structpb.Value_NumberValue:
		return v.NumberValue, nil
	case *structpb.Value_NullValue:
		return nil, nil
	default:
		return nil, fmt.Errorf("value not supported: %v", v)
	}
}
