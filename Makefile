.PHONY: help
help:
	$(info make options:)
	$(info -  make get)
	$(info -  make gen)
	$(info -  make patch)

.PHONY: get
get:
	wget https://alor.dev/rawdocs2/WarpOpenAPIv2.yml -O WarpOpenAPIv2.yml

.PHONY: gen
gen:
	which oapi-codegen
	oapi-codegen -package alor4go -generate types,client ./WarpOpenAPIv2.yml > alor4go.gen.go

.PHONY: patch
patch:
	git apply gen.patch
