package components

import (
	"encoding/json"
	"log"

	"github.com/joinself/self-go-sdk/chat"
)

type UserManagementComponent struct {
	admins        map[string]struct{}
	storeFile     string
	store         Store
	componentName string
}

type UserManagementComponentConfig struct {
	Admins    []string
	StoreFile string
	Store     Store
}

func NewUserManagementComponent(cfg UserManagementComponentConfig) *UserManagementComponent {
	admins := make(map[string]struct{}, len(cfg.Admins))
	for _, v := range cfg.Admins {
		admins[v] = struct{}{}
	}

	um := UserManagementComponent{
		admins:        admins,
		storeFile:     cfg.StoreFile,
		store:         cfg.Store,
		componentName: "user_management_component",
	}
	um.load()

	return &um
}

func (u *UserManagementComponent) RecordCommands(r CommandRecorder) {
	u.recordListUsersCommand(r)
	u.recordAddUserCommand(r)
	u.recordRemoveUserCommand(r)

	r.SetAuth(u.auth)
}

func (i *UserManagementComponent) AfterStartHook(r CommandRecorder) {}

func (u *UserManagementComponent) recordListUsersCommand(r CommandRecorder) {
	r.RecordCommand(Command{
		Name:    "list_users",
		Summary: "Lists all the users with some permissions to manage this app.",
		Callback: func(cmd string, cm *chat.Message) string {
			resp := "This is the list of users with permissions to manage this app..."
			for k, _ := range u.admins {
				resp += "\n - @" + k
			}
			return resp
		},
	})
}

func (u *UserManagementComponent) recordAddUserCommand(r CommandRecorder) {
	r.RecordCommand(Command{
		Name:    "add_user",
		Summary: "Adds a new user to manage this app.",
		Callback: func(cmd string, cm *chat.Message) string {
			params := GetCommandParams(cm.Body)
			if len(params) == 0 {
				return "You must provide the self identifier of the user you want to add"
			}
			u.admins[params[0]] = struct{}{}
			u.save()

			return "@" + params[0] + " has been successfully <b>added</b> to the list of admins"
		},
	})
}

func (u *UserManagementComponent) recordRemoveUserCommand(r CommandRecorder) {
	r.RecordCommand(Command{
		Name:    "delete_user",
		Summary: "Removes a user from the management of this app.",
		Callback: func(cmd string, cm *chat.Message) string {
			params := GetCommandParams(cm.Body)
			if len(params) == 0 {
				return "You must provide the self identifier of the user you want to remove"
			}
			delete(u.admins, params[0])
			u.save()

			return "@" + params[0] + " has been successfully <b>removed</b> to the list of admins"
		},
	})
}

func (u *UserManagementComponent) save() {
	if u.store == nil {
		return
	}

	content, err := json.Marshal(u.admins)
	if err != nil {
		log.Println(err)
		return
	}

	u.store.Set(u.componentName, content)
}

func (u *UserManagementComponent) load() {
	if u.store == nil {
		return
	}

	content, err := u.store.Get(u.componentName)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(content, &u.admins)
	if err != nil {
		log.Println(err)
		return
	}
}

func (u *UserManagementComponent) auth(issuer, cmd string) bool {
	_, ok := u.admins[issuer]
	if !ok {
		return false
	}
	if cmd == "" {
		return true
	}
	//TODO: implement per command auth
	return true
}
