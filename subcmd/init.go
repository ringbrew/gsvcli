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
	Register(Http, func() SubCmd {
		return &GenHttpSubCmd{}
	})
	Register(Install, func() SubCmd {
		return &InstallSubCmd{}
	})
	Register(Version, func() SubCmd {
		return &VersionSubCmd{}
	})
}
