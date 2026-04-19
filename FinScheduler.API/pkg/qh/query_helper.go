package qh

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Enum[T any] interface {
	~string
	IsValid() bool
}

func ParseUUIDs(queryParams url.Values, key string) ([]*uuid.UUID, error) {
	params := queryParams[key]
	if len(params) == 0 {
		return nil, nil
	}

	res := make([]*uuid.UUID, 0, len(params))
	for _, p := range params {
		uuidParsed, err := uuid.Parse(p)
		if err != nil {
			return nil, fmt.Errorf("invalid query parameter %q value %q: %w", key, p, err)
		}

		uuidCopy := uuidParsed
		res = append(res, &uuidCopy)
	}

	return res, nil
}

func ParseString(queryParams url.Values, key string) *string {
	param := queryParams.Get(key)
	if param == "" {
		return nil
	}
	return &param
}

func ParseDecimal(queryParams url.Values, key string) (*decimal.Decimal, error) {
	param, exists := parseSingleParam(queryParams, key)
	if !exists {
		return nil, nil
	}

	res, err := decimal.NewFromString(param)
	if err != nil {
		return nil, fmt.Errorf("invalid query parameter %q value %q: %w", key, param, err)
	}

	return &res, nil
}

func ParseInt32(queryParams url.Values, key string) (*int32, error) {
	param, exists := parseSingleParam(queryParams, key)
	if !exists {
		return nil, nil
	}

	parsed, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid query parameter %q value %q: %w", key, param, err)
	}

	res := int32(parsed)

	return &res, nil
}

func ParseBool(queryParams url.Values, key string) (*bool, error) {
	param, exists := parseSingleParam(queryParams, key)
	if !exists {
		return nil, nil
	}

	res, err := strconv.ParseBool(param)
	if err != nil {
		return nil, fmt.Errorf("invalid query parameter %q value %q: %w", key, param, err)
	}

	return &res, nil
}

func ParseTime(queryParams url.Values, key string) (*time.Time, error) {
	param, exists := parseSingleParam(queryParams, key)
	if !exists {
		return nil, nil
	}

	res, err := time.Parse(time.RFC3339, param)
	if err != nil {
		return nil, fmt.Errorf("invalid query parameter %q value %q: %w", key, param, err)
	}

	return &res, nil
}

func ParseEnums[T Enum[T]](queryParams url.Values, key string) ([]*T, error) {
	params := queryParams[key]
	if len(params) == 0 {
		return nil, nil
	}

	res := make([]*T, 0, len(params))
	for _, p := range params {
		categoryCopy := T(p)
		if !categoryCopy.IsValid() {
			return nil, fmt.Errorf("invalid query parameter %q value %q", key, p)
		}

		res = append(res, &categoryCopy)
	}

	return res, nil
}

func parseSingleParam(queryParams url.Values, key string) (string, bool) {
	params, exists := queryParams[key]
	if !exists || len(params) == 0 {
		return "", false
	}

	return params[0], true
}
