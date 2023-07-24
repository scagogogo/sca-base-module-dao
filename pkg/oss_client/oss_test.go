package oss_client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPutPrivateObjectToOss(t *testing.T) {
	err := PutPrivateObjectToOss("maven-pom", "test", []byte("foo"))
	assert.Nil(t, err)
}
