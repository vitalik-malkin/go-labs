package options

import (
	"flag"
)

type Options struct {
	seedFile string

	generatorVersion int

	maxOfNum int32

	fieldSize    int
	fieldSetSize int

	maxOfRepeatOfNumPerFieldSet    int
	maxOfNeighboringNumsInField    int
	fieldSimilarityDegree          int
	maxOfFieldsWithNeighboringNums int

	genNumAttemptLimit      int
	genFieldAttemptLimit    int
	genFieldSetAttemptLimit int

	dop int
}

func Load() Options {
	var (
		maxOfNum              int
		fieldSetSize          int
		fieldSize             int
		fieldSimilarityDegree int
	)

	flag.IntVar(&maxOfNum, "max-of-num", 1, "")
	flag.IntVar(&fieldSetSize, "field-set-size", 1, "")
	flag.IntVar(&fieldSize, "field-size", 1, "")
	flag.IntVar(&fieldSimilarityDegree, "field-similarity-degree", 1, "")
	flag.Parse()

	opts := Options{
		seedFile: "./../config/seed.txt",

		generatorVersion: 2,

		maxOfNum: int32(maxOfNum),

		fieldSetSize: fieldSetSize,
		fieldSize:    fieldSize,

		maxOfNeighboringNumsInField:    0,
		fieldSimilarityDegree:          fieldSimilarityDegree,
		maxOfFieldsWithNeighboringNums: 0,

		genNumAttemptLimit:      211,
		genFieldAttemptLimit:    307,
		genFieldSetAttemptLimit: 80000000,

		dop: 4,
	}
	opts.maxOfRepeatOfNumPerFieldSet = (opts.fieldSize * opts.fieldSetSize) / int(opts.maxOfNum)
	return opts
}

func (o *Options) SeedFile() string {
	return o.seedFile
}

func (o *Options) FieldSize() int {
	return o.fieldSize
}

func (o *Options) FieldSetSize() int {
	return o.fieldSetSize
}

func (o *Options) SetFieldSetSize(val int) {
	o.fieldSetSize = val
}

func (o *Options) MaxOfNeighboringNumsInField() int {
	return o.maxOfNeighboringNumsInField
}

func (o *Options) FieldSimilarityDegree() int {
	return o.fieldSimilarityDegree
}

func (o *Options) MaxOfRepeatOfNumPerFieldSet() int {
	return o.maxOfRepeatOfNumPerFieldSet
}

func (o *Options) GenNumAttemptLimit() int {
	return o.genNumAttemptLimit
}

func (o *Options) GenFieldAttemptLimit() int {
	return o.genFieldAttemptLimit
}

func (o *Options) GenFieldSetAttemptLimit() int {
	return o.genFieldSetAttemptLimit
}

func (o *Options) SetGenFieldSetAttemptLimit(val int) {
	o.genFieldSetAttemptLimit = val
}

func (o *Options) MaxOfNum() int32 {
	return o.maxOfNum
}

func (o *Options) GeneratorVersion() int {
	return o.generatorVersion
}

func (o *Options) MaxOfFieldsWithNeighboringNums() int {
	return o.maxOfFieldsWithNeighboringNums
}
