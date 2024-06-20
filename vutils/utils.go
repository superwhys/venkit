package vutils

import (
	"math/rand"
	"os"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GenerateRandomString(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Map applies a function to each element in a slice and returns a new slice
func Map[T any, U any](arr []T, fn func(T) U) []U {
	result := make([]U, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}

// Reduce reduces a slice to a single value using a provided function and initial value
func Reduce[T any, U any](arr []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, v := range arr {
		result = fn(result, v)
	}
	return result
}

// Filter returns a new slice containing elements that satisfy a predicate function
func Filter[T any](arr []T, predicate func(T) bool) []T {
	result := []T{}
	for _, v := range arr {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
