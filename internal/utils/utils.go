/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"strings"
)

// Function to convert map[string]string to Java properties format
func MapToJavaProperties(m map[string]string) string {
	var sb strings.Builder

	for key, value := range m {
		// Escape special characters in key and value
		escapedKey := escapeJavaPropertiesString(key)
		escapedValue := escapeJavaPropertiesString(value)

		// Append key-value pair to the string builder
		sb.WriteString(escapedKey)
		sb.WriteString("=")
		sb.WriteString(escapedValue)
		sb.WriteString("\n")
	}

	return sb.String()
}

// Function to escape special characters in Java properties format
func escapeJavaPropertiesString(s string) string {
	// Define special characters that need to be escaped
	escapeChars := map[string]string{
		"\\": "\\\\",
		"\n": "\\n",
		"\r": "\\r",
		"\t": "\\t",
		"=":  "\\=",
		":":  "\\:",
	}

	var sb strings.Builder

	for _, char := range s {
		charStr := string(char)

		// Escape special characters if present in the map
		if escaped, ok := escapeChars[charStr]; ok {
			sb.WriteString(escaped)
		} else {
			sb.WriteString(charStr)
		}
	}

	return sb.String()
}

func GetIntPointer(value int32) *int32 {
	return &value
}
