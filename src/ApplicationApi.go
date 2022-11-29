
/* NOTE This file is containing multiple functions to use for security. other
application functionality will be added shortly.

Funcs:

    -> sha256_: Responsible for hashing incoming sign up passwords. gets s a string and returns a hex digested hash.

    -> StoreTokenInToken: Responsible for generating the JWT for authentication. (Then it will be set in client cookies for better accessibility.)

    -> GetTokenFromJwt: Responsible for parsing the token, then returning the other access token.

*/

package main;

import (
    "crypto/sha256"
    b64 "encoding/base64"
    "fmt"
    "encoding/hex"
    "github.com/golang-jwt/jwt"
    "github.com/gin-gonic/gin"
    "math/rand"
    "time"
    "strconv"
)

var SecretJwtKey []byte = []byte("d700977a3b1e3fd0145853702bdbb2a522530bb9707d314209d07b81dff3c17a")

/* DONE */
func sha256_(s string) string {
    /* gets a string and returned a hash using sh256 */
    hash_ := sha256.New()
    hash_.Write([]byte(s))
    return hex.EncodeToString(hash_.Sum(nil))
}

// Makes a jwt that stores data to be sent to the client.
func StoreTokenInToken(Token string) (string, error) {
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "T": Token,
    })

    tokenString, err := token.SignedString(SecretJwtKey)
    return tokenString, err
}

// Extracts the token of the User that needs to be logged in.
func GetTokenFromJwt(TokenStr string) (string, bool) { 
    token, err := jwt.Parse(TokenStr, func(token *jwt.Token) (interface{}, error) {

        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        return SecretJwtKey, nil
    })

    if err != nil {
        fmt.Println("ERR: ", err)
        return "", false
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims["T"].(string), true
    } else {
	   // Invalid Token, return false.
        return "", false
    }
}


func GetFieldFromContext(c *gin.Context, field string) string {
    return c.Query(field) 
}

func generateAccessToken(salt string) string {
    /* 
        PYTHON bluepring:
            def TokenGen(self, salt: str):
                * _IV: list[str] = [chr(randint(0, 255)) for i in range(32)]
                * _S1 = ''.join(_IV)
                * _S_FINAL = _S1 + b64encode(salt.encode()).decode()
                * return sha256(_S_FINAL.encode()).hexdigest()
    */

    rand.Seed(time.Now().UnixNano())

    var SaltAsBytes []byte = []byte(salt)
    var _IV string = ""
    var threash_hold int = 255;
    var n int
    var nstring string

    for i := 0; i < 32; i++ {
        n = rand.Intn(threash_hold)
        nstring = strconv.Itoa(n)
        _IV += nstring
    }
    
    var final string = _IV + b64.StdEncoding.EncodeToString(SaltAsBytes);
    return sha256_(final);
}


func AuthenticateUserJWT(UserJWT string) Response {
    Token, Ok := GetTokenFromJwt(UserJWT)

    if Ok {
        User, err := getUserByToken(Token)
        
        if err != nil {
            // a db error.
            return MakeServerResponse(500, "Db Error. (line 108).")
        } else {
            // Returns the user if everything was alright.
            return MakeServerResponse(200, User)
        }

    } else {
        // JWT error
        return MakeServerResponse(500, "server could not decode the token. (line 117)")
    }
}