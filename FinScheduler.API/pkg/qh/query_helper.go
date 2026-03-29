package qh

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/url"
	"strconv"
	"time"
)

type Enum[T any] interface {
	~string
	IsValid() bool
}

func ParseUUIDs(queryParams url.Values, key string) []*uuid.UUID {
	params := queryParams[key]
	if len(params) == 0 {
		return nil
	}

	var res []*uuid.UUID
	for _, p := range params {
		if uuidParsed, err := uuid.Parse(p); err == nil {
			uuidCopy := uuidParsed
			res = append(res, &uuidCopy)
		}
	}

	return res
}

func ParseUUID(queryParams url.Values, key string) *uuid.UUID {
	param := queryParams.Get(key)
	if param == "" {
		return nil
	}
	res, err := uuid.Parse(param)
	if err != nil {
		return nil
	}
	return &res
}

func ParseString(queryParams url.Values, key string) *string {
	param := queryParams.Get(key)
	if param == "" {
		return nil
	}
	return &param
}

func ParseDecimal(queryParams url.Values, key string) *decimal.Decimal {
	param := queryParams.Get(key)
	if param == "" {
		return nil
	}
	res, err := decimal.NewFromString(param)
	if err != nil {
		return nil
	}
	return &res
}

func ParseInt32(queryParams url.Values, key string) *int32 {
	param := queryParams.Get(key)
	if param == "" {
		return nil
	}

	parsed, err := strconv.ParseInt(param, 10, 32)

	if err != nil {
		return nil
	}

	res := int32(parsed)

	return &res
}

func ParseBool(queryParams url.Values, key string) *bool {
	param := queryParams.Get(key)
	if param == "" {
		return nil
	}

	res, err := strconv.ParseBool(param)

	if err != nil {
		return nil
	}

	return &res
}

func ParseTime(queryParams url.Values, key string) *time.Time {
	param := queryParams.Get(key)
	if param == "" {
		return nil
	}

	res, err := time.Parse(time.RFC3339, param)

	if err != nil {
		return nil
	}

	return &res
}

func ParseEnums[T Enum[T]](queryParams url.Values, key string) []*T {
	params := queryParams[key]
	if len(params) == 0 {
		return nil
	}

	var res []*T
	for _, p := range params {
		categoryCopy := T(p)
		if categoryCopy.IsValid() {
			res = append(res, &categoryCopy)
		}
	}

	return res
}
