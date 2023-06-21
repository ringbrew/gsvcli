package main

type Example struct {
	A         string `json:"abc" bson:"sg" remark:"123" validate:"123"`
	B         int
	C         Inner
	Ar        [][]string
	MapEntity map[string][][][]Inner
	T         SetG[string]
}

type Inner struct {
	D string
	E int
}

type Outer struct {
	G string
	H int
}

type SetG[T comparable] map[T]struct{}
