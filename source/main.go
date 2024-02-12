package main

import (
	"context"
	"fmt"

	"github.com/LearningMotors/go-genproto/suki/pb/nms/note"
	enginePB "github.com/LearningMotors/go-genproto/suki/pb/nms/property_engine"
	"google.golang.org/protobuf/types/known/anypb"

	"benthos-testing/source/internal/propertyengine"
)

func main() {
	ctx := context.Background()

	err := propertyengine.Initialize()
	if err != nil {
		panic(err)
	}

	pe := propertyengine.Get()

	request := &enginePB.Request{
		Property: &note.Property{
			PropertyType: &note.Property_EntityType{
				EntityType: note.SukiEntityProperties_INVALID_ENTITY,
			},
			Value:    nil,
			Resolved: false,
		},
		Version:  note.Version_V1,
		Operand:  sectionS2AnyPB,
		Contexts: []*anypb.Any{cursorPositionAnyPB, sectionContentAnyPB},
	}

	res, err := pe.ResolveProperty(ctx, request)
	if err != nil {
		panic(err)
	}

	if len(res.Results) == 0 {
		panic("empty result from property engine")
	}

	fmt.Println("test vishnu: ", res.Results)
}
