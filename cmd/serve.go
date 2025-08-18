/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"time"

	_ "github.com/arwoosa/media/internal/db"
	_ "github.com/arwoosa/media/internal/service"
	"github.com/arwoosa/vulpes/codec"
	"github.com/arwoosa/vulpes/db/cache"
	"github.com/arwoosa/vulpes/db/mgo"
	"github.com/arwoosa/vulpes/ezgrpc"
	"github.com/arwoosa/vulpes/log"
	"github.com/arwoosa/vulpes/relation"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the media service server",
	Long: `Start the media service server with REST API and gRPC endpoints.

The server can be started with either REST API or gRPC endpoints:
  - api: Start the REST API server
  - grpc: Start the gRPC server
`,
	Run: func(cmd *cobra.Command, args []string) {
		// initialize cache
		err := cache.InitConnection(
			cache.WithAddr(viper.GetString("cache.address")),
			cache.WithPassword(viper.GetString("cache.password")),
			cache.WithDb(viper.GetInt("cache.db")))
		if err != nil {
			log.Fatal(err.Error())
		}
		// initialize relation
		relation.Initialize(
			relation.WithWriteAddr(viper.GetString("relation.write_uri")),
			relation.WithReadAddr(viper.GetString("relation.read_uri")))

		// initialize mongo
		mongoCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		err = mgo.InitConnection(mongoCtx,
			viper.GetString("database.db"),
			mgo.WithURI(viper.GetString("database.uri")),
			mgo.WithMaxPoolSize(viper.GetUint64("database.max_pool_size")),
			mgo.WithMinPoolSize(viper.GetUint64("database.min_pool_size")),
		)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer func() {
			closeCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			err = mgo.Close(closeCtx)
			if err != nil {
				log.Fatal(err.Error())
			}
		}()

		// initialize mongo indexes
		err = mgo.SyncIndexes(mongoCtx)
		if err != nil {
			log.Fatal(err.Error())
		}
		cancel()

		ezgrpc.InitSessionStore()
		// interceptor.DisableValidateInterceptor()

		codec.WithCodecMethod(codec.GOB)
		ezgrpc.SetServeMuxOpts(
			ezgrpc.DefaultHeaderMatcher,
			ezgrpc.OutgoingHeaderMatcher,
			ezgrpc.SessionCookieForwarder,
			ezgrpc.SessionCookieExtractor,
			ezgrpc.RedirectResponseOption,
			ezgrpc.RedirectResponseModifier,
		)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = ezgrpc.RunGrpcGateway(ctx, viper.GetInt("server.port"))
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
