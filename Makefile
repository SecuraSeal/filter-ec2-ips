BUMP_VERSION := $(GOPATH)/bin/bump_version
RELEASE := $(GOPATH)/bin/github-release

release: $(BUMP_VERSION) $(RELEASE)
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	git tag $(version)
	git push origin --tags
	mkdir -p releases/$(version)
	# Change the binary names below to match your tool name
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/filter-ec2-ips-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/filter-ec2-ips-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/filter-ec2-ips-windows-amd64 .

	# Change the Github username to match your username.
	# These commands are not idempotent, so ignore failures if an upload repeats
	$(RELEASE) release --user SecuraSeal --repo filter-ec2-ips --tag $(version) || true
	$(RELEASE) upload --user SecuraSeal --repo filter-ec2-ips --tag $(version) --name filter-ec2-ips-linux-amd64 --file releases/$(version)/filter-ec2-ips-linux-amd64 || true
	$(RELEASE) upload --user SecuraSeal --repo filter-ec2-ips --tag $(version) --name filter-ec2-ips-darwin-amd64 --file releases/$(version)/filter-ec2-ips-darwin-amd64 || true
	$(RELEASE) upload --user SecuraSeal --repo filter-ec2-ips --tag $(version) --name filter-ec2-ips-windows-amd64 --file releases/$(version)/filter-ec2-ips-windows-amd64 || true
