package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoreValidation(t *testing.T) {
	assert := assert.New(t)

	score := Score{
		Id: 123,
		UserId: 123,
		GameId: 123,
		DatePlayed: "2020-01-01",
		Score: 5,
	}

	ok, err := score.Validate();
	assert.True(ok, "Validation worked")
	assert.Nil(err, "No error message returned")

	score = Score{
		Id: 123,
		UserId: 0,
		GameId: 123,
		DatePlayed: "2020-01-01",
		Score: 5,
	}

	ok, err = score.Validate();
	assert.False(ok, "Validation worked")
	assert.NotNil(err, "No error message returned")
	assert.Regexp(regexp.MustCompile("User ID needs"), err.Error(), "No user id")

	score = Score{
		Id: 123,
		UserId: 123,
		GameId: 0,
		DatePlayed: "2020-01-01",
		Score: 5,
	}

	ok, err = score.Validate();
	assert.False(ok, "Validation worked")
	assert.NotNil(err, "No error message returned")
	assert.Regexp(regexp.MustCompile("Game ID needs"), err.Error(), "No user id")

	score = Score{
		Id: 123,
		UserId: 123,
		GameId: 123,
		DatePlayed: "20-01-01",
		Score: 5,
	}

	ok, err = score.Validate();
	assert.False(ok, "Validation worked")
	assert.NotNil(err, "No error message returned")
	assert.Regexp(regexp.MustCompile("Date Played needs"), err.Error(), "No user id")
}