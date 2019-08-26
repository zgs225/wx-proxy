/*
Copyright © 2019 yuez i@yuez.me

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/zgs225/wxproxy/wx"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动服务",
	Long:  "启动一个微信公众号平台 OAuth2 登录代理",
	Run: func(cmd *cobra.Command, args []string) {
		h := wx.NewWXProxyHTTPServer()
		c := make(chan error)

		go func() {
			c <- http.ListenAndServe(":8080", h.NewHandler())
		}()

		log.Println("微信代理服务器运行在 :8080")
		log.Panic(<-c)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
