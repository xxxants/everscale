package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
)

type Ping struct {
	Message string `json:"message"`
}

type Nft struct {
	DataName string `json:"data_name"`
	Url      string `json:"url"`
	Genome   string `json:"genome"`
}

type Output struct {
	Message string `json:"message"`
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func (s *Server) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		pong := Ping{Message: "pong"}
		c.JSON(http.StatusOK, pong)
	}
}

func (s *Server) GetNftInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		app := "tonos-cli"
		arg0 := "run"
		arg1 := c.Query("addrNft")
		arg2 := "getInfo"
		arg3 := `'{"_answer_id":1}'`
		arg4 := "--abi"
		arg5 := "/root/testlift/contracts/src/compiled/Index.abi.json"
		cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		o := Output{
			Message: string(stdout),
		}
		c.JSON(http.StatusOK, o)
	}
}

func (s *Server) newGenome() gin.HandlerFunc {
	return func(c *gin.Context) {
		var nft Nft
		if err := c.BindJSON(&nft); err != nil {
			s.Logger.Error(err)
			respondWithError(c, 401, err.Error())
			return
		}

		jsdata, err := json.Marshal(nft)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		app := "tonos-cli"
		arg0 := "body"
		arg1 := "--abi"
		arg2 := "/root/testlift/contracts/src/NftRoot.abi.json"
		arg3 := "mintNft"
		arg4 := string(jsdata)

		cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		s.Logger.Debug(cmd.String())
		lastLine := strings.Split(string(stdout), "\n")[0]
		payload := strings.Split(lastLine, " ")[0]

		// tonos-cli call   --abi Wallet.abi.json submitTransaction

		app2 := "tonos-cli"
		arg01 := "call"
		arg11 := "0:41a0006aa2fcff5b91a63f510ee8baeab285229840c4cfa9e8bb9bc78378aba3"
		arg21 := "--abi"
		arg31 := "/root/testlift/contracts/src/Wallet.abi.json"
		arg41 := "submitTransaction"
		arg51 := `'{"dest":"0:73ea2df343d928e9aa4c715cf15dec4f5f90193de714cdfedcae32a563d91347","value":2000000000,"bounce":true,"allBalance":false,"payload":"` + payload + `"`
		arg61 := `--sign $(cat /root/keys)`

		cmd2 := exec.Command(app2, arg01, arg11, arg21, arg31, arg41, arg51, arg61)
		stdout2, err := cmd2.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		s.Logger.Debug(cmd2.String())
		lastLine2 := strings.Split(string(stdout2), "\n")[0]
		payload2 := strings.Split(lastLine2, " ")[0]

		o := Output{
			Message: payload2,
		}
		c.JSON(http.StatusOK, o)
	}
}

func (s *Server) NftList() gin.HandlerFunc {
	return func(c *gin.Context) {
		graphqlClient := graphql.NewClient("https://net.ton.dev/graphql")
		graphqlRequest := graphql.NewRequest(`
			query { 
				accounts 
				(filter : {
					code_hash :{eq : "3d3addc0068a703236d5a6c95d7d74dbe4dda27c1a3893f056f8d8ffcdde84c8"}
				})
			{
				id
			}}
		`)
		var graphqlResponse interface{}
		if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
			s.Logger.Error(err)
			return
		}

		c.JSON(http.StatusOK, graphqlResponse)
	}
}
