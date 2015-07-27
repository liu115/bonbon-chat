package communicate

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitHandler 用於初始化的handler
func InitHandler(c *gin.Context) {
	fake := initMsg{
		Setting: setting{"你每天都那麼努力，忍受那麼多的寂寞和痛苦，但你為什麼還那麼魯。"},
		Friends: []friend{
			friend{ID: 1, Nick: "Apple", Sign: "", Msg: []string{}},
			friend{ID: 2, Nick: "Banana", Sign: "", Msg: []string{"我要跟你分手"}},
			friend{ID: 3, Nick: "Cake", Sign: "", Msg: []string{}},
			friend{ID: 4, Nick: "Dog", Sign: "", Msg: []string{}},
			friend{ID: 5, Nick: "Egg", Sign: "", Msg: []string{"已經結束了"}}}}
	str, _ := json.Marshal(fake)
	fmt.Println(string(str))
	c.JSON(http.StatusOK, fake)
}
