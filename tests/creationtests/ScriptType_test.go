package creationtests

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/alimtvnetwork/enum-v9/scripttype"
)

// Test_AllEnums_ContractsTesting
//
//	Generates testcases by generateAllEnumGeneralTestCases()
//	Runs by allEnumGeneralTestCases
func Test_ScriptType(t *testing.T) {
	// shouldBe := errcore.ShouldBe.StrEqMsg

	for scriptCreationName, expectedScriptType := range allScriptCreationTestCases {
		// Arrange
		createdScriptType, err := scripttype.New(scriptCreationName)

		Convey(scriptCreationName+" should be created properly", t, func() {
			So(err, ShouldBeNil)
			So(createdScriptType, ShouldEqual, expectedScriptType)
		})
	}
}
