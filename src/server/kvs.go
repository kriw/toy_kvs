package server

var database = make(map[string]string)

func get(key string) string {
	return database[key]
}

func set(key string, value string) {
	database[key] = value
}
