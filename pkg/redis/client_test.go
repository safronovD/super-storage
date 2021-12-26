package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	config = &Config{
		Hostname: "0.0.0.0:49154",
		Username: "",
		Password: "",
		HashDB:   0,
		FileDB:   1,
	}
	fileID = "testFile"
	hash   = "0123"
)

func TestRedisHashStorageImpl(t *testing.T) {
	rds, err := NewRedisHashImpl(config)
	if err != nil {
		t.Fatal(err)
	}

	err = rds.WriteNewHash(hash, "ABD", 15)
	if err != nil {
		t.Fatal(err)
	}

	err = rds.WriteBlockPosition(fileID, hash)
	if err != nil {
		t.Fatal(err)
	}

	check, err := rds.Check(hash)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int64(1), check)

	result, err := rds.Read(fileID)
	if err != nil {
		t.Fatal(err)
	}

	exp := []string{"15", "ABD"}
	assert.Equal(t, exp, result[0])
}
