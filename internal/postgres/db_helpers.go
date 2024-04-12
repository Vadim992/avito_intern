package postgres

import (
	"fmt"
	"strings"
)

func CreateWhereReq(b *strings.Builder, column string, val string, symbol string) {
	if b.Len() == 0 {
		b.WriteString("WHERE ")
	} else {
		b.WriteString("AND ")
	}

	b.WriteString(fmt.Sprintf("%s%s%s ", column, symbol, val))
}

func CreateLimitOffsetReq(b *strings.Builder, name string, val int) {
	b.WriteString(fmt.Sprintf("%s %d ", name, val))
}

func CreateInReqFromInt64(tagIds []int64) string {
	var b strings.Builder
	b.WriteString("(")

	for i, tadId := range tagIds {

		if i != len(tagIds)-1 {
			b.WriteString(fmt.Sprintf("%d, ", tadId))
			continue
		}

		b.WriteString(fmt.Sprintf("%d)", tadId))
	}

	return strings.TrimSpace(b.String())
}
