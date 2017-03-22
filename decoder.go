package phpserialize

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Decoder struct {
	source *strings.Reader
}

func Decode(value string) (result interface{}, err error) {
	decoder := &Decoder{
		source: strings.NewReader(value),
	}
	result, err = decoder.DecodeValue()
	return
}

func (decoder *Decoder) DecodeValue() (value interface{}, err error) {
	if token, _, err := decoder.source.ReadRune(); err == nil {
		//type: nil
		if token == 'N' {
			err = decoder.expectElement(VALUES_SEPARATOR)
			return nil, err
		}
		decoder.expectElement(TYPE_VALUE_SEPARATOR)
		switch token {
		//type: bool
		case 'b':
			if rawValue, _, _err := decoder.source.ReadRune(); _err == nil {
				value = (rawValue == '1')	//true/false
			} else {
				err = errors.New("Can not read boolean value")
			}
			if err != nil {
				return nil, err
			}
			err = decoder.expectElement(VALUES_SEPARATOR)
		//type: int, int64, int32, int16, int8
		case 'i':
			if rawValue, _err := decoder.readUntil(VALUES_SEPARATOR); _err == nil {
				if tmpv, _err := strconv.Atoi(rawValue); _err != nil {
					err = fmt.Errorf("Can not convert %v to Int:%v", rawValue, _err)
				} else {
					value = int64(tmpv)
				}
			} else {
				err = errors.New("Can not read int value")
			}
		//type: float32, float64
		case 'd':
			if rawValue, _err := decoder.readUntil(VALUES_SEPARATOR); _err == nil {
				if value, _err = strconv.ParseFloat(rawValue, 64); _err != nil {
					err = fmt.Errorf("Can not convert %v to Float:%v", rawValue, _err)
				}
			} else {
				err = errors.New("Can not read float value")
			}
		//type: string
		case 's':
			value, err = decoder.decodeString()
			if err != nil {
				return nil, err
			}
			err = decoder.expectElement(VALUES_SEPARATOR)
		//type: map[interface{}]interface{}
		case 'a':
			value, err = decoder.decodeArray()
		//type: KvDataMap Object
		case 'O':
			value, err = decoder.decodeObject()
		}
	}
	return value, err
}

func (decoder *Decoder) decodeObject() (*KvDataMap, error) {
	value := &KvDataMap{}
	var err error

	if value.className, err = decoder.decodeString(); err != nil {
		return nil, err
	}
	if err = decoder.expectElement(TYPE_VALUE_SEPARATOR); err != nil {
		return nil, err
	}
	if value.members, err = decoder.decodeArray(); err != nil {
		return nil, err
	}

	return value, err
}

func (decoder *Decoder) decodeArray() (value map[interface{}]interface{}, err error) {
	value = make(map[interface{}]interface{})
	if rawArrlen, _err := decoder.readUntil(TYPE_VALUE_SEPARATOR); _err == nil {
		if arrLen, _err := strconv.Atoi(rawArrlen); _err != nil {
			err = fmt.Errorf("Can not convert array length %v to int:%v", rawArrlen, _err)
		} else {
			decoder.expectElement('{')
			for i := 0; i < arrLen; i++ {
				if k, _err := decoder.DecodeValue(); _err != nil {
					err = fmt.Errorf("Can not read array key %v", _err)
				} else if v, _err := decoder.DecodeValue(); _err != nil {
					err = fmt.Errorf("Can not read array value %v", _err)
				} else {
					switch t := k.(type) {
					default:
						err = fmt.Errorf("Unexpected key type %T", t)
					case string, int64, float64:
						value[k] = v
					}
				}
			}
			decoder.expectElement('}')
		}
	} else {
		err = errors.New("Can not read array length")
	}
	return value, err
}

func (decoder *Decoder) decodeString() (value string, err error) {
	if rawStrlen, _err := decoder.readUntil(TYPE_VALUE_SEPARATOR); _err == nil {
		if strLen, _err := strconv.Atoi(rawStrlen); _err != nil {
			err = errors.New(fmt.Sprintf("Can not convert string length %v to int:%v", rawStrlen, _err))
		} else {
			if err = decoder.expectElement('"'); err != nil {
				return
			}
			tmpValue := make([]byte, strLen, strLen)
			if nRead, _err := decoder.source.Read(tmpValue); _err != nil || nRead != strLen {
				err = errors.New(fmt.Sprintf("Can not read string content %v. Read only: %v from %v", _err, nRead, strLen))
			} else {
				value = string(tmpValue)
				err = decoder.expectElement('"')
			}
		}
	} else {
		err = errors.New("Can not read string length")
	}
	return value, err
}

func (decoder *Decoder) readUntil(stopByte byte) (string, error) {
	result := new(bytes.Buffer)
	var (
		token byte
		err   error
	)
	for {
		if token, err = decoder.source.ReadByte(); err != nil || token == stopByte {
			break
		} else {
			result.WriteByte(token)
		}
	}
	return result.String(), err
}

func (decoder *Decoder) expectElement(expectRune rune) error {
	token, _, err := decoder.source.ReadRune()
	if err != nil {
		err = fmt.Errorf("Can not read expected: %v", expectRune)
	} else if token != expectRune {
		err = fmt.Errorf("Read %v, but expected: %v", token, expectRune)
	}
	return err
}
