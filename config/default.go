package config

var defaultConfigToml = []byte(`
movies_dirs = [
    "/jellyfin/movies",
]

# Shows diretories will promp for a subdir before add
shows_dirs = [
    "/jellyfin/shows",
]
`)
