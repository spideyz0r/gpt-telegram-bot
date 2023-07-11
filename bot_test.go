package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	openai "github.com/spideyz0r/openai-go"
)

func TestRunCommands(t *testing.T) {
	messages := []openai.Message{
		{
			Role:    "system",
			Content: default_system_role,
		},
	}

	tests := []struct {
		model       string
		temperature float64
		command     string
		expected    string
		args        string
	}{
		{
			model:       "gpt-4",
			temperature: 0.5,
			command:     "/info",
			args:        "",
			expected:    fmt.Sprintf("Model: gpt-4\nTemperature: 0.500000\nSystem role: %s", messages[0].Content),
		},
		{
			model:       "test",
			temperature: 1,
			command:     "/info",
			args:        "",
			expected:    fmt.Sprintf("Model: test\nTemperature: 1.000000\nSystem role: %s", messages[0].Content),
		},
		{
			model:       "test",
			temperature: 1,
			command:     "/reset",
			args:        "",
			expected:    "Conversation reset",
		},
		{
			model:       "test",
			temperature: 1,
			command:     "/temperature",
			args:        "",
			expected:    "Syntax is /temperature <float>. What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic. 0.8 is the default.",
		},
		{
			model:       "test",
			temperature: 1,
			command:     "/temperature",
			args:        "0.5",
			expected:    "Temperature set to 0.500000",
		},
		{
			model:       "test",
			temperature: 1,
			command:     "/role",
			args:        "You're a nerd",
			expected:    "Role set to You're a nerd",
		},
		{
			model:       "test",
			temperature: 1,
			command:     "/role",
			args:        "",
			expected:    "Syntax is /role <role>",
		},
	}
	for _, test := range tests {
		actual, updatedMessages, updatedTemperature := runCommands(test.command, test.args, messages, test.temperature, test.model)
		if actual != test.expected {
			t.Errorf("runCommands(%q, %q, %q, %q) returned %v, expected %v", test.command, test.args, messages, test.model, actual, test.expected)
			if !reflect.DeepEqual(updatedMessages, messages) {
				t.Errorf("Expected messages to be the same, but got different messages")
			}
			if test.args == "" {
				if updatedTemperature != test.temperature {
					t.Errorf("Expected temperature to be the same, but got different temperature")
				}
			}
		}
	}

}

func TestReadWhiteList(t *testing.T) {
	whitelist := make(map[int64]bool)
	expectedWhitelist := map[int64]bool{
		123456: true,
		987654: true,
	}

	// Create a temporary file with test data
	tempFile, err := ioutil.TempFile("", "whitelist")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.WriteString("123456\n987654\n")
	tempFile.Close()

	readWhiteList(&whitelist, tempFile.Name())

	if !reflect.DeepEqual(whitelist, expectedWhitelist) {
		t.Errorf("Expected whitelist: %v, but got: %v", expectedWhitelist, whitelist)
	}
}
