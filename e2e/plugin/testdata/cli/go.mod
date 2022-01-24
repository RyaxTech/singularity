module github.com/RyaxTech/singularity/e2e-cli-plugin

go 1.13

require (
	github.com/spf13/cobra v1.0.0
	github.com/RyaxTech/singularity v0.0.0
)

replace github.com/RyaxTech/singularity => ./singularity_source
