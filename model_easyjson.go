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
		case "data":
			if m, ok := out.Data.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Data.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Data = in.Interface()
			}
		case "message":
			out.Message = string(in.String())
		case "status":
			out.Status = int(in.Int())
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
	if in.Data != nil {
		const prefix string = ",\"data\":"
		first = false
		out.RawString(prefix[1:])
		if m, ok := in.Data.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Data.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Data))
		}
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
	if in.Status != 0 {
		const prefix string = ",\"status\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Status))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Rsp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest(w, v)
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

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v P1SignRsp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest1(w, v)
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

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v P1SignReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *P1SignReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest2(l, v)
}
func easyjsonC80ae7adDecodeGithubComBingoohuangFastrest3(in *jlexer.Lexer, out *EncryptRsp) {
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
		case "data":
			out.Data = string(in.String())
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
func easyjsonC80ae7adEncodeGithubComBingoohuangFastrest3(out *jwriter.Writer, in EncryptRsp) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		out.String(string(in.Data))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EncryptRsp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest3(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EncryptRsp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest3(l, v)
}
func easyjsonC80ae7adDecodeGithubComBingoohuangFastrest4(in *jlexer.Lexer, out *EncryptReq) {
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
		case "transId":
			out.TransId = string(in.String())
		case "appId":
			out.AppId = string(in.String())
		case "keyId":
			out.KeyId = string(in.String())
		case "mode":
			out.Mode = string(in.String())
		case "padding":
			out.Padding = string(in.String())
		case "plainText":
			out.PlainText = string(in.String())
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
func easyjsonC80ae7adEncodeGithubComBingoohuangFastrest4(out *jwriter.Writer, in EncryptReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"transId\":"
		out.RawString(prefix[1:])
		out.String(string(in.TransId))
	}
	{
		const prefix string = ",\"appId\":"
		out.RawString(prefix)
		out.String(string(in.AppId))
	}
	{
		const prefix string = ",\"keyId\":"
		out.RawString(prefix)
		out.String(string(in.KeyId))
	}
	{
		const prefix string = ",\"mode\":"
		out.RawString(prefix)
		out.String(string(in.Mode))
	}
	{
		const prefix string = ",\"padding\":"
		out.RawString(prefix)
		out.String(string(in.Padding))
	}
	{
		const prefix string = ",\"plainText\":"
		out.RawString(prefix)
		out.String(string(in.PlainText))
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EncryptReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComBingoohuangFastrest4(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EncryptReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComBingoohuangFastrest4(l, v)
}
