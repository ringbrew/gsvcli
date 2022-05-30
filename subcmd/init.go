package subcmd

func init() {
	Register(Init, func() SubCmd {
		return &InitSubCmd{}
	})
	Register(Gen, func() SubCmd {
		return &GenSubCmd{}
	})
}
