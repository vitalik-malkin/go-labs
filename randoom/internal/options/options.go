package options

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
}

func Load() Options {
	opts := Options{
		seedFile: "./../config/seed.txt",

		generatorVersion: 2,

		maxOfNum: 1000,

		fieldSetSize: 36,
		fieldSize:    100,

		maxOfNeighboringNumsInField:    0,
		fieldSimilarityDegree:          1,
		maxOfFieldsWithNeighboringNums: 0,

		genNumAttemptLimit:      400,
		genFieldAttemptLimit:    400,
		genFieldSetAttemptLimit: 80000000,
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

func (o *Options) MaxOfNum() int32 {
	return o.maxOfNum
}

func (o *Options) GeneratorVersion() int {
	return o.generatorVersion
}

func (o *Options) MaxOfFieldsWithNeighboringNums() int {
	return o.maxOfFieldsWithNeighboringNums
}
