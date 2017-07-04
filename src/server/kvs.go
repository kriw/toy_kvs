package server

var database = make(map[string]string)

func Get(key string) string {
	return database[key]
}

func Set(key string, value string) {
	database[key] = value
}
