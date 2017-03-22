# phpserialize in Go

##Getting Start
phpseriallize is a package for serialize and unserialize an object like php
rules like:
	int  	i:N;
	bool  	b:1;
	string 	s:N;
	float  	d:2.3;
	...

##Example
```Go
package main
import (
    "github.com/sangasong/serialize/phpserialize"
    "fmt"
)
func main() {
        data := make(map[interface{}]interface{})
        data2 := make(map[interface{}]interface{})
        data2["test"] = true
        data2[int64(0)] = int64(5)
        data2["flt32"] = float32(2.3)
        data2["int64"] = int64(45)
        data3 := phpserialize.NewKvDataMap()
        data3.SetClassName("A")
        data3.SetPrivateMemberValue("a", 1)
        data3.SetProtectedMemberValue("b", 3.14)
        data3.SetPublicMemberValue("c", data2)
        data["arr"] = data2
        data["3"] = "s\"tr'}e"
        data["g"] = nil
        data["object"] = data3

        var (
                result    string
                decodeRes interface{}
                err       error
        )
        if result, err = phpserialize.Encode(data); err != nil {
                fmt.Println(fmt.Sprintf("encode data fail %v, %v", err, data))
                return
        }
        fmt.Println(fmt.Sprintf("encode data:%v to result:%v ok", data, result))
        if decodeRes, err = phpserialize.Decode(result); err != nil {
                fmt.Println(fmt.Sprintf("decode data fail %v, %v", err, result))
                return
        }
        fmt.Println(fmt.Sprintf("decode raw:%v to data:%v ok", result, decodeRes))
}

```

