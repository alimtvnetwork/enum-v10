package creationtests

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/core-v9/constants"
	"github.com/alimtvnetwork/core-v9/corecsv"
	"github.com/alimtvnetwork/enum-v10/pathpatterntype"
)

func generatePathPatternTestCases() {
	maxValue := int(pathpatterntype.BasicEnumImpl.Max())
	tab := constants.Tab
	tab2 := strings.Repeat(tab, 2)
	tab3 := strings.Repeat(tab, 3)
	// 	 tab4 := strings.Repeat(tab2, 2)
	// tab8 := strings.Repeat(tab4, 2)

	fmt.Println("var pathPatternTypeCreationTestCases = [...]PathPatternTypeCreationTestWrapper{")

	for i := 0; i <= maxValue; i++ {
		pathType := pathpatterntype.Variant(i)
		fmt.Println(
			tab,
			"pathpatterntype."+pathType.Name()+": {")

		fmt.Println(
			tab2,
			"PathType:                 pathpatterntype."+pathType.Name(),
			",",
		)

		fmt.Println(
			tab2,
			"Name:\""+pathType.Name()+"\",",
		)
		fmt.Println(
			tab2,
			"FullName:\""+pathType.PathFullName()+"\",",
		)

		fmt.Println(
			tab2,
			"CurlyFullName:\""+pathType.CurlyPathFullName()+"\",",
		)

		fmt.Println(
			tab2,
			"CompiledTemplateFullPath:`"+pathType.CompileCurlyTemplate()+"`,",
		)

		fmt.Println(
			tab2,
			"AssociatedTemplatePaths:[]string{",
		)

		fmt.Println(
			tab3,
			corecsv.DefaultCsvUsingJoiner(
				", "+constants.DefaultLine+tab3,
				pathType.SplitExpandedAssocCurlyPathStrings()...)+",",
		)

		fmt.Println(
			tab2,
			"},",
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
