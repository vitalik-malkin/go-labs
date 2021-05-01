package seed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNighbour(t *testing.T) {
	f := newField(30)
	f.Set(1)
	f.Set(30)
	f.Set(15)
	f.Set(16)
	for i := 1; i < 31; i++ {
		var exp = i == 2 || i == 14 || i == 15 || i == 16 || i == 17 || i == 29
		assert.Equal(t, exp, f.IsNeighbour(i))
	}
}

func TestSimilarityDegree(t *testing.T) {

	{
		f1 := newField(30)
		f1.Set(1)
		f1.Set(30)
		f1.Set(15)
		f1.Set(16)
		f1.Reset(30)

		f2 := newField(30)
		f2.Set(1)
		f2.Set(30)
		f2.Set(15)
		f2.Set(16)

		assert.Equal(t, 0, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)
		f1.Set(1)
		f1.Set(30)
		f1.Set(15)
		f1.Set(16)

		f2 := newField(30)
		f2.Set(1)
		f2.Set(30)
		f2.Set(15)
		f2.Set(16)

		assert.Equal(t, 4, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)
		f1.Set(1)
		f1.Set(2)
		f1.Set(3)
		f1.Set(4)

		f2 := newField(30)
		f2.Set(5)
		f2.Set(6)
		f2.Set(7)
		f2.Set(8)

		assert.Equal(t, 0, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)
		f1.Set(1)
		f1.Set(2)
		f1.Set(3)
		f1.Set(4)

		f2 := newField(30)
		f2.Set(1)
		f2.Set(6)
		f2.Set(7)
		f2.Set(8)

		assert.Equal(t, 1, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)
		f1.Set(1)
		f1.Set(2)
		f1.Set(3)
		f1.Set(8)

		f2 := newField(30)
		f2.Set(1)
		f2.Set(6)
		f2.Set(7)
		f2.Set(8)

		assert.Equal(t, 2, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)

		f2 := newField(30)
		f2.Set(1)
		f2.Set(6)
		f2.Set(7)
		f2.Set(8)

		assert.Equal(t, 0, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)
		f1.Set(1)
		f1.Set(2)
		f1.Set(3)
		f1.Set(8)

		f2 := newField(30)

		assert.Equal(t, 0, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)

		f2 := newField(30)

		assert.Equal(t, 0, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)
		f1.Set(1)

		f2 := newField(30)
		f2.Set(1)

		assert.Equal(t, 1, SimilarityDegree(f1, f2))
	}

	{
		f1 := newField(30)
		f1.Set(30)

		f2 := newField(30)
		f2.Set(30)

		assert.Equal(t, 1, SimilarityDegree(f1, f2))
	}

}
