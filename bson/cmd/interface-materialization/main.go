package main

import (
	"fmt"
	"log"
	"reflect"

	bson "go.mongodb.org/mongo-driver/bson"
	bc "go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type birthPlace struct {
	Country string
	City    string
}

type person struct {
	Name       string
	MiddleName string

	BirthPlace birthPlace
}

type data struct {
	Persons []person
}

type md struct {
	Name    string      `bson:"name"`
	Content interface{} `bson:"c"`
}

func x(dctx bc.DecodeContext, vr bsonrw.ValueReader, v reflect.Value) error {
	dr, err := vr.ReadDocument()
	return nil
}

func main() {

	rgb := bc.NewRegistryBuilder()
	rgb.RegisterDecoder(reflect.TypeOf(md{}), bc.ValueDecoderFunc(x))
	rg := rgb.Build()

	d := md{Name: "meta", Content: newData()}
	dBytes, err := bson.Marshal(&d)
	if err != nil {
		log.Fatal(err)
	}

	var d1 md

	err = bson.UnmarshalWithRegistry(rg, dBytes, &d1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%T, %T", d1, d1.Content)

	d1Cont, ok := d1.Content.(primitive.D)
	fmt.Printf("%v, %v", ok, d1Cont)

}

func newData() (d data) {
	d =
		data{
			Persons: []person{
				person{
					Name:       "Ivan",
					MiddleName: "Ivanovich",

					BirthPlace: birthPlace{
						City:    "Moscow",
						Country: "Russia",
					},
				},
			},
		}
	return
}
