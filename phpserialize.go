package phpserialize

import (
	"fmt"
)

const TYPE_VALUE_SEPARATOR = ':'
const VALUES_SEPARATOR = ';'

type KvDataMap struct {
	members   map[interface{}]interface{}
	className string
}

func NewKvDataMap() *KvDataMap {
	d := &KvDataMap{
		members: make(map[interface{}]interface{}),
	}
	return d
}

func (kvd *KvDataMap) GetClassName() string {
	return kvd.className
}

func (kvd *KvDataMap) SetClassName(cName string) {
	kvd.className = cName
}

func (kvd *KvDataMap) GetMembers() map[interface{}]interface{} {
	return kvd.members
}

func (kvd *KvDataMap) SetMembers(members map[interface{}]interface{}) {
	kvd.members = members
}

func (kvd *KvDataMap) GetPrivateMemberValue(memberName string) (value interface{}, found bool) {
	key := fmt.Sprintf("\x00%s\x00%s", kvd.className, memberName)
	value, found = kvd.members[key]
	return
}

func (kvd *KvDataMap) SetPrivateMemberValue(memberName string, value interface{}) {
	key := fmt.Sprintf("\x00%s\x00%s", kvd.className, memberName)
	kvd.members[key] = value
}

func (kvd *KvDataMap) GetProtectedMemberValue(memberName string) (value interface{}, found bool) {
	key := "\x00*\x00" + memberName
	value, found = kvd.members[key]
	return
}

func (kvd *KvDataMap) SetProtectedMemberValue(memberName string, value interface{}) {
	key := "\x00*\x00" + memberName
	kvd.members[key] = value
}

func (kvd *KvDataMap) GetPublicMemberValue(memberName string) (value interface{}, found bool) {
	value, found = kvd.members[memberName]
	return
}

func (kvd *KvDataMap) SetPublicMemberValue(memberName string, value interface{}) {
	kvd.members[memberName] = value
}
