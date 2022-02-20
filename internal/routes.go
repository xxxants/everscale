package server

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

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

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func (s *Server) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		pong := Ping{Message: "pong"}
		c.JSON(http.StatusOK, pong)
	}
}

func (s *Server) GetGenome() gin.HandlerFunc {
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
		c.JSON(http.StatusOK, stdout)
	}
}

func (s *Server) newGenome() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var nft Nft
		// if err := c.BindJSON(&nft); err != nil {
		// 	s.Logger.Error(err)
		// 	respondWithError(c, 401, err.Error())
		// 	return
		// }

		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	return
		// }
		c.JSON(http.StatusOK, "stdout")
	}
}

func (s *Server) NftList() gin.HandlerFunc {
	return func(c *gin.Context) {
		graphqlClient := graphql.NewClient("https://net.ton.dev/graphql")
		graphqlRequest := graphql.NewRequest(`
			query { 
				accounts 
				(filter : {
					code_hash :{eq : "0000c410b6b21a716b351082b99226a7fd150802ef1e7f8760a60cfe2c0ac740"}
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
