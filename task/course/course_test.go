package course

import "testing"

func Test_updateCourses(t *testing.T) {
	type args struct {
		studentID int
		password  string
		pageNo    int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				studentID: 123456789,
				password:  "fghjk",
				pageNo:    10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCourses(tt.args.studentID, tt.args.password, tt.args.pageNo)
		})
	}
}
