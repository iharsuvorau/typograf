package typograf

import "testing"
import "strings"

func TestTypografy(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{"A", args{"у \"окна\" хорошо, а на диване - лучше"}, "<p>у&nbsp;&laquo;окна&raquo; хорошо, а&nbsp;на&nbsp;диване&nbsp;&#151; лучше<br /></p>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := Typografy(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Typografy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// stripping \n from the service output
			gotOut = strings.Replace(gotOut, "\n", "", -1)
			if gotOut != tt.wantOut {
				t.Errorf("Typografy() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
