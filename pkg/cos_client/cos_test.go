package cos_client

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetObjectBytes(t *testing.T) {
	client, err := New()
	assert.Nil(t, err)
	bytes, err := GetObjectBytes(context.Background(), client, "/maven-central/ai.catboost:catboost-spark_2.12:0.25-rc3")
	assert.Nil(t, err)
	assert.True(t, len(bytes) != 0 )

	fmt.Println(string(bytes))
}

func TestPutObject(t *testing.T) {
	client, err := New()
	assert.Nil(t, err)
	err = PutObject(context.Background(), client, "/test", []byte("CC11001100"), ObjectACLPrivate)
	assert.Nil(t, err)
}

func TestPutPrivateObjectToOss(t *testing.T) {
	
}

func TestPutPublicObject(t *testing.T) {

}