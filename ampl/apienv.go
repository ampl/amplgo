package ampl

import "runtime/debug"

const modulePath = "github.com/ampl/amplgo"

// Read by amplkey (inherited by AMPL and amplkey)
func newAPIEnvironment() (*Environment, error) {
	env, err := NewEnv("", "")
	if err != nil {
		return nil, err
	}
	vars := map[string]string{
		"AMPL_API":         "amplgo",
		"AMPL_API_VERSION": moduleVersion(),
		"AMPL_API_ENV":     "",
	}
	for name, value := range vars {
		if err := env.AddEnvironmentVariable(name, value); err != nil {
			env.FreeEnv()
			return nil, err
		}
	}
	return env, nil
}

func moduleVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if info.Main.Path == modulePath {
		return info.Main.Version
	}
	for _, dep := range info.Deps {
		if dep.Path == modulePath {
			if dep.Replace != nil {
				return dep.Replace.Version
			}
			return dep.Version
		}
	}
	return ""
}
