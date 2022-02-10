// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package fastrest

import (
	json "encoding/json"

	easyjson "github.com/bingoohuang/easyjson"
	jlexer "github.com/bingoohuang/easyjson/jlexer"
	jwriter "github.com/bingoohuang/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC80ae7adDecodeGithubComBingoohuangFastrest(in *jlexer.Lexer, out *Rsp) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "status":
			out.Status = int(in.Int())
		case "message":
			out.Message = string(in.String())
		case "data":
			if m, ok := out.Data.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Data.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Data = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjsonC80ae7adEncodeGithubComBingoohuangFastrest(out *jwriter.Writer, in Rsp) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Status != 0 {
		const prefix string = ",\"status\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.Status))
	}
	if in.Message != "" {
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Message))
	}
	if in.Data != nil {
		const prefix string = ",\"data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if m, ok := in.Data.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Data.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Data))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Rsp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Rsp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Rsp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Rsp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest(l, v)
}

func easyjsonC80ae7adDecodeGithubComBingoohuangFastrest1(in *jlexer.Lexer, out *P1SignRsp) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "source":
			out.Source = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjsonC80ae7adEncodeGithubComBingoohuangFastrest1(out *jwriter.Writer, in P1SignRsp) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"source\":"
		out.RawString(prefix[1:])
		out.String(string(in.Source))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v P1SignRsp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v P1SignRsp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *P1SignRsp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *P1SignRsp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest1(l, v)
}

func easyjsonC80ae7adDecodeGithubComBingoohuangFastrest2(in *jlexer.Lexer, out *P1SignReq) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "source":
			out.Source = string(in.String())
		case "bizType":
			out.BizType = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjsonC80ae7adEncodeGithubComBingoohuangFastrest2(out *jwriter.Writer, in P1SignReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"source\":"
		out.RawString(prefix[1:])
		out.String(string(in.Source))
	}
	{
		const prefix string = ",\"bizType\":"
		out.RawString(prefix)
		out.String(string(in.BizType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v P1SignReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v P1SignReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *P1SignReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *P1SignReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest2(l, v)
}
