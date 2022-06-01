package subcmd

func init() {
	Register(Init, func() SubCmd {
		return &InitSubCmd{}
	})
	Register(Domain, func() SubCmd {
		return &GenDomainSubCmd{}
	})
	Register(Grpc, func() SubCmd {
		return &GenGrpcSubCmd{}
	})
	Register(Install, func() SubCmd {
		return &InstallSubCmd{}
	})
}
