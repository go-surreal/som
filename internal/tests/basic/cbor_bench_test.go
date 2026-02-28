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

	c := conv.AllTypes{
		AllTypes: model.AllTypes{
			Node:     som.NewNode[som.ULID](som.ULID(uuid.New().String())),
			FieldString:   "test string",
			FieldInt:      42,
			FieldFloat64:  3.14,
			FieldBool:     true,
			FieldTime:     now,
			FieldDuration: 5 * time.Hour,
			FieldUUID:     uuid.New(),
			FieldMonth:    time.January,
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

	c := conv.AllTypes{
		AllTypes: model.AllTypes{
			Node:     som.NewNode[som.ULID](som.ULID(uuid.New().String())),
			FieldString:   "test string",
			FieldInt:      42,
			FieldFloat64:  3.14,
			FieldBool:     true,
			FieldTime:     now,
			FieldDuration: 5 * time.Hour,
			FieldUUID:     uuid.New(),
			FieldMonth:    time.January,
		},
	}

	data, err := c.MarshalCBOR()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result conv.AllTypes
		err := result.UnmarshalCBOR(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConvRoundTrip benchmarks full marshal+unmarshal cycle
func BenchmarkConvRoundTrip(b *testing.B) {
	now := time.Now()

	c := conv.AllTypes{
		AllTypes: model.AllTypes{
			Node:     som.NewNode[som.ULID](som.ULID(uuid.New().String())),
			FieldString:   "test string",
			FieldInt:      42,
			FieldFloat64:  3.14,
			FieldBool:     true,
			FieldTime:     now,
			FieldDuration: 5 * time.Hour,
			FieldUUID:     uuid.New(),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := c.MarshalCBOR()
		if err != nil {
			b.Fatal(err)
		}

		var result conv.AllTypes
		err = result.UnmarshalCBOR(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConvMarshalSimple benchmarks marshaling of a simple model
func BenchmarkConvMarshalSimple(b *testing.B) {
	c := conv.SpecialTypes{
		SpecialTypes: model.SpecialTypes{
			Node: som.NewNode[som.UUID](som.UUID(uuid.New().String())),
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
	c := conv.SpecialTypes{
		SpecialTypes: model.SpecialTypes{
			Node: som.NewNode[som.UUID](som.UUID(uuid.New().String())),
			Name: "Test Group",
		},
	}

	data, err := c.MarshalCBOR()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result conv.SpecialTypes
		err := result.UnmarshalCBOR(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
