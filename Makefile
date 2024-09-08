REPO_NAME=$(shell basename $(PWD))

rename:
	@echo "renaming repo references in files: ${REPO_NAME}"
	@for f in $(shell find . -not -type d -not -path "./.git/*" -not -path "*.enc"); do \
		sed -i '' 's/{{\.repoName}}/${REPO_NAME}/' $${f}; \
	done
	rm -rf ./Makefile