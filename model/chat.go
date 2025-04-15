package model

import (
	"fmt"
	"strings"
)

type Chat struct {
	Id        int
	CreatedAt string
	Desc      string
}

func (c Chat) ToString() string {
	return fmt.Sprintf("%d %s %s", c.Id, strings.Split(c.CreatedAt, " ")[0], c.Desc)
}
