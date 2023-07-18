package humanize

import (
	"testing"
)

func TestBytes(t *testing.T) {
	type args struct {
		size uint64
	}
	tests := []struct {
		name             string
		args             args
		wantHumanizeSize string
	}{
		{
			name: "test 0",
			args: args{
				size: 0,
			},
			wantHumanizeSize: "0 B",
		},
		{
			name: "test 1",
			args: args{
				size: 1,
			},
			wantHumanizeSize: "1 B",
		},
		{
			name: "test 888",
			args: args{
				size: 888,
			},
			wantHumanizeSize: "888 B",
		},
		{
			name: "test 1024",
			args: args{
				size: 1024,
			},
			wantHumanizeSize: "1.0 kB",
		},
		{
			name: "test 1000*1000*1000",
			args: args{
				size: 1000 * 1000 * 1000,
			},
			wantHumanizeSize: "1.0 GB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHumanizeSize := Bytes(tt.args.size); gotHumanizeSize != tt.wantHumanizeSize {
				t.Errorf("Bytes() = %v, want %v", gotHumanizeSize, tt.wantHumanizeSize)
			}
		})
	}
}

func TestIBytes(t *testing.T) {
	type args struct {
		size uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 0",
			args: args{
				size: 0,
			},
			want: "0 B",
		},
		{
			name: "test 1",
			args: args{
				size: 1,
			},
			want: "1 B",
		},
		{
			name: "test 888",
			args: args{
				size: 888,
			},
			want: "888 B",
		},
		{
			name: "test 1000",
			args: args{
				size: 1000,
			},
			want: "1000 B",
		},
		{
			name: "test 1024",
			args: args{
				size: 1024,
			},
			want: "1.0 KiB",
		},
		{
			name: "test 1000*1000*1000",
			args: args{
				size: 1000 * 1000 * 1000,
			},
			want: "954 MiB",
		},
		{
			name: "test 1024*1024*1024",
			args: args{
				size: 1024 * 1024 * 1024,
			},
			want: "1.0 GiB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IBytes(tt.args.size); got != tt.want {
				t.Errorf("IBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bytes(1000000000000)
	}
}

func BenchmarkIBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IBytes(1000000000000)
	}
}
