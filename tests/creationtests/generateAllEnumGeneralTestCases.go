package creationtests

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/core-v9/codestack"
	"github.com/alimtvnetwork/core-v9/constants"
	"github.com/alimtvnetwork/core-v9/converters"
	"github.com/alimtvnetwork/core-v9/coredata/corestr"
	"github.com/alimtvnetwork/core-v9/coreimpl/enumimpl"
	"github.com/alimtvnetwork/core-v9/coreinterface/enuminf"
	"github.com/alimtvnetwork/enum-v9/quotes"
)

func generateAllEnumGeneralTestCases(
	isSpecificType bool,
	specificTypeNames ...string,
) {
	length := len(allBasicEnumsCollection)
	tab := constants.Tab
	tab2 := strings.Repeat(tab, 2)
	tab3 := strings.Repeat(tab, 3)
	// 	 tab4 := strings.Repeat(tab2, 2)
	// tab8 := strings.Repeat(tab4, 2)

	hashset := corestr.New.Hashset.Strings(specificTypeNames)
	fmt.Println("var allEnumGeneralTestCases = []*EnumTestWrapper{")
	failedBasicEnumer := corestr.New.SimpleSlice.Cap(10)
	doublesWrapFunc := quotes.Double.Wrap
	toFullStringer := converters.AnyTo.ToFullNameValueString

	for i := 0; i < length; i++ {
		item := allBasicEnumsCollection[i]
		typeName := item.TypeName()
		name := item.Name()

		if isSpecificType && hashset.IsMissing(typeName) {
			continue
		}

		fullInvalidName := codestack.NameOf.JoinPackageNameWithRelative(
			typeName,
			name)

		basicEnumer, isSuccess := item.(enuminf.BasicEnumer)

		if !isSuccess {
			fmt.Println("failed", item.TypeName())
			failedBasicEnumer.Add(item.TypeName())

			continue
		}

		fmt.Println(
			tab,
			"{")

		doubleQuoteWrappedTypeName := doublesWrapFunc(typeName)

		fmt.Println(
			tab2,
			"Header: \"Enum("+typeName+") ranges and values verification\",",
		)

		if typeName == "inttype.Variant" {
			fmt.Println(
				tab2,
				"InitialBasicEnumer: inttype.Zero.ToPtr(),",
			)
		} else if typeName == "strtype.Variant" {
			fmt.Println(
				tab2,
				"InitialBasicEnumer: strtype.Variant(\"Invalid\").ToPtr(),",
			)
		} else {
			fmt.Println(
				tab2,
				"InitialBasicEnumer:"+fullInvalidName+".ToPtr(),",
			)
		}

		fmt.Println(
			tab2,
			"TypeName:"+doubleQuoteWrappedTypeName+",",
		)

		enumType := basicEnumer.EnumType()

		fmt.Println(
			tab2,
			"ExpectedEnumType: enumtype."+enumType.Name()+",",
		)

		fmt.Println(
			tab2,
			"ExpectedMapValues: map[string]interface{} {",
		)

		var enumDynamicMap enumimpl.DynamicMap = basicEnumer.RangesDynamicMap()

		if enumType.IsNumber() {
			for _, keyValInteger := range enumDynamicMap.SortedKeyValues() {
				valueString := strconv.Itoa(keyValInteger.ValueInteger)

				fmt.Println(
					tab3,
					keyValInteger.WrapKey()+":",
					valueString,
					",",
				)
			}
		} else if enumType.IsString() {
			for _, anyKeyVal := range enumDynamicMap.SortedKeyAnyValues() {
				fmt.Println(
					tab3,
					toFullStringer(anyKeyVal.KeyString())+":",
					anyKeyVal.WrapValue(),
					",",
				)
			}
		}

		fmt.Println(
			tab2,
			"},", // ending of "ExpectedMapValues: map[string]interface{} {"
		)

		fmt.Println(
			tab2,
			"ExpectedInvalidName:"+doublesWrapFunc(basicEnumer.Name())+",",
		)

		fmt.Println(
			tab2,
			"ExpectedInvalidValueString:"+doublesWrapFunc(basicEnumer.ValueString())+",",
		)

		fmt.Println(
			tab2,
			"ExpectedRangesNamesCsv: "+toFullStringer(basicEnumer.RangeNamesCsv())+",",
		)

		if enumType.IsNumber() {
			fmt.Println(
				tab2,
				"IntegerMinMax: corerange.MinMaxInt {\n",
				tab3+"Min:"+basicEnumer.MinValueString()+",\n",
				tab3+"Max:"+basicEnumer.MaxValueString()+",",
			)

			fmt.Println(
				tab2,
				"},",
			)
		}

		fmt.Println(
			tab2,
			"StringMin:"+doublesWrapFunc(basicEnumer.MinValueString())+",",
		)

		fmt.Println(
			tab2,
			"StringMax:"+doublesWrapFunc(basicEnumer.MaxValueString())+",",
		)

		fmt.Println(
			tab,
			"},",
		)

		/*
			pathpatterntype.Root: {
					PathType:                 pathpatterntype.Root,
					Name:                     "Root",
					FullName:                 "root",
					CurlyFullName:            "{root}",
					CompiledTemplateFullPath: "{root}",
					AssociatedTemplatePaths: []string{
						"{root}",
					},
				}
		*/
	}

	fmt.Println(
		"}",
	)
}
