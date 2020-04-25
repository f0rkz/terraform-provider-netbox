build: test
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/{{.OS}}-{{.Arch}}-terraform-provider-netbox" .

test:
	go test -v $(shell go list ./... | grep -v /vendor/) 

testacc:
	TF_ACC=1 go test -v ./plugin/providers/netbox -run="TestAcc"

install: clean build
	cp pkg/linux_amd64/linux-amd64-terraform-provider-netbox ~/.terraform.d/plugins/terraform-provider-netbox
	
clean:
	rm -rf pkg/
