package utils

import (
	"reflect"
	"testing"
)

func TestProcessLinePairs(t *testing.T) {
	type args struct {
		filename       string
		allowLeftovers bool
	}
	tests := []struct {
		name      string
		args      args
		wantLines [][]string
		wantErr   bool
	}{
		{"data1.txt with leftovers", args{"testdata/data1.txt", true}, [][]string{[]string{"Line 1", "Line 2"}}, false},
		{"data1.txt without leftovers", args{"testdata/data1.txt", false}, [][]string{[]string{"Line 1", "Line 2"}}, true},

		{"data2.txt with leftovers", args{"testdata/data2.txt", true}, [][]string{[]string{"Line 1", "Line 2"}}, false},
		{"data2.txt without leftovers", args{"testdata/data2.txt", false}, [][]string{[]string{"Line 1", "Line 2"}}, false},

		{"non-existent file", args{"testdata/nonexistent.txt", false}, [][]string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLines := [][]string{}
			err := ProcessLinePairs(tt.args.filename, tt.args.allowLeftovers, func(l1 string, l2 string) error {
				gotLines = append(gotLines, []string{l1, l2})
				return nil
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessLinePairs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLines, tt.wantLines) {
				t.Errorf("ProcessLinePairs() lines = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}

func TestIterateFiles(t *testing.T) {
	type args struct {
		dir       string
		extension string
	}
	tests := []struct {
		name      string
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{"testdata .json files", args{"testdata", ".json"}, []string{"testdata/data3.json"}, false},
		{"testdata .txt files", args{"testdata", ".txt"}, []string{"testdata/data1.txt", "testdata/data2.txt", "testdata/nested/more.txt"}, false},
		{"testdata .go files", args{"testdata", ".go"}, []string{}, false},
		{"non-existent directory files", args{"testdata/nonexistent", ".txt"}, []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles := []string{}
			err := IterateFiles(tt.args.dir, tt.args.extension, func(path string) error {
				gotFiles = append(gotFiles, path)
				return nil
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("IterateFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("IterateFiles() files = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}

func TestHashFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		args     args
		wantHash string
		wantErr  bool
	}{
		{"hash file data1.txt", args{"testdata/data1.txt"}, "391ba54caa9e9da3dd31dca1eff275e706979e76c1f60c91401f0624734f52b0", false},
		{"hash file data2.txt", args{"testdata/data2.txt"}, "9140ddc651fb3861322111773bee1afd59db94a6dcba56212a5cabd8aaaf6874", false},
		{"hash file data3.json", args{"testdata/data3.json"}, "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a", false},
		{"hash non-existent file", args{"testdata/nonexistent.txt"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := HashFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHash != tt.wantHash {
				t.Errorf("HashFile() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}
