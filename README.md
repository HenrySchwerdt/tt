# TT (Terminal Time)


## Install

```sh
go install github.com/HenrySchwerdt/tt/cmd/tt@latest
```


## Projects

Projects can be created with the tt create <path> command.

```sh
tt create paper/rag
```

Would create two projects `paper` and `paper/rag` which can be queried serpeatly but when you log time to rag then also paper will increment.


## Logging Time

To log time you can just run `tt start <path>`. It is important to note that you can only start one entry at once. To finish run `tt end -m "What you did"`. Note that the message is optional but can be useful for the log command.