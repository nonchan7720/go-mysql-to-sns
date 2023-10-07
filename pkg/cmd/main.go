package cmd

func Execute() {
	cmd := rootCommand()
	_ = cmd.Execute()
}
