package creationtests

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/alimtvnetwork/enum-v9/pathpatterntype"
)

func Test_PathType_Creation(t *testing.T) {
	maxValue := int(pathpatterntype.BasicEnumImpl.Max())

	for i := 0; i <= maxValue; i++ {
		pathType := pathpatterntype.Variant(i)
		name := pathType.Name()
		testCase := pathPatternTypeCreationTestCases[i]
		testCaseName := testCase.Name

		Convey("Test case equal to PathType Name", t, func() {
			So(name, ShouldEqual, testCaseName)
		})

		Convey("Test case equal to PathType Value", t, func() {
			So(testCase.PathType.Value(), ShouldEqual, pathType.Value())
		})

		joinedAssocPath :=
			strings.Join(
				testCase.AssociatedTemplatePaths,
				"\\")

		Convey("Test case equal to Associated path compiled", t, func() {
			So(joinedAssocPath, ShouldEqual, testCase.CompiledTemplateFullPath)
		})
	}
}
