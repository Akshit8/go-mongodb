# go-mongodb
The all-in-one guide for using Mongodb with Go.

<img src=".github/assets/gomongodbBanner.png">

## Introduction
The repository is a sand-box that illustrates usage pattern of Mongodb in Go using the [go-mongo](https://github.com/mongodb/mongo-go-driver). The code is structured in a way making it easily pluggable with **clean/onion/hexagonal architecture.**
The purpose of creating this sand-box is

- minimal dependency of business logic implementation on Mongodb and mongo-driver.
- proper implementation of singleton client and it's usage across the app.
- solid unit testing of db queries in a concurrent environment.
- easy migration of db layer if required.

## Folder specs
- **.github** - 
- **config** -
- **entity** -
- **random** -
- **repository/mongo** -

## A note on update operation

## Using UUID instead of ObjectID.

## Running tests parallely

## References
[passing-data-to-goroutines](https://stackoverflow.com/questions/40326723/go-vet-range-variable-captured-by-func-literal-when-using-go-routine-inside-of-f)

## Author
**Akshit Sadana <akshitsadana@gmail.com>**

- Github: [@Akshit8](https://github.com/Akshit8)
- LinkedIn: [@akshitsadana](https://www.linkedin.com/in/akshit-sadana-b051ab121/)

## License
Licensed under the MIT License