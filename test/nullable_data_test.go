package test

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNullableDataFromPtr(t *testing.T) {
	type args struct {
		value *Data
	}
	tests := []struct {
		name string
		args args
		want NullableData
	}{
		{
			name: "nilが引数ならinvalidなnullableが生成ができる",
			args: args{
				value: nil,
			},
			want: NullableData{
				value: nil,
				valid: false,
			},
		},
		{
			name: "引数がnilでないならvalidなnullableが生成ができる",
			args: args{
				value: &Data{
					A: "a",
					B: 1,
				},
			},
			want: NullableData{
				value: &Data{
					A: "a",
					B: 1,
				},
				valid: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullableDataFromPtr(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullableDataFromPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableDataFrom(t *testing.T) {
	type args struct {
		value Data
	}
	tests := []struct {
		name string
		args args
		want NullableData
	}{
		{
			name: "必ずvalidなデータが生成される",
			args: args{
				value: Data{
					A: "a",
					B: 1,
				},
			},
			want: NullableData{
				value: &Data{
					A: "a",
					B: 1,
				},
				valid: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullableDataFrom(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullableDataFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableData_Ptr(t *testing.T) {
	type fields struct {
		value *Data
		valid bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *Data
	}{
		{
			name: "validがtrueなら値のポインタが返る",
			fields: fields{
				value: &Data{
					A: "1",
					B: 2,
				},
				valid: true,
			},
			want: &Data{
				A: "1",
				B: 2,
			},
		},
		{
			name: "validがfalseならnilが返る",
			fields: fields{
				value: nil,
				valid: false,
			},
			want: nil,
		},
		{
			name: "validがfalseならnilが返るたとえvalueに値がいても",
			fields: fields{
				value: &Data{
					A: "1",
					B: 2,
				},
				valid: false,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NullableData{
				value: tt.fields.value,
				valid: tt.fields.valid,
			}
			if got := v.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

var sampleDataJSONBytes = []byte(`{"a":"1","b":2}`)

func TestNullableData_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    NullableData
		wantErr bool
	}{
		{
			name:    "空バイトの場合nilになる",
			args:    emptyDataBytes,
			wantErr: false,
			want: NullableData{
				value: nil,
				valid: false,
			},
		},
		{
			name:    "`\"\"`の場合nilになる",
			args:    emptyDataJSON,
			wantErr: false,
			want: NullableData{
				value: nil,
				valid: false,
			},
		},
		{
			name:    "`null`の場合nilになる",
			args:    nullDataBytes,
			wantErr: false,
			want: NullableData{
				value: nil,
				valid: false,
			},
		},
		{
			name:    "jsonのbyte sliceから意図した通りに復元できる",
			args:    sampleDataJSONBytes,
			wantErr: false,
			want: NullableData{
				value: &Data{
					A: "1",
					B: 2,
				},
				valid: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &NullableData{}
			err := v.UnmarshalJSON(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && tt.want.value == nil && v.value == nil {
				return
			}

			if err == nil && (tt.want.valid != v.valid ||
				*tt.want.value != *v.value) {
				t.Errorf("unexpected value. actual = %+v, want %+v", *v.value, *tt.want.value)
			}
		})
	}
}

func TestNullableData_MarshalJSON(t *testing.T) {
	type fields struct {
		value *Data
		valid bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "元structで定義したのと同じようにJSONになる",
			fields: fields{
				value: &Data{
					A: "1",
					B: 2,
				},
				valid: false,
			},
			want:    sampleDataJSONBytes,
			wantErr: false,
		},
		{
			name: "valueがnilならnilのJSONが返る",
			fields: fields{
				value: nil,
				valid: false,
			},
			want:    nullDataBytes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NullableData{
				value: tt.fields.value,
				valid: tt.fields.valid,
			}
			got, err := v.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}
