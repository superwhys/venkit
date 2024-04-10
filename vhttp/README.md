# superwhys/venkit/vhttp

## Example

### Method Get
```go
header := vhttp.DefaultJsonHeader()
params := vhttp.NewParams()

resp := vhttp.Default().Get(ctx, url, params, header)
```

### Method Post
```go
header := vhttp.DefaultFormUrlEncodedHeader()
form := vhttp.NewForm().Add("name", "superwhys").Encode()

resp := vhttp.Default().Post(ctx, url, nil, header, form)
```

### To get string resp
```go
respStr, err := resp.BodyString()
```

### To get Bytes resp
```go
respBytes, err := resp.BodyBytes()
```

### To get json resp
```go
err := resp.BodyJson(&respStruct)
```

### Headers
#### creat a new header
`vhttp.NewHeader()`
#### Add value to header
`header.Add(key, value)`

#### It also has a number of different headers built in
`vhttp.DefaultJsonHeader()`
`vhttp.DefaultFormUrlEncodedHeader()`
`vhttp.DefaultFormHeader()`

### Params
#### creat a new params
`vhttp.NewParams()`
#### Add value to params
`params.Add(key, value)`
#### Get value from params
`params.Get(key)`

### Form
#### creat a new form
`vhttp.NewForm()`
#### Add value to form
`form.Add(key, value)`
#### Encode form
`form.Encode()`
