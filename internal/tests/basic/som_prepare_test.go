package main

import (
	"context"
	"github.com/go-surreal/som/tests/basic/gen/som"
	"regexp"
	"strings"
	"testing"
)

func prepareTestDatabase(ctx context.Context, tb testing.TB) (som.Client, func()) {
	tb.Helper()

	db, cleanup, err := prepareDatabase(ctx, toSlug(tb.Name()))
	if err != nil {
		tb.Fatal(err)
	}

	return db, cleanup
}

func toSlug(input string) string {
	// Remove special characters
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	processedString := reg.ReplaceAllString(input, " ")

	// Remove leading and trailing spaces
	processedString = strings.TrimSpace(processedString)

	// Replace spaces with dashes
	slug := strings.ReplaceAll(processedString, " ", "-")

	// Convert to lowercase
	slug = strings.ToLower(slug)

	return slug
}
