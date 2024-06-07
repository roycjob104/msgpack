# msgpack
To install:

```
go get github.com/roycjob104/msgpack
```

Sample usage model:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	msgpack "github.com/roycjob104/msgpack"
)

func main() {
	jsonData := `{"a": 1}`
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	resultEncode, err := msgpack.EncodeStringInterface(data)
	if err != nil {
		log.Fatalf("Error encoding data: %v", err)
	}
	fmt.Println("Encode json data =" + resultEncode)

	decodedData, err := msgpack.InitDecode(resultEncode)
	if err != nil {
		log.Fatalf("Error decoding data: %v", err)
	}

	fmt.Printf("Decoded data = %+v\n", decodedData)
	resultEncode, err = msgpack.EncodeStringInterface(decodedData.(map[string]interface{}))
	if err != nil {
		log.Fatalf("Error encoding data: %v", err)
	}
	fmt.Println("Encode json data =" + resultEncode)
}
```

### Running tests without larget tests
```
    $ go test
```

### Running tests with larget tests
```
    $ go test -runLargeTests=true
```