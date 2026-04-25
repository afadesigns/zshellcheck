# Minimal runtime image for ZShellCheck.
# Goreleaser drops the pre-built static binary next to this Dockerfile during
# the multi-arch image build, so we don't need a builder stage. CGO_ENABLED=0
# in .goreleaser.yml keeps the binary fully static — no libc, no shell, no
# attack surface beyond zshellcheck itself. Final image is ~2 MB.

FROM scratch

COPY zshellcheck /zshellcheck

USER 65532:65532
ENTRYPOINT ["/zshellcheck"]
