# wrappers

This repo contains installable wrappers for `kubectl` and `helm` that makes it harder to delete resources in user defined protected contexts.

## configuration

Configuration is done through config files in `$XDG_CONFIG_HOME/{kubectl,helm}/wrapper.yml`.

```yaml
# Example config

# command overrides the default path of the binary to run (default is /usr/bin/{kubectl,helm})
command: "/usr/local/bin/kubectl"
protectedContexts:
    - "my-precious-cluster"
# optionally you can define an even more annoying confirm string than the default one.
confirmString: "7H15-15-4-r34lLY-r34lly-4nn0y1n9-57R1N9-70-7YP3-1N"
```

## installation

Put the generated binaries in a location that has priority over `/usr/bin`. 
