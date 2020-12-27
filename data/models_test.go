package data_test

import (
	"testing"

	"github.com/lanwupark/blog-api/data"
	"github.com/stretchr/testify/assert"
)

type A struct {
	Name     string
	Password string
	Age      int
	C        *C
}

type B struct {
	Name     string
	Password string
	Age      int
	C        *C
}

type C struct {
	Name string
}

func TestDuplication(t *testing.T) {
	assert := assert.New(t)
	a := A{
		Name:     "eanson",
		Password: "123456",
		Age:      18,
		C: &C{
			Name: "asasa",
		},
	}
	var b B
	err := data.DuplicateStructField(&a, &b)
	assert.NoError(err)
	assert.Equal(a.Name, b.Name)
	assert.Equal(a.Password, b.Password)
	assert.Equal(a.Age, b.Age)
	assert.Equal(a.C, b.C)
}

func TestDuplicationError(t *testing.T) {
	assert := assert.New(t)
	a := A{
		Name:     "eanson",
		Password: "123456",
		Age:      18,
	}
	b := 0
	err := data.DuplicateStructField(&a, b)
	t.Log(err)
	assert.Error(err)
}
