package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func parseUser(data string) (*User, error) {
	var user User

	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, fmt.Errorf("user parse failed: %w", err)
	}

	return &user, nil
}

func (u User) String() string {
	return fmt.Sprintf("#%s email '%s' age %d", u.Id, u.Email, u.Age)
}

type UserList []User

func loadUserList(fileName string) (*UserList, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("user list file open error: %w", err)
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("user list file read error: %w", err)
	}

	var userList UserList

	if len(data) > 0 {
		err = json.Unmarshal(data, &userList)
		if err != nil {
			return nil, fmt.Errorf("user list file content decode failed: %w", err)
		}
	}

	return &userList, nil
}

func (ul UserList) Get(id string) *User {
	for _, item := range ul {
		if item.Id == id {
			return &item
		}
	}

	return nil
}

func (ul UserList) Has(id string) bool {
	return ul.Get(id) != nil
}

func (ul *UserList) Add(u *User) error {
	if ul.Has(u.Id) {
		return fmt.Errorf("Item with id %s already exists", u.Id)
	}

	*ul = append(*ul, *u)

	return nil
}

func (ul *UserList) Remove(id string) error {
	if !ul.Has(id) {
		return fmt.Errorf("Item with id %s not found", id)
	}

	var newList UserList

	for _, item := range *ul {
		if item.Id != id {
			newList = append(newList, item)
		}
	}

	*ul = newList

	return nil
}

func (ul UserList) DumpItem(id string) ([]byte, error) {
	if !ul.Has(id) {
		return nil, fmt.Errorf("Item with id %s not found", id)
	}

	item := ul.Get(id)

	data, err := json.Marshal((*item))
	if err != nil {
		return nil, fmt.Errorf("user encode error: %w", err)
	}

	return data, err
}

func (ul UserList) Dump() ([]byte, error) {
	if len(ul) == 0 {
		return nil, nil
	}

	data, err := json.Marshal(ul)
	if err != nil {
		return nil, fmt.Errorf("user list encode error: %w", err)
	}

	return data, nil
}

func (ul UserList) saveUserList(fileName string) error {
	data, err := ul.Dump()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("user list file open error: %w", err)
	}

	defer file.Close()

	written, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("user list save error: %w", err)
	}

	if len(data) > 0 && written < len(data) {
		return fmt.Errorf("user list save error: has unsaved data (%d bytes)", len(data)-written)
	}

	return nil
}

func (ul UserList) String() string {
	if len(ul) == 0 {
		return "Users list (empty)"
	}

	var s []string

	for _, i := range ul {
		s = append(s, i.String())
	}

	return fmt.Sprintf("Users list (%d items): %s", len(ul), strings.Join(s, ", "))
}
