package elmohttp

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)

func TestAddParametersToEmptyPath(t *testing.T) {

	path := addParametersToPath("", elmo.ConvertMapToValue(map[string]interface{}{}).(elmo.DictionaryValue))
	if path != "" {
		t.Error("add empty dictionary to empty path should still be an ampty path")
	}

	path = addParametersToPath("", elmo.ConvertMapToValue(map[string]interface{}{
		"pepper": "jalapeno"}).(elmo.DictionaryValue))

	if path != "?pepper=jalapeno" {
		t.Error("expected a jalapeno pepper, not ", path)
	}
}

func TestAddParametersToPathWithoutQuestionMark(t *testing.T) {

	path := addParametersToPath("/spices", elmo.ConvertMapToValue(map[string]interface{}{
		"pepper": "jalapeno"}).(elmo.DictionaryValue))

	if path != "/spices?pepper=jalapeno" {
		t.Error("expected a jalapeno pepper, not ", path)
	}
}

func TestAddMultipleParameters(t *testing.T) {

	path := addParametersToPath("/spices", elmo.ConvertMapToValue(map[string]interface{}{
		"pepper": "jalapeno",
		"amount": "3"}).(elmo.DictionaryValue))

	if path != "/spices?pepper=jalapeno&amount=3" {
		t.Error("expected 3 jalapeno peppers, not ", path)
	}
}

func TestUrlEncodedParameters(t *testing.T) {

	path := addParametersToPath("/spices", elmo.ConvertMapToValue(map[string]interface{}{
		"pepper": "jalapeno?",
		"amount": "\"3\""}).(elmo.DictionaryValue))

	if path != "/spices?pepper=jalapeno%3F&amount=%223%22" {
		t.Error("expected 3 jalapeno peppers, not ", path)
	}
}
