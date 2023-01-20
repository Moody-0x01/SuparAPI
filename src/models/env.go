
package models;


import (
	"log"
	"strings"
	"io/ioutil"
)



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
        if len(Lines[i]) > 0 {
            next := strings.TrimSpace(Lines[i]);
            trimmed := strings.Trim(next, "\n");
            parsedLine := strings.Split(trimmed, "=");
            
            if len(parsedLine) == 2 {
                key, val := parsedLine[0], parsedLine[1];
                Map[key] = val;    
            }
        }
    }

    return Map;
}

func GetEnv(key string) string {
    EnvMap := loadEnv();
    val, ok := EnvMap[key];
    
    if ok {
    	return val;
    } else {
    	panicMsg := "KEY WAS NOT FOUND IN ENV: " + key;
    	panic(panicMsg);
    	return val;
    }

}

