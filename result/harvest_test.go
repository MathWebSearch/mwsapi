package result

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestHarvestElementUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		xmlAsset    string
		jsonAsset   string
		wantXMLErr  bool
		wantJSONErr bool
	}{
		{"test1", "testdata/test1.xml", "testdata/test1.json", false, false},
		{"test2", "testdata/test2.xml", "testdata/test2.json", false, false},
		{"test3", "testdata/test3.xml", "testdata/test3.json", false, false},
		{"test4", "testdata/test4.xml", "testdata/test4.json", false, false},
		{"test5", "testdata/test5.xml", "testdata/test5.json", false, false},
		{"test6", "testdata/test6.xml", "testdata/test6.json", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read the bytes of the xml
			xmlBytes, err := ioutil.ReadFile(tt.xmlAsset)
			if err != nil {
				t.Errorf("Unable to load xml asset: %s", err.Error())
				return
			}

			// parse it
			var gotXML HarvestElement
			err = xml.Unmarshal([]byte(xmlBytes), &gotXML)
			if (err != nil) != tt.wantXMLErr {
				t.Errorf("HarvestElement.UnmarshalXML() error = %v, wantXMLErr %v", err, tt.wantXMLErr)
				return
			}

			// Manual testing; todo: remove me
			// bytes, _ := json.Marshal(gotXML)
			// fmt.Println(string(bytes))

			// read the bytes of the json
			jsonBytes, err := ioutil.ReadFile(tt.jsonAsset)
			if err != nil {
				t.Errorf("Unable to load json asset: %s", err.Error())
				return
			}

			// parse it
			var gotJSON HarvestElement
			err = json.Unmarshal([]byte(jsonBytes), &gotJSON)
			if (err != nil) != tt.wantJSONErr {
				t.Errorf("HarvestElement.UnmarshalJSON() error = %v, wantJSONErr %v", err, tt.wantJSONErr)
				return
			}

			// we do not have Numbers and Paths in the xml
			// so we will update them manually
			gotXML.MWSNumbers = gotJSON.MWSNumbers
			gotXML.MWSPaths = gotJSON.MWSPaths

			// compare them
			if !reflect.DeepEqual(&gotXML, &gotJSON) {
				t.Errorf("HarvestElement.UnmarshalXML() = %v\nHarvestElement.UnmarshalJSON() = %v\nNot identical", gotXML, gotJSON)
				return
			}

		})
	}
}
