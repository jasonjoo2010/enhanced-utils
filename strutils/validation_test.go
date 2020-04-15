package strutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	assert.True(t, IsURL("http://a.dayima.com"))
	assert.True(t, IsURL("http://a.dayima.com/"))
	assert.True(t, IsURL("https://a.dayima.com/a/a/b/d/d-a"))
	assert.True(t, IsURL("https://a.dayima.com/a/a/b/d/d-a/a?topicId=33333&new=1"))
	assert.True(t, IsURL("https://a.dayima.com/a/a/b/d/d-a/a?topicId=33333&new=1#32322"))
	assert.True(t, IsURL("https://a.dayima.com/a/a/b/d/d-a/a?topicId=33333&c=%2f&new=1#32322"))
	assert.True(t, IsURL("dayima://topic/a/a/b/d/d-a/a"))
	assert.True(t, IsURL("dayima://topic/a/a/b/d/d-a/a?topicId=33333&new=1"))
	assert.True(t, IsURL("dayima://topic/a/a/b/d/d-a/a?topicId=33333&new=1&c=%2F"))
	assert.True(t, IsURL("ftp://a.dayima.com"))
	assert.False(t, IsURL("http3://a.dayima.com/"))
	assert.True(t, IsURL("https://a/a/a/b/d/d-a"))
	assert.False(t, IsURL("https:/a.dayima.com/a/a/b/d/d-a/a?topicId=33333&new=1"))
	assert.False(t, IsURL("://a.dayima.com/a/a/b/d/d-a/a?topicId=33333&new=1#32322"))
	assert.False(t, IsURL("dayima:/topic/a/a/b/d/d-a/a"))
	assert.False(t, IsURL("dayima//topic/a/a/b/d/d-a/a?topicId=33333&new=1"))
}

func TestEmail(t *testing.T) {
	assert.False(t, IsEmail(""))
	assert.False(t, IsEmail("test"))
	assert.False(t, IsEmail("test.com"))
	assert.False(t, IsEmail("www.test.com"))
	assert.True(t, IsEmail("test@test.com"))
	assert.True(t, IsEmail("test@www.test.com"))
	assert.True(t, IsEmail("test901@www.test.com"))
	assert.True(t, IsEmail("test901.god@www.test.com"))
	assert.True(t, IsEmail("test901.god-top@www.test.com"))
	assert.True(t, IsEmail("test901.god-top@www.test-first.com"))
}
