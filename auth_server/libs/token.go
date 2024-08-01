package libs

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Token struct {
	Token     string    `json:"token"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
	MatchedAt time.Time `json:"matched_at"`
}

func (t *Token) Auth(publicKey string) {
	t.PublicKey = publicKey
	t.MatchedAt = time.Now()
	fmt.Printf("Token %s authenticated with public key %s\n", t.Token, publicKey)
}

type TokenManager struct {
	tokenList         *list.List
	noMatchTimeoutSec int
	matchTimeoutSec   int
	mu                sync.Mutex
}

func InitTokenManager() *TokenManager {
	return &TokenManager{
		tokenList:         list.New(),
		noMatchTimeoutSec: 60,
		matchTimeoutSec:   600,
	}
}

func (tm *TokenManager) NewToken() string {
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tm.AddToken(token)
	return token
}

func (tm *TokenManager) AddToken(token string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.tokenList.PushBack(&Token{
		Token:     token,
		CreatedAt: time.Now(),
	})
}

// 检查token是否存在，是否超时
// 不要着急退出，扫完成一遍，处理完所有的超时判断后，才退出
func (tm *TokenManager) CheckToken(token string) (*Token, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	var foundToken *Token

	for e := tm.tokenList.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Token)
		var timeoutSec int

		// 根据是否匹配来设置超时时间
		if t.PublicKey != "" {
			timeoutSec = tm.matchTimeoutSec
		} else {
			timeoutSec = tm.noMatchTimeoutSec
		}

		// 检查是否超时
		if time.Since(t.CreatedAt).Seconds() > float64(timeoutSec) {
			// Token 已超时，移除它
			tm.tokenList.Remove(e)
		} else if t.Token == token {
			// Token 存在且未超时，记录找到的 Token
			foundToken = t
		}
	}

	if foundToken != nil {
		fmt.Printf("Token found: token %s and public key %s\n", foundToken.Token, foundToken.PublicKey)
		return foundToken, nil
	}
	return nil, errors.New("token not found or expired")
}

func (tm *TokenManager) RemoveToken(token string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for e := tm.tokenList.Front(); e != nil; e = e.Next() {
		if e.Value.(*Token).Token == token {
			tm.tokenList.Remove(e)
			break
		}
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenManager := c.MustGet("tokenManager").(*TokenManager)

		token, err := c.Cookie("session_token")
		if err != nil {
			c.Set("publickey_hex", "")
			c.Next()
			return
		}

		// 检查token是否有效
		foundToken, err := tokenManager.CheckToken(token)
		if err != nil {
			c.Set("publickey_hex", "")
			c.Next()
			return
		}

		// 更新MatchedAt时间以防止超时
		foundToken.MatchedAt = time.Now()

		// 更新cookie过期时间
		c.SetCookie("session_token", token, 3600, "/", "", false, true)

		// 设置publickey_hex
		c.Set("publickey_hex", foundToken.PublicKey)

		c.Next()
	}
}
