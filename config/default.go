package config

var defaultConfigToml = []byte(`
# You can choose from before download
# Items appear in the order they are declared
download_dirs = [
    "/jellyfin/movies",
    "/jellyfin/shows",
]
`)
