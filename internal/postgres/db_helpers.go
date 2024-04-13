package postgres

import (
	"fmt"
	"github.com/Vadim992/avito/internal/dto"
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

func validateBannersDataPatch(banner dto.PostPatchBanner) (string, error) {
	content := banner.Content
	var b strings.Builder

	if content != nil {
		if content.Title != nil {
			_, err := b.WriteString(fmt.Sprintf("title='%s',", *content.Title))

			if err != nil {
				return "", err
			}
		}

		if content.Text != nil {
			_, err := b.WriteString(fmt.Sprintf("text='%s',", *content.Text))

			if err != nil {
				return "", err
			}
		}

		if content.Url != nil {
			_, err := b.WriteString(fmt.Sprintf("url='%s',", *content.Url))

			if err != nil {
				return "", err
			}
		}

	}

	if banner.IsActive != nil {
		_, err := b.WriteString(fmt.Sprintf("is_active=%t,", *banner.IsActive))

		if err != nil {
			return "", err
		}
	}

	str := strings.TrimSpace(b.String())

	return str, nil
}
