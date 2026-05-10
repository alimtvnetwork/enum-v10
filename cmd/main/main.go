package main

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
	
	"github.com/alimtvnetwork/core-v9/coredata/coredynamic"
	"github.com/alimtvnetwork/core-v9/coredata/corejson"
	"github.com/alimtvnetwork/core-v9/iserror"
	"github.com/alimtvnetwork/core-v9/issetter"
	"github.com/alimtvnetwork/enum-v10/brackets"
	"github.com/alimtvnetwork/enum-v10/dbaction"
	"github.com/alimtvnetwork/enum-v10/instructiontype"
	"github.com/alimtvnetwork/enum-v10/osdetect"
	"github.com/alimtvnetwork/enum-v10/strtype"
)

func main() {
	stringTypeTest()
	bracketsTest()
	
	// Windwos_Version()
	osDetailsTest()
	// serializeDeserializeTester()
	err := os.Setenv("myname", "alim")
	
	fmt.Println(iserror.AllDefined(nil, errors.New("x")))
	
	fmt.Println("myname, err:", err)
}

func stringTypeTest() {
	fmt.Println(instructiontype.New("DependsOnx"))
	fmt.Println(dbaction.New("Create"))
	
	alimStrType := strtype.New("alimx")
	
	fmt.Println(alimStrType.SafeSubString(0, 1))
	fmt.Println(alimStrType.SafeSubString(0, alimStrType.Length()))
	fmt.Println(alimStrType.SafeSubString(-1, alimStrType.Length()+5))
	fmt.Println(alimStrType.SafeSubStringStart(2))
	fmt.Println(alimStrType.SafeSubStringEnd(2))
	fmt.Println(alimStrType.SafeSplit(2))
	fmt.Println(alimStrType.SafeSplit(alimStrType.Length() + 5))
}

func bracketsTest() {
	bracket := brackets.Parenthesis
	
	fmt.Println(bracket.WrapAny("something"))
	
	bracket2 := brackets.ParenthesisStart
	fmt.Printf("sizeof(bracket2) = %d\n", unsafe.Sizeof(bracket2))
	
	someTrue := true
	Val := issetter.True
	
	fmt.Printf("sizeof(someTrue) = %d\n", unsafe.Sizeof(&someTrue))
	fmt.Printf("sizeof(Val) = %d\n", unsafe.Sizeof(&Val))
	
	fmt.Println(bracket2.WrapAny("something2"))
	fmt.Println(bracket2.WrapFmtString("something to do with {wrapped}", "something2"))
	fmt.Println(bracket2.WrapSkipOnExist("(something2)"))
	fmt.Println(bracket2.IsWrapped("(something2)"))
}

func osDetailsTest() {
	osDetail, err := osdetect.GetCurrentOsDetail()
	
	fmt.Println(osDetail.PrettyJsonString())
	fmt.Println("err", err)
	
	fmt.Println("current os mix types", osdetect.CurrentOsType())
	fmt.Println("all os mix types", osdetect.CurrentOsMixTypes())
	fmt.Println("all os mix map", osdetect.CurrentOsTypesMap())
	fmt.Println("osdetect.Ubuntu.IsMajorAtLeast(18) : ", osdetect.Ubuntu.IsMajorAtLeast(18))
	fmt.Println("osdetect.Ubuntu.IsMajorAtLeast(20) : ", osdetect.Ubuntu.IsMajorAtLeast(20))
	fmt.Println("osdetect.Ubuntu.IsMajorAtLeast(21) : ", osdetect.MacOs.IsMajorAtLeast(21))
	fmt.Println("osdetect.Windows.IsMajorAtLeast(10) : ", osdetect.Windows.IsMajorAtLeast(10))
	fmt.Println("osdetect.Windows.IsWindows11() : ", osdetect.Windows.IsWindows11())
	fmt.Println("osdetect.Windows.Name(): ", osdetect.Windows.Name())
	fmt.Println("osdetect.Windows.ProductName(): ", osdetect.Windows.ProductName())
	fmt.Println("osdetect.Windows.RawProductName(): ", osdetect.Windows.RawProductName())
}

func serializeDeserializeTester() {
	osDetail, err := osdetect.GetCurrentOsDetail()
	
	fmt.Println(err)
	
	slice := []interface{}{
		osDetail,
		osDetail.WindowsDetail,
	}
	
	for _, item := range slice {
		serializeDeserializeTesterByInput(item)
		// t := reflect.TypeOf(item)
		// for i := 0; i < t.NumField(); i++ {
		// 	fmt.Printf("%+v\n", t.Field(i))
		// 	fmt.Printf("%+v\n", t.(i))
		// }
	}
	
}

func serializeDeserializeTesterByInput(input interface{}) {
	json := corejson.NewPtr(input)
	
	if json.HasIssuesOrEmpty() {
		fmt.Println(coredynamic.SafeTypeName(input))
		fmt.Println("Marshalling Err", json.MeaningfulErrorMessage())
	}
	
	finalErr := json.Deserialize(input)
	
	if finalErr != nil {
		fmt.Println(finalErr)
		fmt.Println(coredynamic.SafeTypeName(input))
		fmt.Println("Json", json.PrettyJsonString())
	}
}
