## SELF-BOT

Self bot is a library for building chat bots for Self. It provides a list of [components](./components/) you can reuse on your different bots.

Let's see an example on how to use it

```go
client, err := selfsdk.New(cfg)
// ...
b := bot.NewSelfBot(bot.SelfBotConfig{
    Client: client,
    Components: []components.CommandComponent{
        components.NewUserManagementComponent(components.UserManagementComponentConfig{
            Admins:    []string{"111222333"},
            StoreFile: "/tmp/users.json",
        }),
        components.NewHelpComponent(components.HelpComponentConfig{
            Header: "This bot allows you add users to the self sign up incentive program, see below the list of available commands",
        }),
        components.NewWelcomeComponent(components.WelcomeComponentConfig{
            Client:                     client,
            RegisteredUserWelcomeMsg:   rWelcomeMsg,
            UnRegisteredUserWelcomeMsg: uWelcomeMsg,
        }),
    },
})

b.Start()
// ...
```

