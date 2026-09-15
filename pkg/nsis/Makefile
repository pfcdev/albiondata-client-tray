TOP_SRCDIR = ../..
PACKAGE = albiondata-client
PACKAGE_NAME = Albion Data Client
PACKAGE_BUGREPORT = https://github.com/pfcdev/albiondata-client-tray/issues
INSTALLER_COMPRESSION = lzma
PACKAGE_EXE = albiondata-client.exe
RELEASE_VERSION ?= $(GITHUB_REF_NAME)

all: nsis

nsis:
	set -e; if echo "$(RELEASE_VERSION)" | grep -qE '^v?[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		export VERSION=$$(echo "$(RELEASE_VERSION)" | sed 's/^v//'); \
	else \
		export VERSION="0.0.1"; \
	fi; \
	makensis -DPACKAGE_VERSION="$$VERSION" -DPACKAGE="$(PACKAGE)" -DPACKAGE_NAME="$(PACKAGE_NAME)" -DPACKAGE_EXE="$(PACKAGE_EXE)" -DPACKAGE_BUGREPORT="$(PACKAGE_BUGREPORT)" -DTOP_SRCDIR="$(TOP_SRCDIR)" -DOUTFILE=../../albiondata-client-amd64-installer.exe -X"SetCompressor /FINAL /SOLID $(INSTALLER_COMPRESSION)" albiondata-client.nsi $<
