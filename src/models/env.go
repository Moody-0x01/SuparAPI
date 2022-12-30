
package models;

import (
	"log"
	// "github.com/joho/godotenv"
	// "os"
)

/*

func LoadDotEnv() {
	err := godotenv.Load(".env")

	if err != nil {
	    log.Fatalf("Error loading .env file")
	}
}

func GetEnv(k string) string {
	return os.Getenv("PORT")
}

*/


func loadEnv(params ...string) map[string]string {
	var EnvPath string = "./.env";
	
	if(len(params) >= 1) {
		EnvPath = params[0]
	}

	body, err := ioutil.ReadFile(EnvPath)
    
    if err != nil {
        log.Fatalf("unable to read file: %v", err)
    }
    
    Lines := strings.Split(string(body), "\n");
    Map := make(map[string]string)

    for i := 0; i < len(Lines); i++ {
    	// fmt.Println("Line", Lines[i])
    	next := strings.TrimSpace(Lines[i]);
    	trimmed := strings.Trim(next, "\n");
    	parsedLine := strings.Split(trimmed, "=");
    	key, val := parsedLine[0], parsedLine[1];
    	Map[key] = val;
    }

    return Map;
}

func GetEnv(key string) string {
    EnvMap := loadEnv();
    return EnvMap[key];
}