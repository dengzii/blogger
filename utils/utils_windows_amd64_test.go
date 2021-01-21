package utils

import (
	"testing"
	"time"
)

func TestSetCreateTime(t *testing.T) {
	c := time.Now()
	if ChangeFileTimeAttr("", &c, nil, nil) != nil {

	}
}
