package models

import (
	"strconv"
	"strings"
)

// Header ...
type Header struct {
	TraceId string `header:"x-trace-id"`
}

// QueryParams ...
type QueryParams struct {
	Name  string   `json:"name,omitempty" form:"name,omitempty" header:"name,omitempty"`
	Page  int      `json:"page,omitempty" form:"page,omitempty" header:"page,omitempty"`
	Limit int      `json:"limit,omitempty" form:"limit,omitempty" header:"limit,omitempty"`
	Data  []string `json:"data" form:"data"`
}

// CSIntList ...
type CSIntList string

// Values ...
func (c CSIntList) Values() []int {
	res := []int{}

	str := string(c)
	if str != "" {
		values := strings.Split(str, ",")
		for _, v := range values {
			if i, err := strconv.Atoi(v); err == nil {
				res = append(res, i)
			}
		}
	}

	if len(res) == 0 {
		return nil
	}
	return res
}
