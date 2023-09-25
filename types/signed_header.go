package types

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/celestiaorg/go-header"
	"github.com/cometbft/cometbft/crypto/ed25519"
)

func (sH *SignedHeader) New() header.Header {
	return new(SignedHeader)
}

func (sH *SignedHeader) IsZero() bool {
	return sH == nil
}

var (
	ErrAggregatorSetHashMismatch        = errors.New("aggregator set hash in signed header and hash of validator set do not match")
	ErrSignatureVerificationFailed      = errors.New("signature verification failed")
	ErrNoProposerAddress                = errors.New("no proposer address")
	ErrLastHeaderHashMismatch           = errors.New("last header hash mismatch")
	ErrLastCommitHashMismatch           = errors.New("last commit hash mismatch")
	ErrNewHeaderTimeBeforeOldHeaderTime = errors.New("new header has time before old header time")
	ErrNewHeaderTimeFromFuture          = errors.New("new header has time from future")
)

func (sH *SignedHeader) Verify(untrst header.Header) error {
	// Explicit type checks are required due to embedded Header which also does the explicit type check
	untrstH, ok := untrst.(*SignedHeader)
	if !ok {
		// if the header type is wrong, something very bad is going on
		// and is a programmer bug
		panic(fmt.Errorf("%T is not of type %T", untrst, untrstH))
	}
	if err := untrstH.ValidateBasic(); err != nil {
		return &header.VerifyError{
			Reason: err,
		}
	}
	if err := sH.Header.Verify(&untrstH.Header); err != nil {
		return &header.VerifyError{
			Reason: err,
		}
	}

	// TODO: Accept non-adjacent headers until go-header implements feature to accept non-adjacent
	if sH.Height()+1 < untrst.Height() {
		return nil
	}

	sHHash := sH.Header.Hash()
	if !bytes.Equal(untrstH.LastHeaderHash[:], sHHash) {
		return &header.VerifyError{
			Reason: fmt.Errorf("%w: expected %v, but got %v",
				ErrLastHeaderHashMismatch,
				untrstH.LastHeaderHash[:], sHHash,
			),
		}
	}
	sHLastCommitHash := sH.Commit.GetCommitHash(&untrstH.Header, sH.ProposerAddress)
	if !bytes.Equal(untrstH.LastCommitHash[:], sHLastCommitHash) {
		return &header.VerifyError{
			Reason: fmt.Errorf("%w: expected %v, but got %v",
				ErrLastCommitHashMismatch,
				untrstH.LastCommitHash[:], sHHash,
			),
		}
	}
	return nil
}

var _ header.Header = &SignedHeader{}

// ValidateBasic performs basic validation of a signed header.
func (h *SignedHeader) ValidateBasic() error {
	if err := h.Header.ValidateBasic(); err != nil {
		return err
	}

	if err := h.Commit.ValidateBasic(); err != nil {
		return err
	}

	// Handle Based Rollup case
	if h.Validators == nil || len(h.Validators.Validators) == 0 {
		return nil
	}

	if err := h.Validators.ValidateBasic(); err != nil {
		return err
	}

	if !bytes.Equal(h.Validators.Hash(), h.Signatures[0]) {
		return errors.New("aggregator set hash in signed header and hash of validator set do not match")
	}

	// Make sure there is exactly one signature
	if len(h.Commit.Signatures) != 1 {
		return errors.New("expected exactly one signature")
	}

	signature := h.Commit.Signatures[0]
	proposer := h.Validators.GetProposer()
	var pubKey ed25519.PubKey = proposer.PubKey.Bytes()
	msg, err := h.Header.MarshalBinary()
	if err != nil {
		return errors.New("signature verification failed, unable to marshal header")
	}
	if !pubKey.VerifySignature(msg, signature) {
		return ErrSignatureVerificationFailed
	}

	return nil
}
