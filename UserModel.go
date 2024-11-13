package main

type UserModel struct {
	Id      int
	Code    string `validate:"required"`
	Name    string `validate:"required"`
	Program string `validate:"required"`
}
