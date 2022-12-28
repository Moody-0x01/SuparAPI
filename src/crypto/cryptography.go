
/* NOTE This file is containing multiple functions to use for security. other
application functionality will be added shortly.

Funcs:

    -> sha256_: Responsible for hashing incoming sign up passwords. gets s a string and returns a hex digested hash.

    -> StoreTokenInToken: Responsible for generating the JWT for authentication. (Then it will be set in client cookies for better accessibility.)

    -> GetTokenFromJwt: Responsible for parsing the token, then returning the other access token.

*/

package crypto;

import (
    "crypto/sha256"
    b64 "encoding/base64"
    "fmt"
    "encoding/hex"
    "github.com/golang-jwt/jwt"
    "math/rand"
    "time"
    "strconv"
)

var SecretJwtKey []byte = []byte("d700977a3b1e3fd0145853702bdbb2a522530bb9707d314209d07b81dff3c17a")

/* DONE */

func Sha256_(s string) string {
    /* gets a string and returned a hash using sh256 */
    hash_ := sha256.New()
    hash_.Write([]byte(s))
    return hex.EncodeToString(hash_.Sum(nil))
}

// Makes a jwt that stores data to be sent to the client.
func StoreTokenInJWT(Token string) (string, error) {
 
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
        t, _ := claims["T"].(string)
       
        // if err != nil {
        //     return "", false
        // }

        return string(t), true

    } else {
	   // Invalid Token, return false.
        return "", false
    }
}

func GenerateAccessToken(salt string) string {
  
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
    return Sha256_(final);
}