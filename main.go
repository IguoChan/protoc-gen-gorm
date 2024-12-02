package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"

	pgg "github.com/IguoChan/protoc-gen-gorm/internal/protoc-gen-gorm"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.0.0"

var (
	showVersion    = flag.Bool("version", false, "print the version and exit")
	withGormOption = flag.Bool("with_gorm_option", false, "with gorm option")
	withGormDao    = flag.Bool("with_gorm_dao", false, "with gorm dao")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-gorm %v\n", version)
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) | uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
		gen.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO2
		gen.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2023
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			if *withGormDao && !*withGormOption {
				fmt.Println("with_dao_option must be used with with_gorm_option")
				return nil
			}

			opts := []pgg.Option{
				pgg.WithVersion(version),
			}
			if *withGormOption {
				opts = append(opts, pgg.WithGormOption())
			}

			gg := pgg.New(gen, f, opts...)
			gg.GenerateFile()

			if *withGormOption {
				optGG := pgg.NewGormOptionGenerator(gen, f)
				optGG.GenerateFile()
			}

			if *withGormDao {
				daoGG := pgg.NewGormDaoGenerator(gen, f)
				daoGG.GenerateFile()
			}
		}
		return nil
	})
}
