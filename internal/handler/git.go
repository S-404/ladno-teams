package handler

import (
	"bytes"
	"fmt"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/gitstore"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
	"github.com/google/uuid"
)

// GitHandler serves smart HTTP git protocol for workspace bare repos.
type GitHandler struct {
	services interface {
		AuthorizeGitRead(userGuid, workspaceGuid uuid.UUID) *exception.ApiError
		AuthorizeGitWrite(userGuid, workspaceGuid uuid.UUID) *exception.ApiError
		EnsureGitRepo(workspaceGuid uuid.UUID) error
	}
	store *gitstore.Store
}

func NewGitHandler(workspaceSvc interface {
	AuthorizeGitRead(userGuid, workspaceGuid uuid.UUID) *exception.ApiError
	AuthorizeGitWrite(userGuid, workspaceGuid uuid.UUID) *exception.ApiError
	EnsureGitRepo(workspaceGuid uuid.UUID) error
}, store *gitstore.Store) *GitHandler {
	return &GitHandler{services: workspaceSvc, store: store}
}

func (h *GitHandler) InfoRefs(c *gin.Context) {
	guid, ok := h.parseGuid(c)
	if !ok {
		return
	}
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}

	service := c.Query("service")
	switch service {
	case transport.UploadPackServiceName:
		if apiErr := h.services.AuthorizeGitRead(user.Guid, guid); apiErr != nil {
			exception.HttpResponseException(c, apiErr)
			return
		}
	case transport.ReceivePackServiceName:
		if apiErr := h.services.AuthorizeGitWrite(user.Guid, guid); apiErr != nil {
			exception.HttpResponseException(c, apiErr)
			return
		}
	default:
		c.String(http.StatusForbidden, "unsupported service")
		return
	}

	if err := h.services.EnsureGitRepo(guid); err != nil {
		log.Printf("[git] EnsureGitRepo %s: %v", guid, err)
		c.String(http.StatusInternalServerError, "ensure repo: %v", err)
		return
	}

	repoPath := h.store.Path(guid)
	body, err := advertiseRefs(repoPath, service)
	if err != nil {
		log.Printf("[git] advertise %s %s: %v", service, repoPath, err)
		c.String(http.StatusInternalServerError, "git advertise: %v", err)
		return
	}

	ctype := fmt.Sprintf("application/x-%s-advertisement", service)
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, ctype, body)
}

func advertiseRefs(repoPath, service string) ([]byte, error) {
	ep, err := transport.NewEndpoint(repoPath)
	if err != nil {
		return nil, err
	}
	srv := server.NewServer(server.DefaultLoader)

	var ar *packp.AdvRefs
	switch service {
	case transport.UploadPackServiceName:
		sess, err := srv.NewUploadPackSession(ep, nil)
		if err != nil {
			return nil, fmt.Errorf("upload-pack session: %w", err)
		}
		defer sess.Close()
		ar, err = sess.AdvertisedReferences()
		if err != nil {
			return nil, fmt.Errorf("advertised refs: %w", err)
		}
	case transport.ReceivePackServiceName:
		sess, err := srv.NewReceivePackSession(ep, nil)
		if err != nil {
			return nil, fmt.Errorf("receive-pack session: %w", err)
		}
		defer sess.Close()
		ar, err = sess.AdvertisedReferences()
		if err != nil {
			return nil, fmt.Errorf("advertised refs: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported service %s", service)
	}

	ar.Prefix = [][]byte{
		[]byte("# service=" + service),
		pktline.Flush,
	}

	var buf bytes.Buffer
	if err := ar.Encode(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *GitHandler) UploadPack(c *gin.Context) {
	h.servePack(c, transport.UploadPackServiceName, false)
}

func (h *GitHandler) ReceivePack(c *gin.Context) {
	h.servePack(c, transport.ReceivePackServiceName, true)
}

// servePack handles HTTP POST git-upload-pack / git-receive-pack using go-git.
// Unlike the bidirectional git:// protocol, HTTP already advertised refs on info/refs,
// so the body is only the pack request (no second advertisement).
func (h *GitHandler) servePack(c *gin.Context, service string, write bool) {
	guid, ok := h.parseGuid(c)
	if !ok {
		return
	}
	user, err := GetCtxUser(c)
	if err != nil {
		exception.HttpResponseException(c, exception.AuthError("Unauthed"))
		return
	}
	if write {
		if apiErr := h.services.AuthorizeGitWrite(user.Guid, guid); apiErr != nil {
			exception.HttpResponseException(c, apiErr)
			return
		}
	} else {
		if apiErr := h.services.AuthorizeGitRead(user.Guid, guid); apiErr != nil {
			exception.HttpResponseException(c, apiErr)
			return
		}
	}
	if err := h.services.EnsureGitRepo(guid); err != nil {
		log.Printf("[git] EnsureGitRepo %s: %v", guid, err)
		c.String(http.StatusInternalServerError, "ensure repo: %v", err)
		return
	}

	repoPath := h.store.Path(guid)
	ep, err := transport.NewEndpoint(repoPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "endpoint: %v", err)
		return
	}
	srv := server.NewServer(server.DefaultLoader)
	ctx := c.Request.Context()

	c.Header("Content-Type", fmt.Sprintf("application/x-%s-result", service))
	c.Header("Cache-Control", "no-cache")

	if write {
		sess, err := srv.NewReceivePackSession(ep, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "receive-pack session: %v", err)
			return
		}
		defer sess.Close()
		// Initialize session capabilities (HTTP advertise already happened).
		if _, err := sess.AdvertisedReferences(); err != nil {
			c.String(http.StatusInternalServerError, "advertise: %v", err)
			return
		}
		req := packp.NewReferenceUpdateRequest()
		if err := req.Decode(c.Request.Body); err != nil {
			c.String(http.StatusBadRequest, "decode receive-pack: %v", err)
			return
		}
		rs, err := sess.ReceivePack(ctx, req)
		if err != nil {
			log.Printf("[git] receive-pack %s: %v", repoPath, err)
			// Still try to write report-status if present.
			if rs != nil {
				c.Status(http.StatusOK)
				_ = rs.Encode(c.Writer)
				return
			}
			c.String(http.StatusInternalServerError, "receive-pack: %v", err)
			return
		}
		c.Status(http.StatusOK)
		if rs != nil {
			if err := rs.Encode(c.Writer); err != nil {
				log.Printf("[git] encode report-status: %v", err)
			}
		}
		return
	}

	sess, err := srv.NewUploadPackSession(ep, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "upload-pack session: %v", err)
		return
	}
	defer sess.Close()
	if _, err := sess.AdvertisedReferences(); err != nil {
		c.String(http.StatusInternalServerError, "advertise: %v", err)
		return
	}
	req := packp.NewUploadPackRequest()
	if err := req.Decode(c.Request.Body); err != nil {
		c.String(http.StatusBadRequest, "decode upload-pack: %v", err)
		return
	}
	resp, err := sess.UploadPack(ctx, req)
	if err != nil {
		log.Printf("[git] upload-pack %s: %v", repoPath, err)
		c.String(http.StatusInternalServerError, "upload-pack: %v", err)
		return
	}
	c.Status(http.StatusOK)
	if err := resp.Encode(c.Writer); err != nil {
		log.Printf("[git] encode upload-pack response: %v", err)
	}
}

func (h *GitHandler) parseGuid(c *gin.Context) (uuid.UUID, bool) {
	raw := c.Param("guid")
	id, err := uuid.Parse(raw)
	if err != nil {
		exception.HttpResponseException(c, exception.BadRequest("invalid workspace guid"))
		return uuid.Nil, false
	}
	return id, true
}

// gitSubcommand maps protocol service names to `git <subcommand>` forms.
// Kept for documentation/tests: "git-upload-pack" → "upload-pack".
func gitSubcommand(service string) string {
	return strings.TrimPrefix(service, "git-")
}
