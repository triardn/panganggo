package commons

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "doesnt comply min and max rule",
			args: args{
				s: "aBc1@",
			},
			want: false,
		},
		{
			name: "comply with all rules",
			args: args{
				s: "aBc1@se1kn7j",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePassword(tt.args.s); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparePasswords(t *testing.T) {
	type args struct {
		hashedPassword string
		plainPassword  []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				hashedPassword: "$2a$04$51FdsMsF1NCHjDP0VjApzO6o0Z.1Baf1nD5ua7CTO2Pmb0rTi2gs6",
				plainPassword:  []byte("narutoadalahhokage"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePasswords(tt.args.hashedPassword, tt.args.plainPassword); got != tt.want {
				t.Errorf("ComparePasswords() = %v, want %v", got, tt.want)
			}
		})
	}
}
