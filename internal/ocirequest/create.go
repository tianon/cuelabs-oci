package ocirequest

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

func (req *Request) Construct() (method string, ustr string) {
	method, ustr = req.construct()
	u, err := url.Parse(ustr)
	if err != nil {
		panic(err)
	}
	if _, err := Parse(method, u); err != nil {
		panic(fmt.Errorf("invalid request %q %q constructed from %#v: %v", method, ustr, req, err))
	}
	return method, ustr
}

func (req *Request) construct() (method string, url string) {
	switch req.Kind {
	case ReqPing:
		return "GET", "/v2/"
	case ReqBlobGet:
		return "GET", "/v2/" + req.Repo + "/blobs/" + req.Digest
	case ReqBlobHead:
		return "HEAD", "/v2/" + req.Repo + "/blobs/" + req.Digest
	case ReqBlobDelete:
		return "DELETE", "/v2/" + req.Repo + "/blobs/" + req.Digest
	case ReqBlobStartUpload:
		return "POST", "/v2/" + req.Repo + "/blobs/uploads/"
	case ReqBlobUploadBlob:
		return "POST", "/v2/" + req.Repo + "/blobs/uploads/?digest=" + req.Digest
	case ReqBlobMount:
		return "POST", "/v2/" + req.Repo + "/blobs/uploads/?mount=" + req.Digest + "&from=" + req.FromRepo
	case ReqBlobUploadInfo:
		// Note: this is specific to the ociserver implementation.
		return "GET", req.uploadPath()
	case ReqBlobUploadChunk:
		// Note: this is specific to the ociserver implementation.
		return "PATCH", req.uploadPath()
	case ReqBlobCompleteUpload:
		// Note: this is specific to the ociserver implementation.
		// TODO this is bogus when the upload ID contains query parameters.
		return "PUT", req.uploadPath() + "?digest=" + req.Digest
	case ReqManifestGet:
		return "GET", "/v2/" + req.Repo + "/manifests/" + req.tagOrDigest()
	case ReqManifestHead:
		return "HEAD", "/v2/" + req.Repo + "/manifests/" + req.tagOrDigest()
	case ReqManifestPut:
		return "PUT", "/v2/" + req.Repo + "/manifests/" + req.tagOrDigest()
	case ReqManifestDelete:
		return "DELETE", "/v2/" + req.Repo + "/manifests/" + req.tagOrDigest()
	case ReqTagsList:
		return "GET", "/v2/" + req.Repo + "/tags/list" + req.listParams()
	case ReqReferrersList:
		return "GET", "/v2/" + req.Repo + "/referrers/" + req.Digest
	case ReqCatalogList:
		return "GET", "/v2/_catalog" + req.listParams()
	default:
		panic("invalid request kind")
	}
}

func (req *Request) uploadPath() string {
	return "/v2/" + req.Repo + "/blobs/uploads/" + base64.RawURLEncoding.EncodeToString([]byte(req.UploadID))
}

func (req *Request) listParams() string {
	q := make(url.Values)
	if req.ListN >= 0 {
		q.Set("n", fmt.Sprint(req.ListN))
	}
	if req.ListLast != "" {
		q.Set("last", req.ListLast)
	}
	if len(q) > 0 {
		return "?" + q.Encode()
	}
	return ""
}

func (req *Request) tagOrDigest() string {
	if req.Tag != "" {
		return req.Tag
	}
	return req.Digest
}
