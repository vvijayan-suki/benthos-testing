package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/LearningMotors/go-genproto/learningmotors/pb/composer"
	"github.com/LearningMotors/go-genproto/suki/pb/nms/dynamic_data"
	"github.com/LearningMotors/go-genproto/suki/pb/nms/note"
	enginePB "github.com/LearningMotors/go-genproto/suki/pb/nms/property_engine"
	"github.com/LearningMotors/go-genproto/suki/pb/sectioncontent"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"benthos-testing/source/internal/propertyengine"
)

func main() {
	runDynamicLabData()
}

func runDynamicLabData() {
	ctx := context.Background()

	err := propertyengine.Initialize()
	if err != nil {
		panic(err)
	}

	pe := propertyengine.Get()

	sectionS2 := getSectionS2()

	vitalsMapping := getMapping()
	vitalsMappingAnyPB, err := anypb.New(vitalsMapping)
	if err != nil {
		panic(err)
	}

	re, err := regexp.Compile("__[A-Za-z]+_*[A-Za-z]+__")
	if err != nil {
		panic(err)
	}

	plainText := sectionS2.GetPlainText()
	reStrings := re.FindAllString(plainText, -1)

	for _, reString := range reStrings {
		sectionS2AnyPB, err := anypb.New(sectionS2)
		if err != nil {
			panic(err)
		}

		reStringPB := wrapperspb.String(reString)
		reStringAnyPB, err := anypb.New(reStringPB)
		if err != nil {
			panic(err)
		}

		abbreviationsPB, err := structpb.NewStruct(abbreviations)
		if err != nil {
			panic(err)
		}
		abbreviationsAnyPB, err := anypb.New(abbreviationsPB)
		if err != nil {
			panic(err)
		}

		request := &enginePB.Request{
			Property: &note.Property{
				PropertyType: &note.Property_LabResultType{
					LabResultType: note.LabResultProperties_LAB_RESULT_VALUE,
				},
				Value:    nil,
				Resolved: false,
			},
			Version:  note.Version_V1,
			Operand:  sectionS2AnyPB,
			Contexts: []*anypb.Any{vitalsMappingAnyPB, reStringAnyPB, abbreviationsAnyPB},
		}

		fmt.Println("request: ", request)

		res, err := pe.ResolveProperty(ctx, request)
		if err != nil {
			panic(err)
		}

		if len(res.Results) == 0 {
			panic("empty result from property engine")
		}

		// first result is the section content
		firstRes := res.Results[0]

		if firstRes == nil {
			panic("got nil result from property engine")
		}

		switch firstRes.TypeUrl {
		case fmt.Sprintf("type.googleapis.com/%s", proto.MessageName(sectionS2)):
			fmt.Println("here fool")
			resultantEntry := &composer.SectionS2{}
			if err := firstRes.UnmarshalTo(resultantEntry); err != nil {
				panic(err)
			}

			sectionS2 = resultantEntry
		}

		logrus.Info("vishnu: ", res)
	}
	fmt.Println("done")

	fmt.Println("final: ", sectionS2)
}

func getSectionS2() *composer.SectionS2 {
	return &composer.SectionS2{
		Id:                 "section_id_1",
		Name:               "",
		NavigationKeywords: nil,
		ContentS2: &sectioncontent.SectionContent{
			NumberOfStrings:   0,
			TotalStringLength: 64,
			TotalString:       "",
			Content: []*sectioncontent.Content{
				&sectioncontent.Content{
					Id:               0,
					StartOffset:      0,
					EndOffset:        23,
					String_:          "bp: __blood_pressure__  ",
					LengthOfString:   24,
					IsBold:           0,
					IsItalic:         0,
					Source:           sectioncontent.TextSource_DYNAMIC_VITALS,
					RecommendationId: "vishnu_test",
				},
				&sectioncontent.Content{
					Id:               1,
					StartOffset:      24,
					EndOffset:        40,
					String_:          "hr: __heartrate__",
					LengthOfString:   17,
					IsBold:           0,
					IsItalic:         0,
					Source:           0,
					RecommendationId: "",
				},
				&sectioncontent.Content{
					Id:               2,
					StartOffset:      41,
					EndOffset:        51,
					String_:          "bp: __bp__ ",
					LengthOfString:   11,
					IsBold:           0,
					IsItalic:         0,
					Source:           0,
					RecommendationId: "",
				},
			},
		},
		Status:               0,
		Cursors:              nil,
		Hash:                 "",
		DiagnosisEntry:       nil,
		PbcSectionFlag:       false,
		PlainText:            "bp: __blood_pressure__  hr: __heartrate__bp: __bp__ ",
		CursorPosition:       0,
		SectionIndex:         0,
		SubsectionIndex:      0,
		CursorPositionName:   0,
		EditLocation:         0,
		UpdateType:           0,
		OpsStatusFlag:        0,
		NumberOfCursorEvents: 0,
		CursorEndIndex:       0,
		Footer:               nil,
		DictationTag:         "",
		ReadOnly:             false,
		Display:              nil,
	}
}

func getMapping() *dynamic_data.DynamicData {
	vc := &composer.VersionedComposition{
		DynamicData: &dynamic_data.DynamicData{
			Mapping: make(map[string]*dynamic_data.DynamicChartData),
		},
	}

	vc.DynamicData.Mapping["blood_pressure"] = &dynamic_data.DynamicChartData{
		Type:            0,
		Content:         "120/80",
		ResultedDate:    timestamppb.Now(),
		ShouldHaveDates: true,
	}

	vc.DynamicData.Mapping["heartrate"] = &dynamic_data.DynamicChartData{
		Type:          0,
		Content:       "80",
		EffectiveDate: timestamppb.Now(),
		//ShouldHaveDates: true,
	}

	vc.DynamicData.Mapping["respiratory_rate"] = &dynamic_data.DynamicChartData{
		Type:    0,
		Content: "16",
	}

	return vc.GetDynamicData()
}

var abbreviations = map[string]interface{}{
	//"bp":   "blood_pressure",
	"hr":   "heartrate",
	"temp": "temperature",
	"rr":   "respiratory_rate",
	"po":   "pulse_oximetry",
	"ht":   "height",
	"wt":   "weight",
	"bmi":  "bmi",
}
