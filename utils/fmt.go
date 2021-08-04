package utils

import (
	"io"
	"net/http"
	"net/http/httputil"
)

func FmtRequest(writer io.Writer, header string, req *http.Request, rsp http.ResponseWriter) error {
	reqDump, err := httputil.DumpRequest(req, true)
	//rspDump, err := httputil.DumpResponse(&rsp, true)

	if err != nil {
		return err
	}
	writer.Write([]byte("\n-----------------------start--------------------------- \n"))
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write([]byte("\n Request: \n"))
	writer.Write(reqDump)
	//writer.Write([]byte("\n Response: \n"))
	//writer.Write(rspDump)
	writer.Write([]byte("\n------------------------end--------------------------- \n"))
	return nil
}
