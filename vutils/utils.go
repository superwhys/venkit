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
func Reduce[T any, U any](arr []T, initial U, fn func(acc U, v T) U) U {
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

// Any checks if any element in the slice satisfies the predicate function
func Any[T any](arr []T, predicate func(T) bool) bool {
	for _, v := range arr {
		if predicate(v) {
			return true
		}
	}
	return false
}

// All checks if all elements in the slice satisfy the predicate function
func All[T any](arr []T, predicate func(T) bool) bool {
	for _, v := range arr {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Find returns the first element in the slice that satisfies the predicate function
func Find[T any](arr []T, predicate func(T) bool) (T, bool) {
	for _, v := range arr {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// GroupBy groups elements in the slice by a specified key function
func GroupBy[T any, K comparable](arr []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range arr {
		key := keyFunc(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Zip combines two slices into a slice of pairs (tuples)
func Zip[T1 any, T2 any](arr1 []T1, arr2 []T2) [][2]any {
	length := len(arr1)
	if len(arr2) < length {
		length = len(arr2)
	}
	result := make([][2]any, length)
	for i := 0; i < length; i++ {
		result[i] = [2]any{arr1[i], arr2[i]}
	}
	return result
}
