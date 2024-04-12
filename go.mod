module github.com/JesseCoretta/go-schemax


replace internal/rfc2079 => ./internal/rfc2079
replace internal/rfc2307 => ./internal/rfc2307
replace internal/rfc2798 => ./internal/rfc2798
replace internal/rfc3045 => ./internal/rfc3045
replace internal/rfc3671 => ./internal/rfc3671
replace internal/rfc3672 => ./internal/rfc3672
replace internal/rfc4512 => ./internal/rfc4512
replace internal/rfc4517 => ./internal/rfc4517
replace internal/rfc4519 => ./internal/rfc4519
replace internal/rfc4523 => ./internal/rfc4523
replace internal/rfc4524 => ./internal/rfc4524

require (
	github.com/JesseCoretta/go-rfc4512-antlr v0.0.0-20240408235351-7e2a20bb7df5
	github.com/JesseCoretta/go-stackage v1.0.3
	github.com/antlr4-go/antlr/v4 v4.13.0
	internal/rfc2079 v1.0.0
	internal/rfc2307 v1.0.0
	internal/rfc2798 v1.0.0
	internal/rfc3045 v1.0.0
	internal/rfc3671 v1.0.0
	internal/rfc3672 v1.0.0
	internal/rfc4512 v1.0.0
	internal/rfc4517 v1.0.0
	internal/rfc4519 v1.0.0
	internal/rfc4523 v1.0.0
	internal/rfc4524 v1.0.0
	internal/rfc4530 v1.0.0
)

replace internal/rfc4530 => ./internal/rfc4530

require golang.org/x/exp v0.0.0-20240404231335-c0f41cb1a7a0 // indirect

go 1.21.5
