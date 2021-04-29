package options

type Options struct {
	seedFile string

	maxOfNum                    int32
	maxOfRepeatOfNumPerFieldSet int
	fieldSize                   int
	fieldSetSize                int
	maxOfNeighboringNumsInField int
	fieldSimilarityDegree       int

	genNumAttemptLimit      int
	genFieldAttemptLimit    int
	genFieldSetAttemptLimit int
}

func Load() Options {
	opts := Options{
		seedFile: "./../config/seed.txt",

		maxOfNum: 600,

		fieldSetSize: 200,
		fieldSize:    41,

		maxOfNeighboringNumsInField: 0,
		fieldSimilarityDegree:       1,

		genNumAttemptLimit:      300,
		genFieldAttemptLimit:    300,
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
