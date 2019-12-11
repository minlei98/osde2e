/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"fmt"

	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// groupData is the data structure used internally to marshal and unmarshal
// objects of type 'group'.
type groupData struct {
	Kind  *string           "json:\"kind,omitempty\""
	ID    *string           "json:\"id,omitempty\""
	HREF  *string           "json:\"href,omitempty\""
	Users *userListLinkData "json:\"users,omitempty\""
}

// MarshalGroup writes a value of the 'group' to the given target,
// which can be a writer or a JSON encoder.
func MarshalGroup(object *Group, target interface{}) error {
	encoder, err := helpers.NewEncoder(target)
	if err != nil {
		return err
	}
	data, err := object.wrap()
	if err != nil {
		return err
	}
	return encoder.Encode(data)
}

// wrap is the method used internally to convert a value of the 'group'
// value to a JSON document.
func (o *Group) wrap() (data *groupData, err error) {
	if o == nil {
		return
	}
	data = new(groupData)
	data.ID = o.id
	data.HREF = o.href
	data.Kind = new(string)
	if o.link {
		*data.Kind = GroupLinkKind
	} else {
		*data.Kind = GroupKind
	}
	data.Users, err = o.users.wrapLink()
	if err != nil {
		return
	}
	return
}

// UnmarshalGroup reads a value of the 'group' type from the given
// source, which can be an slice of bytes, a string, a reader or a JSON decoder.
func UnmarshalGroup(source interface{}) (object *Group, err error) {
	decoder, err := helpers.NewDecoder(source)
	if err != nil {
		return
	}
	data := new(groupData)
	err = decoder.Decode(data)
	if err != nil {
		return
	}
	object, err = data.unwrap()
	return
}

// unwrap is the function used internally to convert the JSON unmarshalled data to a
// value of the 'group' type.
func (d *groupData) unwrap() (object *Group, err error) {
	if d == nil {
		return
	}
	object = new(Group)
	object.id = d.ID
	object.href = d.HREF
	if d.Kind != nil {
		switch *d.Kind {
		case GroupKind:
			object.link = false
		case GroupLinkKind:
			object.link = true
		default:
			err = fmt.Errorf(
				"expected kind '%s' or '%s' but got '%s'",
				GroupKind,
				GroupLinkKind,
				*d.Kind,
			)
			return
		}
	}
	object.users, err = d.Users.unwrapLink()
	if err != nil {
		return
	}
	return
}
