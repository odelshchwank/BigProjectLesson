package domain

// Пояснение работы Nullable[T any]
/*
1. Поле не передано
------------------
JSON: {}
NullableString:
	- Value: *nil
	- Set: false
------------------
2. Поле передано
------------------
JSON: {
	"phone_number": "+81111111111"
}
NullableString:
	- Value: *"+81111111111"
	- Set: true
------------------
3. В поле передано null
------------------
JSON: {
	"phone_number": null
}
NullableString:
	- Value: *nil
	- Set: true
------------------

По итогу для того, чтобы отличать случаи 1 и 3 нужно написать что-то вроде
if !Set && Value == nil -- случай 1
if Set && Value == nil -- случай 3
*/
type Nullable[T any] struct {
	Value *T
	Set   bool
}
