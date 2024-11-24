package isunippets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupBy(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string
		Age  int
	}
	users := []*User{
		{Name: "Alice", Age: 20},
		{Name: "Bob", Age: 30},
		{Name: "Alice", Age: 25},
	}
	grouped := GroupBy(users, func(u *User) string {
		return u.Name
	})
	assert.Equal(2, len(grouped))
	assert.Equal(2, len(grouped["Alice"]))
	assert.Equal(1, len(grouped["Bob"]))
	assert.Equal(20, grouped["Alice"][0].Age)
	assert.Equal(25, grouped["Alice"][1].Age)
	assert.Equal(30, grouped["Bob"][0].Age)
}

func TestGroupById(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		ID   int64
		Name string
	}
	users := []*User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 1, Name: "Alice"},
	}
	grouped := GroupById(users, func(u *User) int64 {
		return u.ID
	})
	assert.Equal(2, len(grouped))
	assert.Equal(2, len(grouped[1]))
	assert.Equal(1, len(grouped[2]))
	assert.Equal("Alice", grouped[1][0].Name)
	assert.Equal("Alice", grouped[1][1].Name)
	assert.Equal("Bob", grouped[2][0].Name)
}
