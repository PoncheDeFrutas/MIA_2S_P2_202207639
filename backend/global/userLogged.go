package global

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Group struct {
	ID   string
	Type string
	Name string
}

type User struct {
	UserGroup Group
	Username  string
	Password  string
}

var (
	Users           = make(map[string][]User)
	Groups          = make(map[string][]Group)
	LoggedUser      string
	LoggedPartition string
)

func AddGroup(name string) error {
	for _, group := range Groups[name] {
		if group.ID != "0" {
			return fmt.Errorf("group already exists and is active")
		}
	}

	for i, group := range Groups[name] {
		if group.ID == "0" {
			newId := getNextGroupID()
			Groups[name][i].ID = newId
			return nil
		}
	}

	newId := getNextGroupID()
	newGroup := Group{ID: newId, Type: "G", Name: name}
	Groups[name] = append(Groups[name], newGroup)
	return nil
}

func AddUserToGroup(username, password, groupName string) error {
	for _, user := range Users[username] {
		if user.UserGroup.ID != "0" {
			return fmt.Errorf("user already exists and is active")
		}
	}

	for i, user := range Users[username] {
		if user.UserGroup.ID == "0" {
			groupList, exists := Groups[groupName]
			if !exists || len(groupList) == 0 {
				return fmt.Errorf("group does not exist")
			}

			var activeGroup *Group
			for j, group := range groupList {
				if group.ID != "0" {
					activeGroup = &groupList[j]
					break
				}
			}
			if activeGroup == nil {
				return fmt.Errorf("no active group found")
			}

			Users[username][i].UserGroup = *activeGroup
			Users[username][i].Password = password
			return nil
		}
	}

	groupList, exists := Groups[groupName]
	if !exists || len(groupList) == 0 {
		return fmt.Errorf("group does not exist")
	}

	var activeGroup *Group
	for i, group := range groupList {
		if group.ID != "0" {
			activeGroup = &groupList[i]
			break
		}
	}
	if activeGroup == nil {
		return fmt.Errorf("no active group found")
	}

	newUser := User{UserGroup: *activeGroup, Username: username, Password: password}
	Users[username] = append(Users[username], newUser)
	return nil
}

func RemoveGroup(name string) error {
	groupList, exists := Groups[name]
	if !exists {
		return fmt.Errorf("group does not exist")
	}

	groupUpdated := false

	for i := range groupList {
		if groupList[i].ID != "0" {
			groupList[i].ID = "0"
			groupUpdated = true

			for _, userList := range Users {
				for j := range userList {
					if userList[j].UserGroup.Name == name && userList[j].UserGroup.ID != "0" {
						userList[j].UserGroup.ID = "0"
					}
				}
			}
			break
		}
	}

	if groupUpdated {
		Groups[name] = groupList
		return nil
	}

	return fmt.Errorf("no active group found")
}

func RemoveUser(username string) error {
	userList, exists := Users[username]
	if !exists {
		return fmt.Errorf("user does not exist")
	}

	var userFound bool

	updatedUserList := make([]User, len(userList))
	copy(updatedUserList, userList)

	for i, user := range updatedUserList {
		if user.UserGroup.ID != "0" {

			updatedUserList[i].UserGroup.ID = "0"
			userFound = true
		}
	}

	if !userFound {
		return fmt.Errorf("no active user found")
	}

	Users[username] = updatedUserList
	return nil
}

func ChangeUserGroup(username, groupName string) error {
	userList, userExists := Users[username]
	if !userExists {
		return fmt.Errorf("user does not exist")
	}

	groupList, groupExists := Groups[groupName]
	if !groupExists {
		return fmt.Errorf("group does not exist")
	}

	var activeGroup *Group
	for i, group := range groupList {
		if group.ID != "0" {
			activeGroup = &groupList[i]
			break
		}
	}
	if activeGroup == nil {
		return fmt.Errorf("no active group found")
	}

	for i, user := range userList {
		if user.UserGroup.ID != "0" {
			userList[i].UserGroup = *activeGroup
			Users[username] = userList
			return nil
		}
	}
	return fmt.Errorf("no active user found")
}

func GetLoggedUser() (string, string, error) {
	if LoggedUser == "" {
		return "", "", fmt.Errorf("no user logged")
	}

	return LoggedUser, LoggedPartition, nil
}

func GetInfoUser(username string) User {
	for _, user := range Users[username] {
		if user.UserGroup.ID != "0" {
			return user
		}
	}

	return User{}
}

func LogUserIn(username, password, partition string) error {
	userList, exists := Users[username]
	if !exists {
		return fmt.Errorf("invalid user or password")
	}

	var validUser *User
	for _, user := range userList {
		if user.Password == password && user.UserGroup.ID != "0" {
			validUser = &user
			break
		}
	}

	if validUser == nil {
		return fmt.Errorf("invalid user or password")
	}

	LoggedUser = username
	LoggedPartition = partition

	return nil
}

func ClearData() {
	Users = make(map[string][]User)
	Groups = make(map[string][]Group)
}

func LogUserOut() (string, error) {
	if LoggedUser == "" {
		return "", fmt.Errorf("no user logged")
	}

	ClearData()
	temp := LoggedUser
	LoggedUser = ""
	LoggedPartition = ""
	return temp, nil
}

func IsUserLogged() bool {
	return LoggedUser != ""
}

func ParserUserData(data string) {
	ClearData()
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}

		id := strings.TrimSpace(parts[0])
		typ := strings.TrimSpace(parts[1])
		name := strings.TrimSpace(parts[2])

		switch len(parts) {
		case 3:
			Groups[name] = append(Groups[name], Group{ID: id, Type: typ, Name: name})
		case 5:
			username := strings.TrimSpace(parts[3])
			password := strings.TrimSpace(parts[4])
			Users[username] = append(Users[username], User{UserGroup: Group{ID: id, Type: typ, Name: name}, Username: username, Password: password})
		}
	}
}

func getNextGroupID() string {
	maxID := 0

	for _, groupList := range Groups {
		for _, group := range groupList {
			id, err := strconv.Atoi(group.ID)
			if err == nil && id > maxID {
				maxID = id
			}
		}
	}

	return strconv.Itoa(maxID + 1)
}

func ConvertToString() string {
	var sb strings.Builder

	groupMap := make(map[string]Group)
	userMap := make(map[string][]User)
	var usersWithoutGroup []User

	for _, groupList := range Groups {
		for _, group := range groupList {
			groupMap[group.ID] = group
		}
	}

	for _, userList := range Users {
		for _, user := range userList {
			if user.UserGroup.ID == "0" {
				usersWithoutGroup = append(usersWithoutGroup, user)
			} else {
				userMap[user.UserGroup.ID] = append(userMap[user.UserGroup.ID], user)
			}
		}
	}

	var ids []string
	for id := range groupMap {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool {
		id1, _ := strconv.Atoi(ids[i])
		id2, _ := strconv.Atoi(ids[j])
		return id1 < id2
	})

	for _, id := range ids {
		group := groupMap[id]
		sb.WriteString(strings.Join([]string{group.ID, group.Type, group.Name}, ",") + "\n")

		if users, exists := userMap[group.ID]; exists {
			userSet := make(map[string]bool)
			for _, user := range users {
				if !userSet[user.Username] {
					sb.WriteString(strings.Join([]string{
						user.UserGroup.ID,
						"U",
						user.UserGroup.Name,
						user.Username,
						user.Password,
					}, ",") + "\n")
					userSet[user.Username] = true
				}
			}
		}
	}

	for _, user := range usersWithoutGroup {
		sb.WriteString(strings.Join([]string{
			"0",
			"U",
			user.UserGroup.Name,
			user.Username,
			user.Password,
		}, ",") + "\n")
	}

	return sb.String()
}
