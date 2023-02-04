
# AwqatSalah API Go Client

A go client for the Awqat Salah API. 


[![MIT License](https://img.shields.io/badge/License-MIT-orange.svg)](https://choosealicense.com/licenses/mit/)
[![Release Version](https://img.shields.io/github/v/release/cangiremir/awqatsalah-go)](https://img.shields.io/github/v/release/cangiremir/awqatsalah-go)

## Installing

Use `go get` to retrieve the library and add it to the your `GOPATH` workspace, or project's Go module dependencies.   

```bash
go get github.com/cangiremir/awqatsalah-go
```

## Usage

```
client, err := awqatsalah.New(awqatsalah.Credentials{
		Email:    '',
		Password: '',
	})

countries, err := c.Countries()
if err != nil {
    log.Fatalf("error getting countries: %v", err)
}

return countries
```

  
