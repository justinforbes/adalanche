// Code generated by "enumer -type=ObjectType -trimprefix=ObjectType -json"; DO NOT EDIT.

//
package engine

import (
	"encoding/json"
	"fmt"
)

const _ObjectTypeName = "OtherAttributeSchemaClassSchemaControlAccessRightGroupForeignSecurityPrincipalDomainDNSDNSNodeDNSZoneUserComputerManagedServiceAccountOrganizationalUnitBuiltinDomainContainerGroupPolicyContainerCertificateTemplateTrustServiceExecutable"

var _ObjectTypeIndex = [...]uint8{0, 5, 20, 31, 49, 54, 78, 87, 94, 101, 105, 113, 134, 152, 165, 174, 194, 213, 218, 225, 235}

func (i ObjectType) String() string {
	i -= 1
	if i >= ObjectType(len(_ObjectTypeIndex)-1) {
		return fmt.Sprintf("ObjectType(%d)", i+1)
	}
	return _ObjectTypeName[_ObjectTypeIndex[i]:_ObjectTypeIndex[i+1]]
}

var _ObjectTypeValues = []ObjectType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

var _ObjectTypeNameToValueMap = map[string]ObjectType{
	_ObjectTypeName[0:5]:     1,
	_ObjectTypeName[5:20]:    2,
	_ObjectTypeName[20:31]:   3,
	_ObjectTypeName[31:49]:   4,
	_ObjectTypeName[49:54]:   5,
	_ObjectTypeName[54:78]:   6,
	_ObjectTypeName[78:87]:   7,
	_ObjectTypeName[87:94]:   8,
	_ObjectTypeName[94:101]:  9,
	_ObjectTypeName[101:105]: 10,
	_ObjectTypeName[105:113]: 11,
	_ObjectTypeName[113:134]: 12,
	_ObjectTypeName[134:152]: 13,
	_ObjectTypeName[152:165]: 14,
	_ObjectTypeName[165:174]: 15,
	_ObjectTypeName[174:194]: 16,
	_ObjectTypeName[194:213]: 17,
	_ObjectTypeName[213:218]: 18,
	_ObjectTypeName[218:225]: 19,
	_ObjectTypeName[225:235]: 20,
}

// ObjectTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ObjectTypeString(s string) (ObjectType, error) {
	if val, ok := _ObjectTypeNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ObjectType values", s)
}

// ObjectTypeValues returns all values of the enum
func ObjectTypeValues() []ObjectType {
	return _ObjectTypeValues
}

// IsAObjectType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ObjectType) IsAObjectType() bool {
	for _, v := range _ObjectTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for ObjectType
func (i ObjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ObjectType
func (i *ObjectType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ObjectType should be a string, got %s", data)
	}

	var err error
	*i, err = ObjectTypeString(s)
	return err
}