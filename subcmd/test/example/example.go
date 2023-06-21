package main

type Example struct {
	A         string
	B         int
	C         Inner
	Ar        [][]string
	MapEntity map[string][][][]Inner
}

type Inner struct {
	D string
	E int
}
