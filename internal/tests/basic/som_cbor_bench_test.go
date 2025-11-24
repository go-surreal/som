package basic

import (
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/conv"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/uuid"
)

// BenchmarkConvMarshal benchmarks the Conv type marshaling (which uses CBOR internally)
func BenchmarkConvMarshal(b *testing.B) {
	now := time.Now()

	c := conv.AllFieldTypes{
		AllFieldTypes: model.AllFieldTypes{
			Node:     som.NewNode(som.MakeID("all_field_types", uuid.New())),
			String:   "test string",
			Int:      42,
			Float64:  3.14,
			Bool:     true,
			Time:     now,
			Duration: 5 * time.Hour,
			UUID:     uuid.New(),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := c.MarshalCBOR()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConvUnmarshal benchmarks the Conv type unmarshaling
func BenchmarkConvUnmarshal(b *testing.B) {
	now := time.Now()

	c := conv.AllFieldTypes{
		AllFieldTypes: model.AllFieldTypes{
			Node:     som.NewNode(som.MakeID("all_field_types", uuid.New())),
			String:   "test string",
			Int:      42,
			Float64:  3.14,
			Bool:     true,
			Time:     now,
			Duration: 5 * time.Hour,
			UUID:     uuid.New(),
		},
	}

	data, err := c.MarshalCBOR()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result conv.AllFieldTypes
		err := result.UnmarshalCBOR(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConvRoundTrip benchmarks full marshal+unmarshal cycle
func BenchmarkConvRoundTrip(b *testing.B) {
	now := time.Now()

	c := conv.AllFieldTypes{
		AllFieldTypes: model.AllFieldTypes{
			Node:     som.NewNode(som.MakeID("all_field_types", uuid.New())),
			String:   "test string",
			Int:      42,
			Float64:  3.14,
			Bool:     true,
			Time:     now,
			Duration: 5 * time.Hour,
			UUID:     uuid.New(),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := c.MarshalCBOR()
		if err != nil {
			b.Fatal(err)
		}

		var result conv.AllFieldTypes
		err = result.UnmarshalCBOR(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConvMarshalSimple benchmarks marshaling of a simple model
func BenchmarkConvMarshalSimple(b *testing.B) {
	c := conv.Group{
		Group: model.Group{
			Node: som.NewNode(som.MakeID("group", uuid.New())),
			Name: "Test Group",
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := c.MarshalCBOR()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConvUnmarshalSimple benchmarks unmarshaling of a simple model
func BenchmarkConvUnmarshalSimple(b *testing.B) {
	c := conv.Group{
		Group: model.Group{
			Node: som.NewNode(som.MakeID("group", uuid.New())),
			Name: "Test Group",
		},
	}

	data, err := c.MarshalCBOR()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result conv.Group
		err := result.UnmarshalCBOR(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
